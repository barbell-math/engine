package model;

import (
    "fmt"
    "time"
    "testing"
    "database/sql"
    "github.com/barbell-math/block/db"
    "github.com/barbell-math/block/util/test"
)

func TestGeneratePrediction(t *testing.T){
    sgId,e:=db.Create(&testDB,db.StateGenerator{T: "TestSG"});
    msId,e1:=db.Create(&testDB,db.ModelState{
        ClientID: 1, ExerciseID: 15, StateGeneratorID: sgId[0],
        Date: time.Now().AddDate(0, 0, -1),
        Eps: 0, Eps1: 0, Eps2: 0, Eps3: 0, Eps4: 0,
        Eps5: 0, Eps6: 0, Eps7: 0,
    });
    fmt.Println(e);
    fmt.Println(e1);
    fmt.Println(sgId[0]);
    fmt.Println(msId[0]);
    tl:=db.TrainingLog{
        ClientID: 1,
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0,
        InterWorkoutFatigue: 0, InterExerciseFatigue: 0,
        ExerciseID: 15, DatePerformed: time.Now(),
    };
    sg,_:=db.GetStateGeneratorByName(&testDB,"TestSG");
    fmt.Println(sg);
    p,err:=GeneratePrediction(&testDB,&tl,&sg);
    fmt.Println(err,p);
    db.Delete[db.ModelState](&testDB,db.ModelState{ Id: msId[0] },db.OnlyIDFilter);
    db.Delete[db.StateGenerator](
        &testDB,
        db.StateGenerator{ Id: sgId[0] },
        db.OnlyIDFilter,
    );
}

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
