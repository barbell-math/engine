package model;

import (
    "fmt"
    "time"
    stdMath "math"
    "github.com/barbell-math/block/db"
    "github.com/barbell-math/block/util/dataStruct"
    mathUtil "github.com/barbell-math/block/util/math"
)

type SavedIntensityValue struct {
    intensity float64;
    date time.Time;
};

type SlidingWindowStateGen struct {
    allotedThreads int;
    windowLimits dataStruct.Pair[int];
    timeFrameLimits dataStruct.Pair[int];
    windowValues []SavedIntensityValue;
    optimalMs db.ModelState;
    lr mathUtil.LinearReg[float64];
    withinWindowLimits (func(t time.Time) bool);
};

//The lr and window values are not created until the generate prediction method
//in order to make sure the slice values are unique. (ie. It makes sure multiple
//threads don't access the same slices because they are not deep copied even
//though the method does not have a pointer receiver.)
func NewSlidingWindowStateGen(
        timeFrameLimits dataStruct.Pair[int],
        windowLimits dataStruct.Pair[int],
        allotedThreads int) (SlidingWindowStateGen,error) {
    rv:=SlidingWindowStateGen{
        allotedThreads: mathUtil.Constrain(allotedThreads,1,stdMath.MaxInt),
        timeFrameLimits: dataStruct.Pair[int]{
            First: -mathUtil.Abs(timeFrameLimits.First),
            Second: -mathUtil.Abs(timeFrameLimits.Second),
        }, windowLimits: dataStruct.Pair[int]{
            First: -mathUtil.Abs(windowLimits.First),
            Second: -mathUtil.Abs(windowLimits.Second),
        }, optimalMs: db.ModelState{
            Mse: stdMath.Inf(1),
        },
    };
    if rv.timeFrameLimits.First<=rv.timeFrameLimits.Second {
        return rv,InvalidPredictionState("min time frame > max time frame");
    } else if rv.windowLimits.First<=rv.windowLimits.Second {
        return rv,InvalidPredictionState("min window > max window");
    } else if rv.windowLimits.Second<rv.timeFrameLimits.Second {
        return rv,InvalidPredictionState("max window > max time frame");
    } else if rv.windowLimits.First>rv.timeFrameLimits.First {
        return rv,InvalidPredictionState("min window < min time frame");
    }
    return rv,nil;
}

func (s SlidingWindowStateGen)GenerateClientModelStates(
        d *db.DB,
        c db.Client,
        ch chan<- []error){
    stateGenType,err:=db.GetStateGeneratorByName(d,"Sliding Window");
    err=db.CustomReadQuery(d,missingModelStatesForGivenStateGenQuery(),[]any{
        c.Id,stateGenType.Id,
    }, func (m *missingModelStateData){
        fmt.Printf("Need ms for %+v\n",m);
        DEBUG.Log("Need ms: %+v\n",m);
    });
    fmt.Println(err);
}

func (s SlidingWindowStateGen)GenerateModelState(
        d *db.DB,
        missingData missingModelStateData,
        ch chan<- StateGeneratorRes){
    var curDate time.Time;
    s.setWithinWindowLimits(missingData.Date);
    s.lr=fatigueAwareModel();
    err:=db.CustomReadQuery(d,timeFrameQuery(),[]any{
        missingData.Date,
        missingData.Date.AddDate(0, 0, s.timeFrameLimits.Second),
        missingData.ExerciseID,
        missingData.ClientID,
    }, func (d *dataPoint){
        if !curDate.Equal(d.DatePerformed) && !missingData.Date.AddDate(
            0, 0, s.windowLimits.First,
        ).Before(curDate) {
            //calc model state, set
            curDate=d.DatePerformed;
        }
        //if within window, add to window vals
        if s.withinWindowLimits(curDate) {
            s.updateWindowValues(d);
        }
        s.updateLrSummations(d);
    });
    fmt.Println(err);
}

func (s *SlidingWindowStateGen)updateWindowValues(d *dataPoint){
    s.windowValues=append(s.windowValues,SavedIntensityValue{
        intensity: d.Intensity,
        date: d.DatePerformed,
    });
    DEBUG.Log("Added to window: %+v\n",d);
}

func (s *SlidingWindowStateGen)updateLrSummations(d *dataPoint){
    s.lr.UpdateSummations(map[string]float64{
        "I": d.Intensity, "R": d.Reps, "E": d.Effort, "S": d.Sets,
        "F_w": d.InterWorkoutFatigue, "F_e": d.InterExerciseFatigue,
    });
    DEBUG.Log("Added dp: %+v\n",d);
}

func (s *SlidingWindowStateGen)setWithinWindowLimits(
        missingDataTime time.Time){
    fmt.Println(missingDataTime.AddDate(0 ,0, s.windowLimits.First));
    fmt.Println(missingDataTime.AddDate(0 ,0, s.windowLimits.Second));
    s.withinWindowLimits=between(
        missingDataTime.AddDate(0, 0, s.windowLimits.First),
        missingDataTime.AddDate(0, 0, s.windowLimits.Second),
    );
}

//this will eventually need to be moved to some common utility file
//Due to the confusing nature of time, both after and before are **inclusive**
//After is the 'future' time, before is the 'past' time
func between(after time.Time, before time.Time) (func(t time.Time) bool) {
    //if the 'past' time is after the 'future' time then switch them
    if before.After(after) {
        tmp:=before;
        before=after;
        after=tmp;
    }
    return func(t time.Time) bool {
        return (before.AddDate(0, 0, -1).Before(t) &&
            after.AddDate(0, 0, 1).After(t));
    }
}
