package db;

import (
    "time"
    "database/sql"
    "github.com/carmichaeljr/powerlifting-engine/util"

    _ "github.com/lib/pq"
)

func GetClientID(c *CRUD, email string) (int,error) {
    if c,err:=getRowFromUniqueVal(c,Client{Email: email},"Email"); c!=nil {
        return c.Id,err;
    } else {
        return -1,err;
    }
}

//func GetExerciseTypeID(c *CRUD, type string) (int,error)
//func GetExerciseFocusID(c *CRUD, focus string) (int,error)

func GetExerciseID(c *CRUD, n string) (int,error) {
    if e,err:=getRowFromUniqueVal(c,Exercise{Name: n},"Name"); e!=nil {
        return e.Id,err;
    } else {
        return -1,err;
    }
}

func getRowFromUniqueVal[R DBTable](c *CRUD, row R, col string) (*R,error) {
    var rv *R=nil;
    err:=Read(c,row,GenColFilter(false,col),
        func(r *R){
            rv=r;
    });
    if rv==nil {
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
    return util.ChainedErrorOps(
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
}

func RmClient(crud *CRUD, c *Client) (int64,error) {
    var rv int64=0;
    err:=util.ChainedErrorOps(
        func(r ...any) (any,error) {
            return Delete(
                crud,TrainingLog{ClientID: c.Id},GenColFilter(false,"ClientID"),
            );
        }, func(r ...any) (any,error) {
            return Delete(
                crud,Rotation{ClientID: c.Id},GenColFilter(false,"ClientID"),
            );
        }, func(r ...any) (any,error) {
            return Delete(
                crud,BodyWeight{ClientID: c.Id},GenColFilter(false,"ClientID"),
            );
        }, func(r ...any) (any,error) {
            return Delete(
                crud,ModelState{ClientID: c.Id},GenColFilter(false,"ClientID"),
            );
        }, func(r ...any) (any,error) { return Delete(crud,*c,OnlyIDFilter); },
        func(r ...any) (any,error) {
            for _,v:=range(r) {
                rv+=v.(int64);
            }
            return nil,nil;
        },
    );
    return rv,err;
}

//func UpdateTrainingLogUsingCurRot(c *Client, e *Exercise, t *TrainingLog) (int,error) {
//
//}
