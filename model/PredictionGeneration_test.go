package model;

import (
    "time"
    "testing"
    "database/sql"
    "github.com/barbell-math/block/db"
    "github.com/barbell-math/block/util/test"
)

const TEST_SG_NAME string="TestSG";

func createPredictionData(createMs bool, createTl bool) (func ()){
    sgId,msId,tlId:=[]int{}, []int{}, []int{};
    sgId,_=db.Create(&testDB,db.StateGenerator{T: TEST_SG_NAME});
    if createMs {
        msId,_=db.Create(&testDB,db.ModelState{
            ClientID: 1, ExerciseID: 15, StateGeneratorID: sgId[0],
            Date: time.Now().AddDate(0, 0, -1),
            Eps: 0, Eps1: 0, Eps2: 0, Eps3: 0, Eps4: 0,
            Eps5: 0, Eps6: 0, Eps7: 0,
        });
    }
    if createTl {
        tlId,_=db.Create(&testDB,db.TrainingLog{
            ClientID: 1, ExerciseID: 15, RotationID: 1,
            DatePerformed: time.Now().AddDate(0, 0, -1),
        });
    }
    return func(){
        db.Delete[db.StateGenerator](
            &testDB,
            db.StateGenerator{ Id: sgId[0] },
            db.OnlyIDFilter,
        );
        if createMs {
            db.Delete[db.ModelState](
                &testDB,
                db.ModelState{ Id: msId[0] },
                db.OnlyIDFilter,
            );
        }
        if createTl {
            db.Delete[db.TrainingLog](
                &testDB,
                db.TrainingLog{ Id: tlId[0] },
                db.OnlyIDFilter,
            );
        }
    }
}

func TestGeneratePrediction(t *testing.T){
    defer createPredictionData(true,true)();
    tl:=db.TrainingLog{
        ClientID: 1,
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0,
        InterWorkoutFatigue: 0, InterExerciseFatigue: 0,
        ExerciseID: 15, DatePerformed: time.Now(),
    };
    sg,_:=db.GetStateGeneratorByName(&testDB,TEST_SG_NAME);
    _,err:=GeneratePrediction(&testDB,&tl,&sg);
    test.BasicTest(nil,err,
        "Generate prediction returned an error when it was not supposed to.",t,
    );
}

func TestGeneratePredictionNoStateGenerator(t *testing.T){
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0,
        InterWorkoutFatigue: 0, InterExerciseFatigue: 0,
        DatePerformed: time.Now(),
    };
    sg:=db.StateGenerator{ T: "Non-existent state generator" };
    _,err:=GeneratePrediction(&testDB,&tl,&sg);
    test.BasicTest(sql.ErrNoRows,err,
        "The wrong error was raised when no predictions where found.",t,
    );
}

func TestGeneratePredictionNoModelState(t *testing.T){
    defer createPredictionData(false,true)();
    tl:=db.TrainingLog{
        ClientID: 1,
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0,
        InterWorkoutFatigue: 0, InterExerciseFatigue: 0,
        ExerciseID: 15, DatePerformed: time.Now(),
    };
    sg,_:=db.GetStateGeneratorByName(&testDB,TEST_SG_NAME);
    _,err:=GeneratePrediction(&testDB,&tl,&sg);
    test.BasicTest(sql.ErrNoRows,err,
        "Generate prediction returned incorrect error.",t,
    );
}
