package db;

import (
    "database/sql"
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
)

var testDB CRUD;

func TestMain(m *testing.M){
    setup();
    m.Run();
    teardown();
}

func setup(){
    var err error=nil;
    testDB,err=NewCRUD("localhost",5432,"carmichaeljr","test");
    if err!=nil && err!=util.DataVersionNotAvailable {
        panic("Could not open database for testing.");
    }
    err=testDB.execSQLScript("../sql/globalInit.sql");
    if util.IsSqlScriptNotFound(err) {
        panic("Could not find 'globalInit.sql' file for testing.");
    } else if err!=nil {
        panic("An error occurred running the 'globalInit.sql' script.");
    }
}

func teardown(){
    testDB.execSQLScript("../sql/globalInit.sql");
    testDB.Close();
}

func TestVersion(t *testing.T){
    val,err:=testDB.getDataVersion();
    testUtil.BasicTest(
        sql.ErrNoRows,err,
        "Attempting to get version before it was added was successful.",t,
    );
    err=testDB.setDataVersion(-1);
    testUtil.BasicTest(nil,err,"Could not set data version.",t);
    val,err=testDB.getDataVersion();
    testUtil.BasicTest(nil,err,"Could not access version.",t);
    testUtil.BasicTest(-1,val,"Version was not set correctly.",t);
    err=testDB.setDataVersion(-2);
    testUtil.BasicTest(nil,err,"Could not set data version.",t);
    val,err=testDB.getDataVersion();
    testUtil.BasicTest(nil,err,"Could not access version.",t);
    testUtil.BasicTest(-2,val,"Version was not set correctly.",t);
    err=testDB.db.QueryRow("SELECT COUNT(*) FROM Version;").Scan(&val);
    testUtil.BasicTest(nil,err,"Could not access version for counting.",t);
    testUtil.BasicTest(1,val,"Version table was not limited to only one entry.",t);
}

func TestCreateExerciseType(t *testing.T){
    var id int=0;
    _,err:=testDB.ReadExerciseType(1);
    testUtil.BasicTest(
        sql.ErrNoRows,err,
        "Reading exercise type before any were added was successful.",t,
    );
    id,err=testDB.CreateExerciseType(
        ExerciseType{
            _type: "TestType",
            description: "TestTypeDescription",
        },
    );
    testUtil.BasicTest(nil,err,"Could not create exercise type.",t);
    testUtil.BasicTest(1 ,id,"Exercise was not created correctly.",t);
    id,err=testDB.CreateExerciseType(
        ExerciseType{
            _type: "TestType1",
            description: "TestTypeDescription1",
        },
    );
    testUtil.BasicTest(nil,err,"Could not create exercise type.",t);
    testUtil.BasicTest(2 ,id,"Exercise was not created correctly.",t);
    val,err:=testDB.ReadExerciseType(1);
    testUtil.BasicTest(
        nil,err,
        "Could not read exercise type that was previously inserted.",t,
    );
}
