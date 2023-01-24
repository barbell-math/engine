package model;

import (
    "fmt"
    "time"
    stdMath "math"
    "github.com/barbell-math/block/db"
    mathUtil "github.com/barbell-math/block/util/math"
    "github.com/barbell-math/block/util/dataStruct"
)

type SavedIntensityValue struct {
    intensity float64;
    date time.Time;
};

type SlidingWindowStateGen struct {
    window int;
    allotedThreads int;
    windowLimits dataStruct.Pair[int];
    timeFrameLimits dataStruct.Pair[int];
    windowValues []SavedIntensityValue;
    optimalMs db.ModelState;
    lr mathUtil.LinearReg[float64];
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
        allotedThreads: allotedThreads,
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
    }, func (t *db.TrainingLog){
        fmt.Printf("Need ms for %+v\n",t);
    });
    fmt.Println(err);
}
