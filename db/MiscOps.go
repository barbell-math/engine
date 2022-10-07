package db;

import (
    //"fmt"
    "time"
    "database/sql"
    "github.com/carmichaeljr/powerlifting-engine/util"

    _ "github.com/lib/pq"
)

func GetExerciseID(c *CRUD, n string) (int,error) {
    rv:=-1;
    err:=Read(c,Exercise{Name: n},GenColFilter(false,"Name"),
        func(e *Exercise){
            rv=e.Id;
    });
    if rv==-1 {
        return rv,sql.ErrNoRows;
    }
    return rv,err;
}

func InitClient(
        crud *CRUD,
        c *Client,
        sMax float32,
        bMax float32,
        dMax float32) error {
    f:=func(cId int, rId int, eId int, m float32) TrainingLog {
        return TrainingLog{
            ClientID: cId,
            ExerciseID: eId,
            RotationID: rId,
            DatePerformed: time.Now().AddDate(0, 0, -1),
            Weight: m,
            Sets: float32(1),
            Reps: 1,
            Intensity: float64(1),
        };
    }
    err:=util.ChainedErrorOps(
        func(r ...any) (any,error) { return GetExerciseID(crud,"Squat"); },
        func(r ...any) (any,error) { return GetExerciseID(crud,"Bench"); },
        func(r ...any) (any,error) { return GetExerciseID(crud,"Deadlift"); },
        func(r ...any) (any,error) { return Create(crud,*c); },
        func(r ...any) (any,error) {
            return Create(crud,Rotation{
                ClientID: r[3].([]int)[0],
                StartDate: time.Now().AddDate(0, 0, -1),
                EndDate: time.Now(),
            });
        }, func(r ...any) (any,error) {
            s:=f(r[3].([]int)[0],r[4].([]int)[0],r[0].(int),sMax);
            b:=f(r[3].([]int)[0],r[4].([]int)[0],r[1].(int),bMax);
            d:=f(r[3].([]int)[0],r[4].([]int)[0],r[2].(int),dMax);
            return Create(crud,s,d,b);
        },
    );
    return err;
}

//func AddToTrainingLog(c *Client, e *Exercise, r *Rotation, m float32) (int,error) {
//
//}
