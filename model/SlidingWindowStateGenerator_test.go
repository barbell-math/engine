package model;

import (
    "testing"
    "github.com/barbell-math/block/db"
    "github.com/barbell-math/block/util/test"
    "github.com/barbell-math/block/util/dataStruct"
)

func invalidCheck(slidingWindowSg SlidingWindowStateGen, err error) (func(t *testing.T)){
    return func(t *testing.T){
        if !IsInvalidPredictionState(err) {
            test.FormatError(InvalidPredictionState(""),err,
                "The wrong error was raised when creating an invalid prediction generator.",t,
            );
        }
    }
}
func TestNewSlidingWindowStateGenInvalidTimeFrameLimits(t *testing.T){
    invalidCheck(NewSlidingWindowStateGen(
        dataStruct.Pair[int]{1,0},dataStruct.Pair[int]{0, 1},1,
    ))(t);
}
func TestNewSlidingWindowStateGenInvalidWindowLimits(t *testing.T){
    invalidCheck(NewSlidingWindowStateGen(
        dataStruct.Pair[int]{0, 1},dataStruct.Pair[int]{1, 0},1,
    ))(t);
}
func TestNewSlidingWindowStateGenInvalidWindowSize(t *testing.T){
    invalidCheck(NewSlidingWindowStateGen(
        dataStruct.Pair[int]{0, 1},dataStruct.Pair[int]{0, 2},1,
    ))(t);
    invalidCheck(NewSlidingWindowStateGen(
        dataStruct.Pair[int]{1, 2},dataStruct.Pair[int]{0, 2},1,
    ))(t);
}
func TestNewSlidingWindowValid(t *testing.T){
    _,err:=NewSlidingWindowStateGen(
        dataStruct.Pair[int]{0, 1},dataStruct.Pair[int]{0, 1},1,
    );
    test.BasicTest(nil,err,
        "Creating a sliding window resulted in an error when it shouldn't have.",t,
    );
}

func TestNewSlidingWindowConstrainedThreadAllocation(t *testing.T){
    sw,err:=NewSlidingWindowStateGen(
        dataStruct.Pair[int]{0, 1},dataStruct.Pair[int]{0, 1},0,
    );
    test.BasicTest(nil,err,
        "Creating a sliding window resulted in an error when it shouldn't have.",t,
    );
    test.BasicTest(1,sw.allotedThreads,
        "The sliding window was allotted the wrong number of threads.",t,
    );
}

func TestGenerateModelStates(t *testing.T){
    ch:=make(chan<- []error);
    sw,_:=NewSlidingWindowStateGen(
        dataStruct.Pair[int]{0, 1},dataStruct.Pair[int]{0, 1},0,
    );
    sw.GenerateClientModelStates(&testDB,db.Client{ Id: 1 },ch);
}
