package db;

import (
    "time"
    "database/sql"
    "github.com/barbell-math/block/util/algo"
    customerr "github.com/barbell-math/block/util/err"
)

// func GetById[R DBTable](c *DB, id int) (R,error) {
//     return getRowFromUniqueValGenerator(func(searchId int) (R,string) {
//         var tmp R;
//         tmp.SetId(searchId);
//         return tmp,"Id";
//     })(c,id);
// }

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
var GetStateGeneratorByName=getRowFromUniqueValGenerator(
    func(name string) (StateGenerator,string) {
        return StateGenerator{T: name},"T";
});

type ValGenerator[R DBTable, V any] func(data V) (R,string);
type RowFromUniqueVal[R DBTable, V any] func(c *DB, data V) (R,error);
func getRowFromUniqueValGenerator[
        R DBTable,
        V any,
    ](valGen ValGenerator[R,V]) RowFromUniqueVal[R,V] {
    return func(c *DB, data V) (R,error){
        searchR,col:=valGen(data);
        if rv,err,found:=Read(c,searchR,algo.GenFilter(false,col)).Nth(0); rv!=nil && found {
            return *rv,err;
        } else if rv==nil && err==nil && !found {
            return *new(R),sql.ErrNoRows;
        } else {
            return *new(R),err;
        }
    }
}

func InitClient(
        db *DB,
        c *Client,
        sMax float64,
        bMax float64,
        dMax float64) error {
    f:=func(cId int, rId int, eId int, m float64) TrainingLog {
        return TrainingLog{
            ClientID: cId,
            ExerciseID: eId,
            RotationID: rId,
            DatePerformed: time.Now().AddDate(0, 0, -1),
            Weight: m,
            Sets: 1,
            Reps: 1,
            Intensity: float64(1),
        };
    }
    return customerr.ChainedErrorOps(
        func(r ...any) (any,error) { return GetExerciseByName(db,"Squat"); },
        func(r ...any) (any,error) { return GetExerciseByName(db,"Bench"); },
        func(r ...any) (any,error) { return GetExerciseByName(db,"Deadlift"); },
        func(r ...any) (any,error) { return Create(db,*c); },
        func(r ...any) (any,error) {
            return Create(db,Rotation{
                ClientID: r[3].([]int)[0],
                StartDate: time.Now().AddDate(0, 0, -1),
                EndDate: time.Now(),
            });
        }, func(r ...any) (any,error) {
            s:=f(r[3].([]int)[0],r[4].([]int)[0],r[0].(Exercise).Id,sMax);
            b:=f(r[3].([]int)[0],r[4].([]int)[0],r[1].(Exercise).Id,bMax);
            d:=f(r[3].([]int)[0],r[4].([]int)[0],r[2].(Exercise).Id,dMax);
            return Create(db,s,d,b);
        },
    );
}

func RmClient(db *DB, c *Client) (int64,error) {
    var rv int64=0;
    err:=customerr.ChainedErrorOps(
        func(r ...any) (any,error) {
            return CustomDeleteQuery(db,
                `DELETE FROM Prediction
                 WHERE Id IN (
                    SELECT Prediction.Id
                    FROM Prediction
                    JOIN TrainingLog
                    ON Prediction.TrainingLogID=TrainingLog.Id
                    WHERE TrainingLog.ClientID=$1
                 );`,[]any{c.Id},
            );
        }, func(r ...any) (any,error) {
            return Delete(
                db,TrainingLog{ClientID: c.Id},algo.GenFilter(false,"ClientID"),
            );
        }, func(r ...any) (any,error) {
            return Delete(
                db,Rotation{ClientID: c.Id},algo.GenFilter(false,"ClientID"),
            );
        }, func(r ...any) (any,error) {
            return Delete(
                db,BodyWeight{ClientID: c.Id},algo.GenFilter(false,"ClientID"),
            );
        }, func(r ...any) (any,error) {
            return Delete(
                db,ModelState{ClientID: c.Id},algo.GenFilter(false,"ClientID"),
            );
        }, func(r ...any) (any,error) { return Delete(db,*c,OnlyIDFilter); },
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
