package model;

import (
    "math"
    "time"
    "github.com/carmichaeljr/powerlifting-engine/db"
    "github.com/carmichaeljr/powerlifting-engine/mathUtil"
)

//The struct that holds values when linear regression is performed.
//Note - THE ORDER OF THE STRUCT FIELDS MUST MATCH THE ORDER OF THE VALUES
//IN THE QUERY. Otherwise the values returned will be all jumbled up.
type dataPoint struct {
    DatePerformed time.Time;
    Sets float64;
    Reps float64;
    Effort float64;
    Intensity float64;
    FatigueIndex float64;
};

//The model equation is as follows:
//  I=d-a(s-1)^2*(r-1)^2-b(s-1)^2-c(r-1)^2-eps_1*E-eps_2*F
//Where:
//  d,a,b,c,eps_1,eps_2 are the constants linear reg will find
//  s is sets
//  r is reps
//  E is effort (RPE)
//  F is the fatigue index
func IntensityPrediction(ms *db.ModelState, tl *db.TrainingLog) float64 {
    return (ms.D-
            ms.A*math.Pow(float64(tl.Sets-1),2)*math.Pow(float64(tl.Reps-1),2)-
            ms.B*math.Pow(float64(tl.Sets-1),2)-
            ms.C*math.Pow(float64(tl.Reps-1),2)-
            ms.Eps*tl.Effort-
            ms.Eps2*float64(tl.FatigueIndex));
}

//Returns non-standard linear regression for the model according to the
//model equation.
func fatigueAwareModel() mathUtil.LinearReg[float64] {
    return mathUtil.NewLinearReg[float64](fatigueAwareSumOpGen());
}
//The ordering of the functions makes for this ordering of constants:
//  d,a,b,c,eps_1,eps_2
func fatigueAwareSumOpGen() ([]mathUtil.SummationOp[float64],
        mathUtil.SummationOp[float64]) {
    return []mathUtil.SummationOp[float64]{
        mathUtil.ConstSummationOp[float64](1),
        func(vals map[string]float64) (float64,error) {
            s,err:=mathUtil.VarAcc(vals,"S");
            if err!=nil {
                return 0, err;
            }
            r,err:=mathUtil.VarAcc(vals,"R");
            if err!=nil {
                return 0, err;
            }
            return -(math.Pow(s-1,2)*math.Pow(r-1,2)),nil;
        }, func(vals map[string]float64) (float64,error) {
            s,err:=mathUtil.VarAcc(vals,"S");
            if err!=nil {
                return 0, err;
            }
            return -math.Pow(s-1,2),nil;
        }, func(vals map[string]float64) (float64,error) {
            r,err:=mathUtil.VarAcc(vals,"R");
            if err!=nil {
                return 0, err;
            }
            return -math.Pow(r-1,2),nil;
        }, mathUtil.NegatedLinearSummationOp[float64]("E"),
        mathUtil.NegatedLinearSummationOp[float64]("F"),
    },mathUtil.LinearSummationOp[float64]("I");
}

func getPredFromLinRegResult(
        res mathUtil.LinRegResult[float64],
        tl *db.TrainingLog) float64 {
    rv,_:=res.Predict(map[string]float64{
        "F": float64(tl.FatigueIndex), "I": tl.Intensity, "R": float64(tl.Reps),
        "E": tl.Effort, "S": float64(tl.Sets),
    });
    return rv;
}
