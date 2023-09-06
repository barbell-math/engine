package bookDataGen;

import (
	"fmt"
	"testing"

	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/settings"
	"github.com/barbell-math/block/util/algo/iter"
	customerr "github.com/barbell-math/block/util/err"
	"github.com/barbell-math/block/util/io/csv"
)

var testDB db.DB;

func TestMain(m *testing.M){
    settings.ReadSettings("testData/bookTestSettings.json");
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
                "../../data/testData/ClientTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.Client) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
        }, func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.StateGenerator](csv.CSVFileSplitter(
                "../../data/testData/StateGeneratorTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.StateGenerator) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
        }, func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.ExerciseType](csv.CSVFileSplitter(
                "../../data/testData/ExerciseTypeTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.ExerciseType) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
        },func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.ExerciseFocus](csv.CSVFileSplitter(
                "../../data/testData/ExerciseFocusTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.ExerciseFocus) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
        },func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.Exercise](csv.CSVFileSplitter(
                "../../data/testData/ExerciseTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.Exercise) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
        },func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.Rotation](csv.CSVFileSplitter(
                "../../data/testData/RotationTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.Rotation) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
        },func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.TrainingLog](csv.CSVFileSplitter(
                "../../data/testData/AugmentedTrainingLogTestData.csv",',','#',
            ),"1/2/2006").ForEach(
            func(index int, val db.TrainingLog) (iter.IteratorFeedback, error) {
                db.Create(&testDB,val);
                return iter.Continue,nil;
            });
    });
}

func teardown(){
    testDB.ResetDB();
    testDB.Close();
}
