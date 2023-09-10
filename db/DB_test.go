package db;

import (
    "fmt"
    "testing"
    "github.com/barbell-math/engine/settings"
    "github.com/barbell-math/engine/util/test"
)

var testDB DB;

func TestMain(m *testing.M){
    //setup();
    settings.ReadSettings("testData/dbTestSettings.json");
    m.Run();
    teardown();
}

func setup(){
    var err error=nil;
    testDB,err=NewDB(settings.DBHost(),settings.DBPort(),settings.DBName());
    if err!=nil && err!=DataVersionNotAvailable {
        panic("Could not open database for testing.");
    }
    if err=testDB.ResetDB(); err!=nil {
        panic(fmt.Sprintf("Could not reset DB for testing. Check location of global init SQL file relative to the ./testData/dbTestSettings.json file. \n  | Given err: %v",err));
    }
}

func teardown(){
    testDB.ResetDB();
    testDB.Close();
}

func createExerciseTestData(){
    Create(&testDB,ExerciseFocus{Focus: "Squat"});
    Create(&testDB,ExerciseType{T: "Accessory"});
    Create(&testDB,
        Exercise{Name: "Squat", FocusID: 1, TypeID: 1},
        Exercise{Name: "Bench", FocusID: 1, TypeID: 1},
        Exercise{Name: "Deadlift", FocusID: 1, TypeID: 1},
    );
}

func TestVersion(t *testing.T){
    setup();
    //val,err:=testDB.getDataVersion();
    //test.BasicTest(
    //    sql.ErrNoRows,err,
    //    "Attempting to get version before it was added was successful.",t,
    //);
    err:=testDB.setDataVersion(-1);
    test.BasicTest(nil,err,"Could not set data version.",t);
    val,err:=testDB.getDataVersion();
    test.BasicTest(nil,err,"Could not access version.",t);
    test.BasicTest(-1,val,"Version was not set correctly.",t);
    err=testDB.setDataVersion(-2);
    test.BasicTest(nil,err,"Could not set data version.",t);
    val,err=testDB.getDataVersion();
    test.BasicTest(nil,err,"Could not access version.",t);
    test.BasicTest(-2,val,"Version was not set correctly.",t);
    err=testDB.db.QueryRow("SELECT COUNT(*) FROM Version;").Scan(&val);
    test.BasicTest(nil,err,"Could not access version for counting.",t);
    test.BasicTest(1,val,"Version table was not limited to only one entry.",t);
}
