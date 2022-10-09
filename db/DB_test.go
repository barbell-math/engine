package db;

import (
    "fmt"
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
    "github.com/carmichaeljr/powerlifting-engine/settings"
)

var testDB CRUD;

func TestMain(m *testing.M){
    //setup();
    settings.ReadSettings("../testData/dbTestSettings.json");
    m.Run();
    teardown();
}

func setup(){
    var err error=nil;
    testDB,err=NewCRUD(settings.DBHost(),settings.DBPort(),settings.DBName());
    if err!=nil && err!=util.DataVersionNotAvailable {
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

func TestVersion(t *testing.T){
    setup();
    //val,err:=testDB.getDataVersion();
    //testUtil.BasicTest(
    //    sql.ErrNoRows,err,
    //    "Attempting to get version before it was added was successful.",t,
    //);
    err:=testDB.setDataVersion(-1);
    testUtil.BasicTest(nil,err,"Could not set data version.",t);
    val,err:=testDB.getDataVersion();
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
