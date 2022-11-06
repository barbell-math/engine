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
    if err:=uploadTestData(); err!=nil {
        panic(fmt.Sprintf(
            "Could not upload data for testing. Check location of testData folder. | %s",
            err,
        ));
    }
    //addFatigueIndex();
}

//This function can be used to generate the augmented training log file
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
                "../testData/model/AugmentedTrainingLogTestData.csv",',',"1/2/2006",
                func(t *db.TrainingLog){
                    //fmt.Println(*t);
                    db.Create(&testDB,*t);
            });
    });
}

func teardown(){
    //testDB.ResetDB();
    testDB.Close();
}

//func TestPlaceHolder(t *testing.T){
//    //setup();
//    var tmp db.TrainingLog;
//    latestSquat,_:=time.Parse("1/2/2006","7/24/2022");
//    err:=db.Read(&testDB,db.TrainingLog{
//        DatePerformed: latestSquat,
//        ExerciseID: 14,
//    },util.GenFilter(false,"DatePerformed","ExerciseID"),func(t *db.TrainingLog){
//        tmp=*t;
//    });
//    fmt.Println(err);
//    fmt.Println(tmp);
//    ms,err:=GenerateModelState(&testDB,&tmp);
//    fmt.Printf("%+v\n",ms);
//    fmt.Println(err);
//}

func TestPlaceHolder2(t *testing.T){
    var avg float64=0;
    var cntr int=0;
    db.Read(&testDB,db.TrainingLog{
        ExerciseID: 14,
    },util.GenFilter(false,"ExerciseID"),func(t *db.TrainingLog){
        modelState,_:=GenerateModelState(&testDB,t);
        //fmt.Printf("%+v\n",modelState);
        if modelState.Difference>0 {
            avg+=modelState.Difference*modelState.Difference;
            cntr++;
        }
        db.Create(&testDB,modelState);
    });
    //fmt.Println("Avg diff: ",avg/float64(cntr));
}

func TestPrediction(t *testing.T){
    for i:=0; i<10; i++ {
        setup();
        db.Read(&testDB,db.TrainingLog{
            ExerciseID: 14,
        },util.GenFilter(false,"ExerciseID"),func(t *db.TrainingLog){
            modelState,_:=GenerateModelState(&testDB,t);
            //fmt.Printf("%+v\n",modelState);
            db.Create(&testDB,modelState);
        });
        db.Read(&testDB,db.TrainingLog{
            ExerciseID: 14,
        },util.GenFilter(false,"ExerciseID"),func(tl *db.TrainingLog){
            p,err:=GeneratePrediction(&testDB,tl);
            //fmt.Println(err);
            //fmt.Printf("%+v\n",tl);
            if err==nil {
                //fmt.Printf("%+v\n",p);
                db.Create(&testDB,p);
            }
        });
        type ErrResult struct { Mse float64; };
        query:=`SELECT AVG(
            POWER(TrainingLog.Intensity-Prediction.IntensityPred,2)
        ) FROM TrainingLog
        JOIN Prediction ON Prediction.TrainingLogID=TrainingLog.Id
        WHERE Prediction.IntensityPred>0;`
        db.CustomReadQuery(&testDB,query,[]any{},func(r *ErrResult){
            fmt.Printf("%+v\n",r);
        });
    }
}
