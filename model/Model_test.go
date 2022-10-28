package model;

import (
    "fmt"
    "time"
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
    if err:=uploadTestData(); err!=nil {
        panic(fmt.Sprintf(
            "Could not upload data for testing. Check location of testData folder. | %s",
            err,
        ));
    }
}

func uploadTestData() error {
    return util.ChainedErrorOps(
        func(r ...any) (any,error) {
            return db.Create(&testDB,db.Client{
                Id: 1,
                FirstName: "testF",
                LastName: "testL",
                Email: "test@test.com",
            });
        }, func(r ...any) (any,error) {
            return nil,util.CSVToStruct(
                "../testData/model/ExerciseTypeTestData.csv",',',"",
                func(e *db.ExerciseType){
                    //fmt.Println(*e);
                    db.Create(&testDB,*e);
            });
        },func(r ...any) (any,error) {
            return nil,util.CSVToStruct(
                "../testData/model/ExerciseFocusTestData.csv",',',"",
                func(e *db.ExerciseFocus){
                    //fmt.Println(*e);
                    db.Create(&testDB,*e);
            });
        },func(r ...any) (any,error) {
            return nil,util.CSVToStruct(
                "../testData/model/ExerciseTestData.csv",',',"",
                func(e *db.Exercise){
                    //fmt.Println(*e);
                    db.Create(&testDB,*e);
            });
        },func(r ...any) (any,error) {
            return nil,util.CSVToStruct(
                "../testData/model/RotationTestData.csv",',',"1/2/2006",
                func(r *db.Rotation){
                    //fmt.Println(*r);
                    db.Create(&testDB,*r);
            });
        },func(r ...any) (any,error) {
            return nil,util.CSVToStruct(
                "../testData/model/TrainingLogTestData.csv",',',"1/2/2006",
                func(t *db.TrainingLog){
                    //fmt.Println(*t);
                    db.Create(&testDB,*t);
            });
    });
}

func teardown(){
    testDB.ResetDB();
    testDB.Close();
}

func TestPlaceHolder(t *testing.T){
    setup();
    var tmp db.TrainingLog;
    latestSquat,_:=time.Parse("1/2/2006","7/24/2022");
    err:=db.Read(&testDB,db.TrainingLog{
        DatePerformed: latestSquat,
        ExerciseID: 42,
    },util.GenFilter(false,"DatePerformed","ExerciseID"),func(t *db.TrainingLog){
        tmp=*t;
    });
    fmt.Println(err);
    fmt.Println(tmp);
    GenerateModelState(&testDB,&tmp);
}
