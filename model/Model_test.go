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

//func TestPlaceHolder2(t *testing.T){
//    var avg float64=0;
//    var cntr int=0;
//    db.Read(&testDB,db.TrainingLog{
//        ExerciseID: 14,
//    },util.GenFilter(false,"ExerciseID"),func(t *db.TrainingLog){
//        modelState,_:=GenerateModelState(&testDB,t);
//        //fmt.Printf("%+v\n",modelState);
//        if modelState.Difference>0 {
//            avg+=modelState.Difference*modelState.Difference;
//            cntr++;
//        }
//        db.Create(&testDB,modelState);
//    });
//    //fmt.Println("Avg diff: ",avg/float64(cntr));
//}

//func TestPrediction(t *testing.T){
//    ch:=make(chan ModelStateGenerationRes);
//    cntr:=0;
//    for i:=0; i<1; i++ {
//        setup();
//        //t,_:=time.Parse("1/2/2006","7/14/2022");
//        fmt.Println("UPDATING MODEL STATE ===================");
//        db.Read(&testDB,db.TrainingLog{
//            ExerciseID: 14,
//            //DatePerformed: t,
//        },util.GenFilter(false,"ExerciseID"),func(t *db.TrainingLog){
//            go NewPredictionState(38,500,13).GenerateModelState(&testDB,*t,ch);
//            cntr++;
//            //if res:=<-ch; res.Err==nil {
//            //    //fmt.Printf("%+v\n",res.Ms);
//            //    db.Create(&testDB,res.Ms);
//            //}
//        });
//        for i:=0; i<cntr; i++ {
//            if res:=<-ch; res.Err==nil {
//                db.Create(&testDB,res.Ms);
//            }
//        }
//        //t,_=time.Parse("1/2/2006","7/24/2022");
//        fmt.Println("UPDATING MODEL PRED ===================");
//        err:=db.Read(&testDB,db.TrainingLog{
//            ExerciseID: 14,
//            //DatePerformed: t,
//        },util.GenFilter(false,"ExerciseID"),func(tl *db.TrainingLog){
//            p,err:=GeneratePrediction(&testDB,tl);
//            //fmt.Println(err);
//            //fmt.Printf("%+v\n",tl);
//            if err==nil {
//                //fmt.Printf("%+v\n",p);
//                db.Create(&testDB,p);
//            }
//        });
//        fmt.Println(err);
//        type ErrResult struct { Mse float64; };
//        query:=`SELECT SQRT(AVG(
//            POWER(TrainingLog.Intensity-Prediction.IntensityPred,2)
//        )) FROM TrainingLog
//        JOIN Prediction ON Prediction.TrainingLogID=TrainingLog.Id
//        WHERE Prediction.IntensityPred>0;`
//        db.CustomReadQuery(&testDB,query,[]any{},func(r *ErrResult){
//            fmt.Printf("%+v\n",r);
//        });
//    }
//}

func TestPlaceholder(t *testing.T){
    setup();
    ch:=make(chan error);
    go NewPredictionState(38,500,13).UpdateModelStates(&testDB,1,ch);
    err:=<-ch;
    fmt.Println(err);
}
