package stateGenerator

import (
	"fmt"
	"testing"

	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/settings"
	"github.com/barbell-math/block/util/algo/iter"
	customerr "github.com/barbell-math/block/util/err"
	"github.com/barbell-math/block/util/io/csv"
	logUtil "github.com/barbell-math/block/util/io/log"
)

var testDB db.DB;
const (
    DP_DEBUG int=1<<iota
    MS_DEBUG
    MS_PARALLEL_RESULT_DEBUG
)

func setupLog[T any](dest *logUtil.Logger[T], path string){
    customerr.PanicOnError(func() error {
        var err error;
        *dest,err=logUtil.NewLog[T](logUtil.Debug,path,false);
        return err;
    });
}
func setupLogs(debugFile string, logFlags int) (func()) {
    if DP_DEBUG&logFlags==DP_DEBUG {
        setupLog(&SLIDING_WINDOW_DP_DEBUG,
            fmt.Sprintf("%s.dataPoint.log",debugFile),
        );
    }
    if MS_DEBUG&logFlags==MS_DEBUG {
        setupLog(&SLIDING_WINDOW_MS_DEBUG,
            fmt.Sprintf("%s.modelState.log",debugFile),
        );
    }
    if MS_PARALLEL_RESULT_DEBUG&logFlags==MS_PARALLEL_RESULT_DEBUG {
        setupLog(&SLIDING_WINDOW_MS_PARALLEL_RESULT_DEBUG,
            fmt.Sprintf("%s.modelState.log",debugFile),
        );
    }
    return func(){
        if DP_DEBUG&logFlags==DP_DEBUG {
            SLIDING_WINDOW_DP_DEBUG.Close();
        }
        if MS_DEBUG&logFlags==MS_DEBUG {
            SLIDING_WINDOW_MS_DEBUG.Close();
        }
        if MS_PARALLEL_RESULT_DEBUG&logFlags==MS_PARALLEL_RESULT_DEBUG {
            SLIDING_WINDOW_MS_PARALLEL_RESULT_DEBUG.Close();
        }
    }
}

func TestMain(m *testing.M){
    settings.ReadSettings("testData/modelTestSettings.json");
    setup();
    m.Run();
    teardown();
}

func setup(){
    var err error=nil;
    testDB,err=db.NewDB(settings.DBHost(),settings.DBPort(),settings.DBName());
    if err!=nil && err!=db.DataVersionNotAvailable {
        panic("Could not open database for testing.");
    }
    if err=testDB.ResetDB(); err!=nil {
        panic("Could not reset DB for testing. Check location of global init SQL file relative to the ./testData/modelTestSettings.json file.");
    }
    if err:=uploadTestData(); err!=nil {
        panic(fmt.Sprintf(
            "Could not upload data for testing. Check location of testData folder. | %s",
            err,
        ));
    }
}

func uploadTestData() error {
    return customerr.ChainedErrorOps(
        func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.Client](csv.CSVFileSplitter(
                "../../../data/testData/ClientTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.Client) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
        }, func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.StateGenerator](csv.CSVFileSplitter(
                "../../../data/testData/StateGeneratorTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.StateGenerator) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
        }, func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.ExerciseType](csv.CSVFileSplitter(
                "../../../data/testData/ExerciseTypeTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.ExerciseType) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
        },func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.ExerciseFocus](csv.CSVFileSplitter(
                "../../../data/testData/ExerciseFocusTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.ExerciseFocus) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
        },func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.Exercise](csv.CSVFileSplitter(
                "../../../data/testData/ExerciseTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.Exercise) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
        },func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.Rotation](csv.CSVFileSplitter(
                "../../../data/testData/RotationTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.Rotation) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
        },func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.TrainingLog](csv.CSVFileSplitter(
                "../../../data/testData/AugmentedTrainingLogTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.TrainingLog) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
    });
}

func teardown(){
    //testDB.ResetDB();
    testDB.Close();
}
