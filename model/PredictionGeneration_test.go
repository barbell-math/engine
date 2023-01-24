package model;

import (
    //"fmt"
    "time"
    "testing"
    "database/sql"
    "github.com/barbell-math/block/db"
    "github.com/barbell-math/block/util/test"
)

//func TestGeneratePrediction(t *testing.T){
//    msId,e1:=db.Create(&testDB,db.ModelState{
//        ClientID: 1, ExerciseID: 15, StateGeneratorID: 1,
//        Date: time.Now().AddDate(0, 0, -1),
//        A: 0, B: 0, C: 0, D: 0, Eps: 0, Eps2: 0,
//    });
//    fmt.Println(e1);
//    tl:=db.TrainingLog{
//        ClientID: 1,
//        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0, FatigueIndex: 0,
//        ExerciseID: 15, DatePerformed: time.Now(),
//    };
//    sg:=db.StateGenerator{ Id: 1, T: "Sliding Window" };
//    p,err:=GeneratePrediction(&testDB,&tl,&sg);
//    fmt.Println(err,p);
//    db.Delete[db.ModelState](&testDB,db.ModelState{ Id: msId[0] },db.OnlyIDFilter);
//}

func TestGeneratePredictionNoState(t *testing.T){
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0,
        InterWorkoutFatigue: 0, InterExerciseFatigue: 0,
        DatePerformed: time.Now(),
    };
    sg:=db.StateGenerator{ T: "Non-existant state generator" };
    _,err:=GeneratePrediction(&testDB,&tl,&sg);
    test.BasicTest(sql.ErrNoRows,err,
        "The wrong error was raised when no predictions where found.",t,
    );
}
