package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/util/algo/iter"
	"github.com/barbell-math/block/util/dataStruct"
	"github.com/barbell-math/block/util/dataStruct/types"
	"github.com/barbell-math/block/util/io/log"
	"github.com/barbell-math/block/util/test"
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
        dataStruct.Pair[int,int]{A: 1,B: 0},dataStruct.Pair[int,int]{A: 0, B: 1},1,
    ))(t);
}
func TestNewSlidingWindowStateGenInvalidWindowLimits(t *testing.T){
    invalidCheck(NewSlidingWindowStateGen(
        dataStruct.Pair[int,int]{A: 0, B: 1},dataStruct.Pair[int,int]{A: 1, B: 0},1,
    ))(t);
}
func TestNewSlidingWindowStateGenInvalidWindowSize(t *testing.T){
    invalidCheck(NewSlidingWindowStateGen(
        dataStruct.Pair[int,int]{A: 0, B: 1},dataStruct.Pair[int,int]{A: 0, B: 2},1,
    ))(t);
    invalidCheck(NewSlidingWindowStateGen(
        dataStruct.Pair[int,int]{A: 1, B: 2},dataStruct.Pair[int,int]{A: 0, B: 2},1,
    ))(t);
}
func TestNewSlidingWindowValid(t *testing.T){
    _,err:=NewSlidingWindowStateGen(
        dataStruct.Pair[int,int]{A: 0, B: 1},dataStruct.Pair[int,int]{A: 0, B: 1},1,
    );
    test.BasicTest(nil,err,
        "Creating a sliding window resulted in an error when it shouldn't have.",t,
    );
}

func TestNewSlidingWindowConstrainedThreadAllocation(t *testing.T){
    sw,err:=NewSlidingWindowStateGen(
        dataStruct.Pair[int,int]{A: 0, B: 1},dataStruct.Pair[int,int]{A: 0, B: 1},0,
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
    tmp:=setupLogs("./debugLogs/SlidingWindowStateGeneratorGood");
    baseTime,_:=time.Parse("01/02/2006","09/10/2022");
    ch:=make(chan<- StateGeneratorRes);
    sw,_:=NewSlidingWindowStateGen(
        //dataStruct.Pair[int,int]{4, 500},dataStruct.Pair[int,int]{5, 10},0,
        //dataStruct.Pair[int,int]{0, 500},dataStruct.Pair[int,int]{0, 1},0,
        dataStruct.Pair[int,int]{A: 0, B: 500},dataStruct.Pair[int,int]{A: 0, B: 10},0,
    );
    err:=sw.GenerateModelState(&testDB,missingModelStateData{
        ClientID: 1,
        ExerciseID: 15,
        Date: baseTime,
    },ch);
    fmt.Println("ERR: ",err);
    tmp();
    iter.Join[log.LogEntry[*dataPoint],log.LogEntry[db.ModelState]](
        log.LogElems(SLIDING_WINDOW_DP_DEBUG),log.LogElems(SLIDING_WINDOW_MS_DEBUG),
        dataStruct.Variant[log.LogEntry[*dataPoint],log.LogEntry[db.ModelState]]{},
        log.JoinLogByTime[*dataPoint,db.ModelState],
    ).ForEach(func(index int,
        val types.Variant[log.LogEntry[*dataPoint],log.LogEntry[db.ModelState]],
    ) (iter.IteratorFeedback, error) {
        if val.HasA() {
            fmt.Printf("%+v\n",val.ValA());
        } else {
            fmt.Printf("%+v\n",val.ValB());
        }
        return iter.Continue,nil;
    })
}
