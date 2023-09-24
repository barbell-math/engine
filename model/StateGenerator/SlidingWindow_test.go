package stateGenerator

import (
	"fmt"
	"testing"
	"time"

	"github.com/barbell-math/engine/db"
	"github.com/barbell-math/engine/util/algo/iter"
	"github.com/barbell-math/engine/util/dataStruct"
	"github.com/barbell-math/engine/util/io/log"
	"github.com/barbell-math/engine/util/test"
	potSurf "github.com/barbell-math/engine/model/potentialSurface"
	customerr "github.com/barbell-math/engine/util/err"
)

func invalidCheck(slidingWindowSg SlidingWindowStateGen, err error) (func(t *testing.T)){
    return func(t *testing.T){
        if !customerr.IsInvalidValue(err) {
            test.FormatError(customerr.InvalidValue(""),err,
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

func TestNoDataForModelState(t *testing.T){
    baseTime,_:=time.Parse("01/02/2006","09/09/2022");
    timeFrame:=dataStruct.Pair[int,int]{A: 0, B: 500};
    window:=dataStruct.Pair[int,int]{A: 0, B: 1};
    _,err:=generateModelStateHelper("noWinData",baseTime,timeFrame,window,0,t);
    if !IsNoDataInSelectedWindow(err) {
        test.FormatError(NoDataInSelectedWindow(""),err,
            "The incorrect error was returned.",t,
        );
    }
}

func TestGenerateModelStateScenario1(t *testing.T){
    //Window limits: 8/31/2022-9/10/2022
    //Looking at the data there are four deadlift entries in that time span:
    //  - Two on 9/1/2022
    //  - One on 9/7/2022
    //  - One on 9/10/2022
    baseTime,_:=time.Parse("01/02/2006","09/10/2022");
    timeFrame:=dataStruct.Pair[int,int]{A: 0, B: 500};
    window:=dataStruct.Pair[int,int]{A: 0, B: 10};
    _,err:=generateModelStateHelper("scenario1",baseTime,timeFrame,window,4,t);
    test.BasicTest(nil,err,
        "Running the sliding window model state generator returned an error when it shouldn't have.",t,
    );
}

func TestGenerateModelStateScenario2(t *testing.T){
    //Window limits: 8/31/2022-9/5/2022
    //Looking at the data there are three deadlift entries in that time span:
    //  - Two on 9/1/2022
    baseTime,_:=time.Parse("01/02/2006","09/10/2022");
    timeFrame:=dataStruct.Pair[int,int]{A: 4, B: 500};
    window:=dataStruct.Pair[int,int]{A: 5, B: 10};
    _,err:=generateModelStateHelper("scenario2",baseTime,timeFrame,window,2,t);
    test.BasicTest(nil,err,
        "Running the sliding window model state generator returned an error when it shouldn't have.",t,
    );
}

func generateModelStateHelper(scenarioName string,
        baseTime time.Time,
        timeFrame dataStruct.Pair[int,int],
        window dataStruct.Pair[int,int],
        numWindowVals int,
        t *testing.T) (db.ModelState,error) {
    closeLogs:=setupLogs(fmt.Sprintf(
        "./debugLogs/SlidingWindowStateGenerator.%s",scenarioName,
    ),DP_DEBUG|MS_DEBUG);
    sw,_:=NewSlidingWindowStateGen(timeFrame,window,1);
    missingData:=missingModelStateData{
        ClientID: 1,
        ExerciseID: 15,
        Date: baseTime,
    };
    s:=[]potSurf.Surface{potSurf.NewBasicSurface().ToGenericSurf()};
    ms,err:=sw.GenerateModelState(&testDB,s,&missingData);
    closeLogs();
    if len(ms)>0 {
        runModelStateDebugLogTests(baseTime,
            missingData.ClientID,missingData.ExerciseID,int(SlidingWindowStateGenId),
            timeFrame,window,
            ms[0].Mse,potSurf.BasicSurfaceCalculation.Stability(&ms[0]),
            t,
        );
    }
    runDataPointDebugLogTests(baseTime,t);
    runWindowDataPointDebugLogTests(baseTime,window,numWindowVals,t);
    if len(ms)>0 {
        return ms[0],err;
    }
    return db.ModelState{},err;
}

func runDataPointDebugLogTests(baseTime time.Time, t *testing.T){
    initialDate:=time.Time{};
    baseTime=baseTime.AddDate(0,0,1);
    err:=log.LogElems(SLIDING_WINDOW_DP_DEBUG).Filter(
    func(index int, val log.LogEntry[*dataPoint]) bool {
        return val.Message=="DataPoint";
    }).Next(func(index int, 
        val log.LogEntry[*dataPoint],
        status iter.IteratorFeedback,
    ) (iter.IteratorFeedback, log.LogEntry[*dataPoint], error) {
        if status!=iter.Break {
            test.BasicTest(
                true,
                val.Val.DatePerformed.Before(baseTime),
                "Window values occurred at or after the current time! (Implies generated data is no longer a prediction!!)",t,
            );
        }
        return iter.Continue,val,nil;
    }).Filter(
    func(index int, val log.LogEntry[*dataPoint]) bool {
        if index==0 {
            initialDate=val.Val.DatePerformed;
            return false;
        }
        return true;
    }).ForEach(
    func(index int, val log.LogEntry[*dataPoint]) (iter.IteratorFeedback, error) {
        test.BasicTest(true,initialDate.Sub(val.Val.DatePerformed)>=0,
            "Training log dates did not continually decrease from query.",t,
        );
        initialDate=val.Val.DatePerformed;
        return iter.Continue,nil;
    });
    test.BasicTest(nil,err,
        "Iterating over a log generated an error when it should not have.",t,
    );
}

func runWindowDataPointDebugLogTests(baseTime time.Time,
        window dataStruct.Pair[int,int],
        numWindowVals int,
        t *testing.T){
    initialDate:=time.Time{};
    baseTime=baseTime.AddDate(0,0,1);
    cnt,err:=log.LogElems(SLIDING_WINDOW_DP_DEBUG).Filter(
    func(index int, val log.LogEntry[*dataPoint]) bool {
        return val.Message=="WindowDataPoint";
    }).Next(func(index int, 
        val log.LogEntry[*dataPoint],
        status iter.IteratorFeedback,
    ) (iter.IteratorFeedback, log.LogEntry[*dataPoint], error) {
        if status!=iter.Break {
            test.BasicTest(true,
                val.Val.DatePerformed.After(baseTime.AddDate(0, 0, -window.B-1)),
                "Window value is not after oldest allowed window value.",t,
            );
            test.BasicTest(true,
                val.Val.DatePerformed.Before(baseTime.AddDate(0, 0, -window.A)),
                "Window value is not before earliest allowed window value.",t,
            );
            test.BasicTest(true,val.Val.DatePerformed.Before(baseTime),
                "Window values occurred after the current time! (Implies generated data is no longer a prediction!!)",t,
            );
        }
        return iter.Continue,val,nil;
    }).Filter(func(index int, val log.LogEntry[*dataPoint]) bool {
        if index==0 {
            initialDate=val.Val.DatePerformed;
            return false;
        }
        return true;
    }).Next(func(index int,
        val log.LogEntry[*dataPoint],
        status iter.IteratorFeedback,
    ) (iter.IteratorFeedback, log.LogEntry[*dataPoint], error) {
        if status!=iter.Break {
            test.BasicTest(true,initialDate.Sub(val.Val.DatePerformed)>=0,
                "Training log dates did not continually decrease from query.",t,
            );
            initialDate=val.Val.DatePerformed;
        }
        return iter.Continue,val,nil;
    }).Count();
    if numWindowVals>0 {
        //The first filter removes the first value, making for an off by one error here
        test.BasicTest(numWindowVals-1,cnt,
            "The correct number of window values were not generated.",t,
        );
    }
    test.BasicTest(nil,err,
        "Iterating over a log generated an error when it should not have.",t,
    );
}

func runModelStateDebugLogTests(baseTime time.Time, cId int, eId int, sId int,
        timeFrame dataStruct.Pair[int,int],
        window dataStruct.Pair[int,int],
        optimalMse float64,
        optimalStability int,
        t *testing.T){
    initialMse:=0.0;
    initialStability:=0;
    err:=log.LogElems(SLIDING_WINDOW_MS_DEBUG).Next(
    func(index int,
        val log.LogEntry[db.ModelState],
        status iter.IteratorFeedback,
    ) (iter.IteratorFeedback, log.LogEntry[db.ModelState], error) {
        if status!=iter.Break {
            test.BasicTest(true,val.Val.TimeFrame>=timeFrame.A,
                "A model state had a time frame less than the selected lowest value.",t,
            );
            test.BasicTest(true,val.Val.TimeFrame<=timeFrame.B,
                "A model state had a time frame greater than the selected highest value.",t,
            );
            test.BasicTest(true,val.Val.Win>=window.A,
                "A model state had a window less than the selected lowest value.",t,
            );
            test.BasicTest(true,val.Val.Win<=window.B,
                "A model state had a window greater than the selected highest value.",t,
            );
            test.BasicTest(cId,val.Val.ClientID,
                "A model state had the incorrect client ID.",t,
            );
            test.BasicTest(eId,val.Val.ExerciseID,
                "A model state had the incorrect client ID.",t,
            );
            test.BasicTest(sId,val.Val.StateGeneratorID,
                "A model state had the incorrect state generator ID.",t,
            );
            y1,m1,d1:=val.Val.Date.Date();
            y2,m2,d2:=baseTime.Date();
            test.BasicTest(y2,y1,"A model state had an incorrect year.",t);
            test.BasicTest(m2,m1,"A model state had an incorrect month.",t);
            test.BasicTest(d2,d1,"A model state had an incorrect day.",t);
            curStability:=potSurf.BasicSurfaceCalculation.Stability(&val.Val);
            test.BasicTest(true,optimalStability>=curStability,
                "The optimal model state stability was not correctly found.",t,
            );
            if optimalStability==curStability {
                test.BasicTest(true,optimalMse<=val.Val.Mse,
                    "The optimal model state mse was not correctly found.",t,
                );
            }
        }
        return iter.Continue,val,nil;
    }).Filter(func(index int, val log.LogEntry[db.ModelState]) bool {
        if index==0 {
            initialMse=val.Val.Mse;
            return false;
        }
        return true;
    }).ForEach(
    func(index int, val log.LogEntry[db.ModelState]) (iter.IteratorFeedback, error) {
        newStability:=potSurf.BasicSurfaceCalculation.Stability(&val.Val);
        test.BasicTest(true,initialStability<=newStability,
            "Stability values did not continually increase.",t,
        );
        if initialStability==newStability {
            test.BasicTest(true,initialMse>val.Val.Mse,
                "Mse values did not continually decrease given equivalent stability.",t,
            );
        }
        initialStability=newStability;
        initialMse=val.Val.Mse;
        return iter.Continue,nil;
    });
    test.BasicTest(nil,err,
        "Iterating over a log generated an error when it should not have.",t,
    );
}

func generateAllModelStatesHelper(scenarioName string,
        timeFrame dataStruct.Pair[int,int],
        window dataStruct.Pair[int,int],
        numThreads int,
        t *testing.T) error {
    db.DeleteAll[db.ModelState](&testDB);
    closeLogs:=setupLogs(fmt.Sprintf(
        "./debugLogs/SlidingWindowStateGenerator.%s",scenarioName,
    ),MS_PARALLEL_RESULT_DEBUG);
    sw,_:=NewSlidingWindowStateGen(timeFrame,window,numThreads);
    c,_:=db.GetClientByEmail(&testDB,"testing@testing.com")
    // Earilest data point is 8/10/2021, this date is small enough to get all values
    sw.GenerateClientModelStates(&testDB,c,time.Date(
        2020,time.Month(1),1,0,0,0,0,time.UTC,
    ),func() []potSurf.Surface {
        return []potSurf.Surface{ potSurf.NewBasicSurface().ToGenericSurf() };
    });
    //fmt.Println(cnts,err);
    closeLogs();
    return nil;
}

func TestGenerateClientModelStatesSingleThread(t *testing.T){
    timeFrame:=dataStruct.Pair[int,int]{A: 1, B: 5000};
    window:=dataStruct.Pair[int,int]{A: 1, B: 30};
    generateAllModelStatesHelper("allValues1Thread",timeFrame,window,1,t);
    runAllModelStateDebugLogTests(t);
}

func TestGenerateClientModelStates10Threads(t *testing.T){
    timeFrame:=dataStruct.Pair[int,int]{A: 1, B: 5000};
    window:=dataStruct.Pair[int,int]{A: 1, B: 30};
    generateAllModelStatesHelper("allValues10Threads",timeFrame,window,10,t);
    threadOutputChecker(t);
    runAllModelStateDebugLogTests(t);
}

func TestGenerateClientModelStates100Threads(t *testing.T){
    timeFrame:=dataStruct.Pair[int,int]{A: 1, B: 5000};
    window:=dataStruct.Pair[int,int]{A: 1, B: 30};
    generateAllModelStatesHelper("allValues100Threads",timeFrame,window,100,t);
    threadOutputChecker(t);
    runAllModelStateDebugLogTests(t);
}

// Yes, this is an n^2 algorithm. It's just for testing. There is no way to index
// or sort the values in the log because the order is not guaranteed by either the 
// sql query or the multi-threaded logs.
func threadOutputChecker(t *testing.T){
    log.LogElems(SLIDING_WINDOW_MS_PARALLEL_RESULT_DEBUG).ForEach(
    func(index int, val log.LogEntry[db.ModelState]) (iter.IteratorFeedback, error) {
        l,_:=log.NewLog[db.ModelState](log.Debug,
            "./debugLogs/SlidingWindowStateGenerator.allValues1Thread.modelState.log",
            true,
        );
        _,_,found:=log.LogElems(l).Find(func(val2 log.LogEntry[db.ModelState]) (bool, error) {
            return (val.Val.Date==val2.Val.Date &&
                val.Val.ClientID==val2.Val.ClientID &&
                val.Val.Mse==val2.Val.Mse),nil;
        });
        test.BasicTest(true,found,
            "The value from the single thread was not found in the multi-thread test.",t,
        );
        return iter.Continue,nil;
    });
}

func runAllModelStateDebugLogTests(t *testing.T){
    type tmp struct { V int; };
    ids,_:=db.CustomReadQuery[tmp](&testDB,`SELECT Exercise.Id 
        FROM Exercise 
        JOIN ExerciseType 
        ON Exercise.TypeID=ExerciseType.Id
        WHERE ExerciseType.T='Main Compound' 
            OR ExerciseType.T='Main Compound Accessory';`,
        []any{},
    ).Collect();
    log.LogElems(SLIDING_WINDOW_MS_PARALLEL_RESULT_DEBUG).ForEach(
    func(index int, val log.LogEntry[db.ModelState]) (iter.IteratorFeedback, error) {
        _,_,found:=iter.SliceElems[*tmp](ids).Find(func(idVal *tmp) (bool, error) {
            return val.Val.ExerciseID==idVal.V,nil;
        });
        test.BasicTest(true,found,"ExerciseID was not found in valid ID list.",t);
        return iter.Continue,nil;
    });
}
