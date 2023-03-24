package model;

import (
    "fmt"
    "time"
    "testing"
    //"github.com/barbell-math/block/db"
    "github.com/barbell-math/block/util/log"
    "github.com/barbell-math/block/util/test"
    "github.com/barbell-math/block/util/dataStruct/base"
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
        base.Pair[int,int]{1,0},base.Pair[int,int]{0, 1},1,
    ))(t);
}
func TestNewSlidingWindowStateGenInvalidWindowLimits(t *testing.T){
    invalidCheck(NewSlidingWindowStateGen(
        base.Pair[int,int]{0, 1},base.Pair[int,int]{1, 0},1,
    ))(t);
}
func TestNewSlidingWindowStateGenInvalidWindowSize(t *testing.T){
    invalidCheck(NewSlidingWindowStateGen(
        base.Pair[int,int]{0, 1},base.Pair[int,int]{0, 2},1,
    ))(t);
    invalidCheck(NewSlidingWindowStateGen(
        base.Pair[int,int]{1, 2},base.Pair[int,int]{0, 2},1,
    ))(t);
}
func TestNewSlidingWindowValid(t *testing.T){
    _,err:=NewSlidingWindowStateGen(
        base.Pair[int,int]{0, 1},base.Pair[int,int]{0, 1},1,
    );
    test.BasicTest(nil,err,
        "Creating a sliding window resulted in an error when it shouldn't have.",t,
    );
}

func TestNewSlidingWindowConstrainedThreadAllocation(t *testing.T){
    sw,err:=NewSlidingWindowStateGen(
        base.Pair[int,int]{0, 1},base.Pair[int,int]{0, 1},0,
    );
    test.BasicTest(nil,err,
        "Creating a sliding window resulted in an error when it shouldn't have.",t,
    );
    test.BasicTest(1,sw.allotedThreads,
        "The sliding window was allotted the wrong number of threads.",t,
    );
}

//func TestGenerateModelStates(t *testing.T){
//    setLogs("./debugLogs/SlidingWindow.log");
//    ch:=make(chan<- []error);
//    sw,_:=NewSlidingWindowStateGen(
//        dataStruct.Pair[int]{0, 1},dataStruct.Pair[int]{0, 1},0,
//    );
//    sw.GenerateClientModelStates(&testDB,db.Client{ Id: 1 },ch);
//}

func TestGenerateModelState(t *testing.T){
    tmp:=setupLogs("./debugLogs/SlidingWindowStateGeneratorGood.log");
    baseTime,_:=time.Parse("01/02/2006","09/10/2022");
    ch:=make(chan<- StateGeneratorRes);
    sw,_:=NewSlidingWindowStateGen(
        //base.Pair[int,int]{4, 500},base.Pair[int,int]{5, 10},0,
        //base.Pair[int,int]{0, 500},base.Pair[int,int]{0, 1},0,
        base.Pair[int,int]{0, 500},base.Pair[int,int]{0, 10},0,
    );
    err:=sw.GenerateModelState(&testDB,missingModelStateData{
        ClientID: 1,
        ExerciseID: 15,
        Date: baseTime,
    },ch);
    fmt.Println("ERR: ",err);
    tmp();
    log.LogElems[dataPoint](DEBUG).ForEach(func(index int, val log.LogEntry[dataPoint]) error {
        fmt.Printf("%+v\t %s\t %s\n",val,val.T,val.V.DatePerformed);
        return nil;
    });
}
