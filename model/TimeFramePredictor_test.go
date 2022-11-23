package model;

func TestNewTimeFramePredictor(t *testing.T){
    tmp,err:=NewTimeFramePredictor(1,2,3);
    testUtil.BasicTest(nil,err,
        "An error was raised when it shouldn't have been.",t,
    );
    testUtil.BasicTest(-1,tmp.minTimeFrame,"Min time frame not set correctly.",t);
    testUtil.BasicTest(-2,tmp.maxTimeFrame,"Max time frame not set correctly.",t);
    testUtil.BasicTest(-3,tmp.window,"Window not set correctly.",t);
    tmp,err=NewTimeFramePredictor(2,1,3);
    if !util.IsInvalidPredictionState(err) {
        testUtil.FormatError(util.InvalidPredictionState(""),err,
            "An error was not raised when it should have been.",t,
        );
    }
    testUtil.BasicTest(-2,tmp.minTimeFrame,"Min time frame not set correctly.",t);
    testUtil.BasicTest(-1,tmp.maxTimeFrame,"Max time frame not set correctly.",t);
    testUtil.BasicTest(-3,tmp.window,"Window not set correctly.",t);
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
//            go NewPredictionState(38,500,0).GenerateModelState(&testDB,*t,ch);
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

func TestTimeFramePredictor(t *testing.T){
    setup();
    ch:=make(chan []error);
    go NewPredictionState(38,500,13).UpdateModelStates(&testDB,1,ch);
    errs:=<-ch;
    fmt.Println(errs);
}
