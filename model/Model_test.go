package model;

import (
    "os"
    "fmt"
    "time"
    "errors"
    "strconv"
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/db"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/settings"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
)

var testDB db.CRUD;

func TestMain(m *testing.M){
    settings.ReadSettings("../testData/modelTestSettings.json");
    setup();
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
    //Enable this code as needed, it will re-generate the augmented model data
    //if addFatigueIndex(); err!=nil {
    //    panic(fmt.Sprintf(
    //        "Could not generate the augmented model data | %s",err,
    //    ));
    //}
    if err:=uploadTestData(); err!=nil {
        panic(fmt.Sprintf(
            "Could not upload data for testing. Check location of testData folder. | %s",
            err,
        ));
    }
}

func addFatigueIndex() error {
    f,err:=os.Create("../testData/model/AugmentedTrainingLogTestData.csv");
    if err!=nil {
        return errors.New("Err occured opening augmented training log data file.");
    }
    cntr:=0;
    var curDate time.Time;
    dayVals:=make(map[string]int,0);
    return util.CSVFileSplitter("../testData/model/TrainingLogTestData.csv",
        ',',false,func(cols []string) bool {
            if cntr==0 {
                f.WriteString(util.CSVGenerator(",",func(iter int) (string,bool){
                    if iter==len(cols) {
                        return "FatigueIndex",false;
                    }
                    return cols[iter],true;
                }));
            } else {
                iterDate,_:=time.Parse("1/2/2006",cols[4]);
                if !curDate.Equal(iterDate) {
                    dayVals=make(map[string]int,0);
                    curDate=iterDate;
                }
                if val,ok:=dayVals[cols[2]]; ok {
                    dayVals[cols[2]]=val+1;
                } else {
                    dayVals[cols[2]]=0;
                }
                f.WriteString(util.CSVGenerator(",",func(iter int) (string,bool){
                    if iter==len(cols) {
                        return strconv.Itoa(dayVals[cols[2]]),false;
                    }
                    return cols[iter],true;
                }));
            }
            f.WriteString("\n");
            cntr++;
            return true;
    });
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
                    db.Create(&testDB,*e);
            });
        },func(r ...any) (any,error) {
            return nil,util.CSVToStruct(
                "../testData/model/ExerciseFocusTestData.csv",',',"",
                func(e *db.ExerciseFocus){
                    db.Create(&testDB,*e);
            });
        },func(r ...any) (any,error) {
            return nil,util.CSVToStruct(
                "../testData/model/ExerciseTestData.csv",',',"",
                func(e *db.Exercise){
                    db.Create(&testDB,*e);
            });
        },func(r ...any) (any,error) {
            return nil,util.CSVToStruct(
                "../testData/model/RotationTestData.csv",',',"1/2/2006",
                func(r *db.Rotation){
                    db.Create(&testDB,*r);
            });
        },func(r ...any) (any,error) {
            return nil,util.CSVToStruct(
                "../testData/model/AugmentedTrainingLogTestData.csv",',',"1/2/2006",
                func(t *db.TrainingLog){
                    db.Create(&testDB,*t);
            });
    });
}

func teardown(){
    //testDB.ResetDB();
    testDB.Close();
}

func TestModelCreation(t *testing.T){
    cntr:=0;
    m:=fatigueAwareModel();
    m.IterLHS(func(r int, c int, v float64){
        cntr++;
    });
    testUtil.BasicTest(36,cntr,"LHS Lin reg wrong size for model.",t);
    cntr=0;
    m.IterRHS(func(r int, c int, v float64){
        cntr++;
    });
    testUtil.BasicTest(6,cntr,"RHS Lin reg wrong size for model.",t);
}
