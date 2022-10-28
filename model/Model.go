package model;

import (
    "fmt"
    "time"
    "math"
    "github.com/carmichaeljr/powerlifting-engine/db"
    "github.com/carmichaeljr/powerlifting-engine/mathUtil"
)

//func MakeIntensityPrediction(date time.Time){
//
//}

//Generates the model state given the date and exercise specified in the
//training log. Uses the training log data as the data that is being predicted,
//which means it needs to have all **VALID** values.
func GenerateModelState(c *db.CRUD, tl *db.TrainingLog) (db.ModelState,error) {
    type DataPoint struct {
        DatePerformed time.Time;
        Sets float64;
        Reps float64;
        Effort float64;
        Intensity float64;
    };
    query:=fmt.Sprintf(`SELECT TrainingLog.DatePerformed,
            TrainingLog.Sets,
            TrainingLog.Reps,
            TrainingLog.Effort,
            TrainingLog.Intensity
        FROM TrainingLog
        WHERE TrainingLog.ExerciseID=$1
            AND TrainingLog.DatePerformed<'%s'::date
        ORDER BY TrainingLog.DatePerformed DESC;`,
        tl.DatePerformed.Format("01/02/2006"),
    );
    //    tl.DatePerformed.Month(),tl.DatePerformed.Day(),tl.DatePerformed.Year(),
    //);
    lr:=fatigueAwareModel();
    fmt.Println(lr);
    //var curDate time.Time;
    err:=db.CustomReadQuery(c,query,[]any{tl.ExerciseID},func(d *DataPoint){
        //if curDate
        fmt.Println(d);
    });
    fmt.Println(err);
    return db.ModelState{},err;
}

//Returns non-standard linear regression for the model according to the
//equation below:
//  I=d-a(s-1)^2*(r-1)^2-b(s-1)^2-c(r-1)^2-eps_1*E-eps_2*F
//Where:
//  d,a,b,c,eps_1,eps_2 are the constants linear reg will find
//  s is sets
//  r is reps
//  E is effort (RPE)
//  F is the fatigue index
func fatigueAwareSumOpGen() ([]mathUtil.SummationOp[float64],
        mathUtil.SummationOp[float64]) {
    return []mathUtil.SummationOp[float64]{
        func(vals map[string]float64) (float64,error) {
            s,err:=mathUtil.VarAcc(vals,"S");
            if err!=nil {
                return 0, err;
            }
            r,err:=mathUtil.VarAcc(vals,"R");
            if err!=nil {
                return 0, err;
            }
            return math.Pow(s-1,2)*math.Pow(r-1,2),nil;
        }, func(vals map[string]float64) (float64,error) {
            s,err:=mathUtil.VarAcc(vals,"S");
            if err!=nil {
                return 0, err;
            }
            return math.Pow(s-1,2),nil;
        }, func(vals map[string]float64) (float64,error) {
            r,err:=mathUtil.VarAcc(vals,"R");
            if err!=nil {
                return 0, err;
            }
            return math.Pow(r-1,2),nil;
        }, mathUtil.LinearSummationOp[float64]("E"),
        mathUtil.LinearSummationOp[float64]("F"),
    },mathUtil.LinearSummationOp[float64]("I");
}
func fatigueAwareModel() mathUtil.LinearReg[float64] {
    return mathUtil.NewLinearReg[float64](fatigueAwareSumOpGen());
}
