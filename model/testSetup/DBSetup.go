package testSetup

import (
    "fmt"

    "github.com/barbell-math/engine/db"
    "github.com/barbell-math/engine/settings"
    "github.com/barbell-math/engine/util/algo/iter"
    "github.com/barbell-math/engine/util/io/csv"
    customerr "github.com/barbell-math/engine/util/err"
)


func SetupDB() db.DB {
    var testDB db.DB;
    var err error=nil;
    fmt.Print("Connecting to db...");
    testDB,err=db.NewDB(settings.DBHost(),settings.DBPort(),settings.DBName());
    if err!=nil && err!=db.DataVersionNotAvailable {
        panic("Could not open database for testing.");
    } else {
        fmt.Println("Done");
    }
    fmt.Print("Resetting db...");
    if err=testDB.ResetDB(); err!=nil {
        panic("Could not reset DB for testing. Check location of global init SQL file relative to the ./testData/testSettings.json file.");
    } else {
        fmt.Println("Done");
    }
    fmt.Print("Uploading test data...\r");
    if err:=uploadTestData(&testDB,"Uploading test data..."); err!=nil {
        panic(fmt.Sprintf(
            "Could not upload data for testing. Check location of testData folder. | %s",
            err,
        ));
    } else {
        fmt.Print("Uploading test data...Done");
        for i:=0; i<80; i++ { fmt.Print(" "); }
        fmt.Println();
    }
    return testDB;
}

func uploadFunc[R db.DBTable](testDB *db.DB, progressLineHeader string, fName string) (
    func(index int, val R) (iter.IteratorFeedback,error),
) {
    return func(index int, val R) (iter.IteratorFeedback, error) {
        _,err:=db.Create(testDB,val);
        fmt.Printf(
            "%s (File: %s, Line: %d)\r",
            progressLineHeader,
            fName,
            index+1,
        );
        return iter.Continue,err;
    }
}

func uploadTestData(testDB *db.DB, progressLineHeader string) error {
    return customerr.ChainedErrorOps(
        func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.Client](csv.CSVFileSplitter(
                settings.ClientInitData(),',','#',
            ),"1/2/2006").ForEach(uploadFunc[db.Client](
                testDB,progressLineHeader,"ClientTestData",
            ));
        }, func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.StateGenerator](csv.CSVFileSplitter(
                settings.StateGeneratorInitData(),',','#',
            ),"1/2/2006").ForEach(uploadFunc[db.StateGenerator](
                testDB,progressLineHeader,"StateGeneratorTestData",
            ));
        }, func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.PotentialSurface](csv.CSVFileSplitter(
                settings.PotentialSurfaceInitData(),',','#',
            ),"1/2/2006").ForEach(uploadFunc[db.PotentialSurface](
                testDB,progressLineHeader,"PotentialSurfaceTestData",
            ));
        }, func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.ExerciseType](csv.CSVFileSplitter(
                settings.ExerciseTypeInitData(),',','#',
            ),"1/2/2006").ForEach(uploadFunc[db.ExerciseType](
                testDB,progressLineHeader,"ExerciseTypeTestData",
            ));
        },func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.ExerciseFocus](csv.CSVFileSplitter(
                settings.ExerciseFocusInitData(),',','#',
            ),"1/2/2006").ForEach(uploadFunc[db.ExerciseFocus](
                testDB,progressLineHeader,"ExerciseFocusTestData",
            ));
        },func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.Exercise](csv.CSVFileSplitter(
                settings.ExerciseInitData(),',','#',
            ),"1/2/2006").ForEach(uploadFunc[db.Exercise](
                testDB,progressLineHeader,"ExerciseTestData",
            ));
        },func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.Rotation](csv.CSVFileSplitter(
                settings.RotationInitData(),',','#',
            ),"1/2/2006").ForEach(uploadFunc[db.Rotation](
                testDB,progressLineHeader,"RotationTestData",
            ));
        },func(r ...any) (any,error) {
            return nil,csv.CSVToStruct[db.TrainingLog](csv.CSVFileSplitter(
                settings.TrainingLogInitData(),',','#',
            ),"1/2/2006").ForEach(uploadFunc[db.TrainingLog](
                testDB,progressLineHeader,"AugmentedTrainingLogTestData",
            ));
    });
}

func TeardownDB(testDB *db.DB){
    //testDB.ResetDB();
    testDB.Close();
}
