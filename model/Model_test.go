package model;

import (
    "fmt"
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/db"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/settings"
)

var testDB db.CRUD;

func TestMain(m *testing.M){
    //setup();
    settings.ReadSettings("../testData/modelTestSettings.json");
    m.Run();
    teardown();
}

func setup(){
    var err error=nil;
    testDB,err=db.NewCRUD(settings.DBHost(),settings.DBPort(),settings.DBName());
    if err!=nil && err!=util.DataVersionNotAvailable {
        panic("Could not open database for testing.");
    }
    if err=testDB.ResetDB(); err!=nil {
        panic("Could not reset DB for testing. Check location of global init SQL file relative to the ./testData/modelTestSettings.json file.");
    }
    if err=testDB.ExecSQLScript("../sql/uploadModelTestData.sql"); err!=nil {
        fmt.Println(err);
        panic("Could not upload test data to run tests on the model. Check the location of the 'uploadModelTestData.sql' file relative to the ./sql folder.");
    }
}

func teardown(){
    testDB.ResetDB();
    testDB.Close();
}


func TestPlaceHolder(t *testing.T){
    setup();
}
