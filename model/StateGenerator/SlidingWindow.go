package stateGenerator

import (
	"database/sql"
	"fmt"
	stdMath "math"
	stdTime "time"

	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/util/algo/iter"
	"github.com/barbell-math/block/util/dataStruct"
	mathUtil "github.com/barbell-math/block/util/math/numeric"
	timeUtil "github.com/barbell-math/block/util/time"
	potSurf "github.com/barbell-math/block/model/potentialSurface"
	customerr "github.com/barbell-math/block/util/err"
)

type SlidingWindowStateGen struct {
    allotedThreads int;
    windowLimits dataStruct.Pair[int,int];
    timeFrameLimits dataStruct.Pair[int,int];
    windowValues [][]dataPoint;
    optimalMs []db.ModelState;
    models []potSurf.Surface;
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
        allotedThreads: mathUtil.Constrain(allotedThreads,dataStruct.Pair[int,int]{
            A: 1, B: stdMath.MaxInt,
        }), timeFrameLimits: dataStruct.Pair[int,int]{
            A: -mathUtil.Abs(timeFrameLimits.A),
            B: -mathUtil.Abs(timeFrameLimits.B),
        }, windowLimits: dataStruct.Pair[int,int]{
            A: -mathUtil.Abs(windowLimits.A),
            B: -mathUtil.Abs(windowLimits.B),
        }, 
    };
    if rv.timeFrameLimits.A<rv.timeFrameLimits.B {
        return rv,customerr.InvalidValue("min time frame > max time frame");
    } else if rv.windowLimits.A<rv.windowLimits.B {
        return rv,customerr.InvalidValue("min window >= max window, should be <");
    } else if rv.windowLimits.B<rv.timeFrameLimits.B {
        return rv,customerr.InvalidValue("max window >= max time frame, should be <");
    } else if rv.windowLimits.A>rv.timeFrameLimits.A {
        return rv,customerr.InvalidValue("min window <= min time frame, should be >");
    }
    return rv,nil;
}

func (s SlidingWindowStateGen)Id() StateGeneratorId {
    return SlidingWindowStateGenId;
}

//The method receiver is not a pointer so that the object will be copied. It is
//meant to be called in parallel (i.e. multiple clients) so the copy is necessary.
func (s SlidingWindowStateGen)GenerateClientModelStates(
        d *db.DB,
        c db.Client,
        minTime stdTime.Time,
        surfaceFactory func() []potSurf.Surface) (dataStruct.Pair[int,int],error) {
    rv:=dataStruct.Pair[int,int]{A: 0, B: 0};
    bufCreator,err:=db.NewBufferedCreate[db.ModelState](100);
    if err!=nil {
        return rv,err;
    }
    err=iter.Parallel[*missingModelStateData,[]db.ModelState](
        db.CustomReadQuery[missingModelStateData](d,
            missingModelStatesForGivenStateGenQuery(),[]any{
                c.Id,s.Id(),minTime,
        }),func(val *missingModelStateData) ([]db.ModelState, error) {
            return s.GenerateModelState(d,surfaceFactory(),val);
        },func(val *missingModelStateData, res []db.ModelState, err error) {
            //fmt.Println(err);
            if err==nil {
                for _,r:=range(res) {
                    bufCreator.Write(d,r);
                    SLIDING_WINDOW_MS_PARALLEL_RESULT_DEBUG.Log("Optimal MS",r);
                }
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

//The method receiver is not a pointer so that the object will be copied. It is
//meant to be called in parallel (i.e. multiple dates/exercises) so the copy
//is necessary.
func (s SlidingWindowStateGen)GenerateModelState(
        d *db.DB,
        surface []potSurf.Surface,
        missingData *missingModelStateData) ([]db.ModelState,error) {
    s.setInitialOptimalMsValues(missingData,surface);
    s.setWithinWindowLimits(missingData.Date);
    _,err:=s.runAlgo(d,missingData);
    if err==sql.ErrNoRows {
        err=NoDataInSelectedTimeFrame(fmt.Sprintf(
            "Date: %s Min time frame: %d Max time frame: %d Exercise: %d Client: %d",
            missingData.Date, s.timeFrameLimits.A, s.timeFrameLimits.B,
            missingData.ExerciseID, missingData.ClientID,
        ));
        return []db.ModelState{},err;
    }
    return s.optimalMs,err;
}

func (s *SlidingWindowStateGen)setInitialOptimalMsValues(
        missingData *missingModelStateData,
        surfaces []potSurf.Surface){
    s.models=surfaces;
    s.optimalMs=make([]db.ModelState,len(s.models));
    for i,_:=range(s.models) {
        s.optimalMs[i]=db.ModelState{
            Mse: stdMath.Inf(1),
            Date: missingData.Date,
            ClientID: missingData.ClientID,
            ExerciseID: missingData.ExerciseID,
            StateGeneratorID: int(SlidingWindowStateGenId),
            PotentialSurfaceID: int(s.models[i].Id()),
        };
    }
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
        s.updateModel(val);    //Time frame limits guaranteed by query
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
    for i,m:=range(s.models) {
        var numPoints float64=0;
        var cumulativeSe float64=0.0;
        rcond,_:=m.Run();
        curStability:=m.Stability();
        for j,w:=range(s.windowValues) {
            totalSe,err:=s.getActualAndPredSe(m,j);
            if err!=nil {
                return err;
            }
            numPoints+=float64(len(w));
            cumulativeSe+=totalSe;
            oldStability:=m.Calculations().Stability(&s.optimalMs[i]);
            if err==nil && (
                curStability>oldStability || (curStability==oldStability && 
                cumulativeSe/numPoints<s.optimalMs[i].Mse)) {
                s.saveModelState(i,rcond,cumulativeSe/numPoints,
                    timeUtil.DaysBetween(missingData.Date,w[0].DatePerformed),
                    timeUtil.DaysBetween(missingData.Date,d.DatePerformed),
                );
            }
        }
    }
    return nil;
}

func (s *SlidingWindowStateGen)getActualAndPredSe(m potSurf.Surface, 
        windowIndex int) (float64,error) {
    actual,pred:=make(
        []float64,len(s.windowValues[windowIndex]),
    ),make([]float64,len(s.windowValues[windowIndex]));
    for i,v:=range(s.windowValues[windowIndex]) {
        actual[i]=v.Intensity;
        if iterPred,err:=m.PredictIntensity(map[string]float64{
            "I": v.Intensity,
            "E": v.Effort,
            "R": float64(v.Reps),
            "S": float64(v.Sets),
            "F_w": float64(v.InterWorkoutFatigue),
            "F_e": float64(v.InterExerciseFatigue),
        }); err==nil {
            pred[i]=iterPred;
        }
    }
    cumulativeSe:=0.0;
    se,err:=mathUtil.SqErr(actual,pred);
    if err==nil {
        for _,v:=range(se) {
            cumulativeSe+=v;
        }
    }
    return cumulativeSe,err;
}

func (s *SlidingWindowStateGen)saveModelState(
        i int,
        rcond float64,
        mse float64,
        winLen int,
        timeFrameLen int){
    s.optimalMs[i].Eps=s.models[i].GetConstant(0);
    s.optimalMs[i].Eps1=s.models[i].GetConstant(1);
    s.optimalMs[i].Eps2=s.models[i].GetConstant(2);
    s.optimalMs[i].Eps3=s.models[i].GetConstant(3);
    s.optimalMs[i].Eps4=s.models[i].GetConstant(4);
    s.optimalMs[i].Eps5=s.models[i].GetConstant(5);
    s.optimalMs[i].Eps6=s.models[i].GetConstant(6);
    s.optimalMs[i].Eps7=s.models[i].GetConstant(7);
    s.optimalMs[i].TimeFrame=timeFrameLen;
    s.optimalMs[i].Win=winLen;
    s.optimalMs[i].Rcond=rcond;
    s.optimalMs[i].Mse=mse;
    SLIDING_WINDOW_MS_DEBUG.Log("ModelState",s.optimalMs[i]);
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

func (s *SlidingWindowStateGen)updateModel(d *dataPoint){
    for _,m:=range(s.models) {
        m.Update(map[string]float64{
            "I": d.Intensity, "R": d.Reps, "E": d.Effort, "S": d.Sets,
            "F_w": d.InterWorkoutFatigue, "F_e": d.InterExerciseFatigue,
        });
    }
    SLIDING_WINDOW_DP_DEBUG.Log("DataPoint",d);
}

func (s *SlidingWindowStateGen)setWithinWindowLimits(
        missingDataTime stdTime.Time){
    s.withinWindowLimits=timeUtil.Between(
        missingDataTime.AddDate(0, 0, s.windowLimits.A),
        missingDataTime.AddDate(0, 0, s.windowLimits.B),
    );
}
