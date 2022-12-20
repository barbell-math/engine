package model;

import (
    "fmt"
    "time"
    "testing"
    "github.com/barbell-math/block/db"
    //"github.com/barbell-math/block/util/test"
)

func TestGeneratePrediction(t *testing.T){
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0, FatigueIndex: 0,
        DatePerformed: time.Now(),
    };
    p,err:=GeneratePrediction(&testDB,&tl);
    fmt.Println(err,p);
}
