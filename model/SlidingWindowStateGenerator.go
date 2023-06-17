package model

import (
	"database/sql"
	"fmt"
	"math"
	stdMath "math"
	stdTime "time"

	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/util/algo/iter"
	"github.com/barbell-math/block/util/dataStruct"
	mathUtil "github.com/barbell-math/block/util/math"
	timeUtil "github.com/barbell-math/block/util/time"
)

type SlidingWindowStateGen struct {
    allotedThreads int;
    windowLimits dataStruct.Pair[int,int];
    timeFrameLimits dataStruct.Pair[int,int];
    windowValues [][]dataPoint;
    optimalMs db.ModelState;
    lr mathUtil.LinearReg[float64];
    withinWindowLimits (func(t stdTime.Time) bool);
};

//The lr and window values are not created until the generate prediction method
//in order to make sure the slice values are unique. (ie. It makes sure multiple
//threads don't access the same slices because they are not deep copied even
//though the method does not have a pointer receiver.)
func NewSlidingWindowStateGen(
        timeFrameLimits dataStruct.Pair[int,int],
        windowLimits dataStruct.Pair[int,int],
        allotedThreads int) (SlidingWindowStateGen,error) {
    rv:=SlidingWindowStateGen{
        allotedThreads: mathUtil.Constrain(allotedThreads,1,stdMath.MaxInt),
        timeFrameLimits: dataStruct.Pair[int,int]{
            A: -mathUtil.Abs(timeFrameLimits.A),
            B: -mathUtil.Abs(timeFrameLimits.B),
        }, windowLimits: dataStruct.Pair[int,int]{
            A: -mathUtil.Abs(windowLimits.A),
            B: -mathUtil.Abs(windowLimits.B),
        }, optimalMs: db.ModelState{
            Mse: stdMath.Inf(1),
        },
    };
    if rv.timeFrameLimits.A<rv.timeFrameLimits.B {
        return rv,InvalidPredictionState("min time frame > max time frame");
    } else if rv.windowLimits.A<rv.windowLimits.B {
        return rv,InvalidPredictionState("min window >= max window, should be <");
    } else if rv.windowLimits.B<rv.timeFrameLimits.B {
        return rv,InvalidPredictionState("max window >= max time frame, should be <");
    } else if rv.windowLimits.A>rv.timeFrameLimits.A {
        return rv,InvalidPredictionState("min window <= min time frame, should be >");
    }
    return rv,nil;
}

func (s SlidingWindowStateGen)Id() StateGeneratorId {
    return SlidingWindowStateGenId;
}

func (s SlidingWindowStateGen)GenerateClientModelStates(
        d *db.DB,
        c db.Client,
        minTime stdTime.Time) (dataStruct.Pair[int,int],error) {
    rv:=dataStruct.Pair[int,int]{A: 0, B: 0};
    bufCreator,err:=db.NewBufferedCreate[db.ModelState](100);
    if err!=nil {
        return rv,err;
    }
    err=iter.Parallel[*missingModelStateData,db.ModelState](
        db.CustomReadQuery[missingModelStateData](d,
            missingModelStatesForGivenStateGenQuery(),[]any{
                c.Id,s.Id(),minTime,
        }),func(val *missingModelStateData) (db.ModelState, error) {
            return s.GenerateModelState(d,val);
        },func(val *missingModelStateData, res db.ModelState, err error) {
            SLIDING_WINDOW_MS_PARALLEL_RESULT_DEBUG.Log("Optimal MS",res);
            //fmt.Println(err);
            if err==nil {
                bufCreator.Write(d,res);
            } else {
                rv.B++;
            }
        },s.allotedThreads,
    );
    bufCreator.Flush(d);
    rv.A=bufCreator.Succeeded();
    rv.B+=bufCreator.Failed();
    return rv,err;
}

func (s SlidingWindowStateGen)GenerateModelState(
        d *db.DB,
        missingData *missingModelStateData) (db.ModelState,error) {
    s.setInitialOptimalMsValues(missingData);
    s.setWithinWindowLimits(missingData.Date);
    s.lr=fatigueAwareModel();
    cntr,err:=s.runAlgo(d,missingData);
    if err==sql.ErrNoRows {
        err=NoDataInSelectedTimeFrame(fmt.Sprintf(
            "Date: %s Min time frame: %d Max time frame: %d Exercise: %d Client: %d",
            missingData.Date, s.timeFrameLimits.A, s.timeFrameLimits.B,
            missingData.ExerciseID, missingData.ClientID,
        ));
    } else if s.optimalMs.Mse==math.Inf(1) && err==nil {
        err=NotEnoughData(fmt.Sprintf(
            "Num data pnts: %d Date: %s Min time frame: %d Max time frame: %d Exercise: %d Client: %d",
            cntr,missingData.Date, s.timeFrameLimits.A, s.timeFrameLimits.B,
            missingData.ExerciseID, missingData.ClientID,
        ));
    }
    return s.optimalMs,err;
}

func (s *SlidingWindowStateGen)setInitialOptimalMsValues(
        missingData *missingModelStateData){
    s.optimalMs.ClientID=missingData.ClientID;
    s.optimalMs.ExerciseID=missingData.ExerciseID;
    s.optimalMs.StateGeneratorID=int(SlidingWindowStateGenId);
    s.optimalMs.Date=missingData.Date;
}

func (s *SlidingWindowStateGen)runAlgo(
    d *db.DB,
    missingData *missingModelStateData,
) (int,error) {
    cntr:=0;
    var curDate stdTime.Time;
    err:=db.CustomReadQuery[dataPoint](d,timeFrameQuery(),[]any{
        missingData.Date.AddDate(0, 0, s.timeFrameLimits.A),
        missingData.Date.AddDate(0, 0, s.timeFrameLimits.B),
        missingData.ExerciseID,
        missingData.ClientID,
    }).ForEach(func(index int, val *dataPoint) (iter.IteratorFeedback, error) {
        if !curDate.Equal(val.DatePerformed) {
            if len(s.windowValues)>0 {
                s.calcAndSetModelState(val,missingData);
            }
            curDate=val.DatePerformed;
        }
        if s.withinWindowLimits(curDate) {
            s.updateWindowValues(val);
        }
        s.updateLrSummations(val);    //Time frame limits guaranteed by query
        cntr++;
        if curDate.Before(
            missingData.Date.AddDate(0, 0, s.windowLimits.B),
        ) && len(s.windowValues)==0 {
            return iter.Break,NoDataInSelectedWindow(fmt.Sprintf(
                "Date: %s Min window: %d Max window: %d Exercise: %d Client: %d",
                missingData.Date, s.windowLimits.A, s.windowLimits.B,
                missingData.ExerciseID, missingData.ClientID,
            ));
        }
        return iter.Continue,nil;
    });
    return cntr,err;
}

//Algo steps:
//  1. For each day in the window values array
//      1. Generate a prediction for every value in that day
//      2. Calculate the square error (SE) for every value in that day. If the
//         pred and actual lists are different lengths then it means there was
//         an error generating intensity predictions and the model state should
//         be dis-regarded
//      3. Assuming the previous step succeeded, calculate the total SE
//      4. Calculate the mean SE (MSE)
//      5. If the MSE is less than the current lowest MSE then save a new
//         model state
func (s *SlidingWindowStateGen)calcAndSetModelState(
        d *dataPoint,
        missingData *missingModelStateData) error {
    var numPoints float64=0;
    var cumulativeSe float64=0.0;
    res,rcond,_:=s.lr.Run();
    for _,w:=range(s.windowValues) {
        actual,pred:=make([]float64,len(w)),make([]float64,len(w));
        for i,v:=range(w) {
            actual[i]=v.Intensity;
            if iterPred,err:=intensityPredFromLinReg(res,&v); err==nil {
                pred[i]=iterPred;
            }
        }
        if se,err:=mathUtil.SqErr(actual,pred); err==nil {
            seTot:=0.0;
            for _,v:=range(se) {
                seTot+=v;
            }
            cumulativeSe+=seTot;
            numPoints+=float64(len(w));
            if cumulativeSe/numPoints<s.optimalMs.Mse {
                s.saveModelState(res,rcond,cumulativeSe/numPoints,
                    timeUtil.DaysBetween(missingData.Date,w[0].DatePerformed),
                    timeUtil.DaysBetween(missingData.Date,d.DatePerformed),
                );
            }
        } else {
            return err;
        }
    }
    return nil;
}

func (s *SlidingWindowStateGen)saveModelState(
        res mathUtil.LinRegResult[float64],
        rcond float64,
        mse float64,
        winLen int,
        timeFrameLen int){
    s.optimalMs.Eps=stdMath.Max(res.GetConstant(0),0);
    //s.optimalMs.Eps1=res.GetConstant(1);
    s.optimalMs.Eps2=res.GetConstant(1);
    s.optimalMs.Eps3=res.GetConstant(2);
    s.optimalMs.Eps4=res.GetConstant(3);
    s.optimalMs.Eps5=stdMath.Max(res.GetConstant(4),0);
    s.optimalMs.Eps6=stdMath.Max(res.GetConstant(5),0);
    s.optimalMs.Eps7=stdMath.Max(res.GetConstant(6),0);
    s.optimalMs.TimeFrame=timeFrameLen;
    s.optimalMs.Win=winLen;
    s.optimalMs.Rcond=rcond;
    s.optimalMs.Mse=mse;
    if mse==math.Inf(1) {
        fmt.Println("Saved model state is inf");
    }
    SLIDING_WINDOW_MS_DEBUG.Log("ModelState",s.optimalMs);
}

func (s *SlidingWindowStateGen)updateWindowValues(d *dataPoint){
    l:=len(s.windowValues);
    if l==0 || !s.windowValues[l-1][0].DatePerformed.Equal(d.DatePerformed) {
        s.windowValues=append(s.windowValues,[]dataPoint{*d});
    } else {
        s.windowValues[l-1]=append(s.windowValues[l-1],*d);
    }
    SLIDING_WINDOW_DP_DEBUG.Log("WindowDataPoint",d);
}

func (s *SlidingWindowStateGen)updateLrSummations(d *dataPoint){
    s.lr.UpdateSummations(map[string]float64{
        "I": d.Intensity, "R": d.Reps, "E": d.Effort, "S": d.Sets,
        "F_w": d.InterWorkoutFatigue, "F_e": d.InterExerciseFatigue,
    });
    SLIDING_WINDOW_DP_DEBUG.Log("DataPoint",d);
}

func (s *SlidingWindowStateGen)setWithinWindowLimits(
        missingDataTime stdTime.Time){
    //fmt.Println(missingDataTime.AddDate(0 ,0, s.windowLimits.A));
    //fmt.Println(missingDataTime.AddDate(0 ,0, s.windowLimits.B));
    s.withinWindowLimits=timeUtil.Between(
        missingDataTime.AddDate(0, 0, s.windowLimits.A),
        missingDataTime.AddDate(0, 0, s.windowLimits.B),
    );
}
