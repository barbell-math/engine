package db;

import (
    "time"
    "database/sql"
    "github.com/barbell-math/block/util/algo"
    customerr "github.com/barbell-math/block/util/err"
)

var GetClientByEmail=getRowFromUniqueValGenerator(
    func(email string) (Client,string) {
        return Client{Email: email},"Email";
});
var GetExerciseByName=getRowFromUniqueValGenerator(
    func(name string) (Exercise,string) {
        return Exercise{Name: name},"Name";
});
var GetExerciseTypeByName=getRowFromUniqueValGenerator(
    func(_type string) (ExerciseType,string) {
        return ExerciseType{T: _type},"T";
});
var GetExerciseFocusByName=getRowFromUniqueValGenerator(
    func(focus string) (ExerciseFocus,string) {
        return ExerciseFocus{Focus: focus},"Focus";
});

type ValGenerator[R DBTable, V any] func(data V) (R,string);
type RowFromUniqueVal[R DBTable, V any] func(c *DB, data V) (R,error);
func getRowFromUniqueValGenerator[
        R DBTable,
        V any,
    ](valGen ValGenerator[R,V]) RowFromUniqueVal[R,V] {
    return func(c *DB, data V) (R,error){
        var rv *R=nil;
        searchR,col:=valGen(data);
        if err:=Read(c,searchR,algo.GenFilter(false,col),func(r *R){
            rv=r;
        }); rv!=nil {
            return *rv,err;
        } else if rv==nil && err==nil {
            return *new(R),sql.ErrNoRows;
        } else {
            return *new(R),err;
        }
    }
}

func InitClient(
        crud *DB,
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
    return customerr.ChainedErrorOps(
        func(r ...any) (any,error) { return GetExerciseByName(crud,"Squat"); },
        func(r ...any) (any,error) { return GetExerciseByName(crud,"Bench"); },
        func(r ...any) (any,error) { return GetExerciseByName(crud,"Deadlift"); },
        func(r ...any) (any,error) { return Create(crud,*c); },
        func(r ...any) (any,error) {
            return Create(crud,Rotation{
                ClientID: r[3].([]int)[0],
                StartDate: time.Now().AddDate(0, 0, -1),
                EndDate: time.Now(),
            });
        }, func(r ...any) (any,error) {
            s:=f(r[3].([]int)[0],r[4].([]int)[0],r[0].(Exercise).Id,sMax);
            b:=f(r[3].([]int)[0],r[4].([]int)[0],r[1].(Exercise).Id,bMax);
            d:=f(r[3].([]int)[0],r[4].([]int)[0],r[2].(Exercise).Id,dMax);
            return Create(crud,s,d,b);
        },
    );
}

func RmClient(crud *DB, c *Client) (int64,error) {
    var rv int64=0;
    err:=customerr.ChainedErrorOps(
        func(r ...any) (any,error) {
            return Delete(
                crud,TrainingLog{ClientID: c.Id},algo.GenFilter(false,"ClientID"),
            );
        }, func(r ...any) (any,error) {
            return Delete(
                crud,Rotation{ClientID: c.Id},algo.GenFilter(false,"ClientID"),
            );
        }, func(r ...any) (any,error) {
            return Delete(
                crud,BodyWeight{ClientID: c.Id},algo.GenFilter(false,"ClientID"),
            );
        }, func(r ...any) (any,error) {
            return Delete(
                crud,ModelState{ClientID: c.Id},algo.GenFilter(false,"ClientID"),
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
