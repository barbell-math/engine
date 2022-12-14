package model;

import (
    stdMath "math"
    "github.com/barbell-math/block/db"
    mathUtil "github.com/barbell-math/block/util/math"
)

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
        ms.A*stdMath.Pow(float64(tl.Sets-1),2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.B*stdMath.Pow(float64(tl.Sets-1),2)-
        ms.C*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps*tl.Effort-
        ms.Eps2*float64(tl.FatigueIndex));
}
func EffortPrediction(ms *db.ModelState, tl *db.TrainingLog) float64 {
    return (ms.D-
        tl.Intensity-
        ms.A*stdMath.Pow(float64(tl.Sets-1),2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.B*stdMath.Pow(float64(tl.Sets-1),2)-
        ms.C*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps2*float64(tl.FatigueIndex))/(ms.Eps);
}
func SetsPrediction(ms *db.ModelState, tl *db.TrainingLog) float64 {
    return stdMath.Pow((ms.D-
        tl.Intensity-
        ms.C*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps*tl.Effort-
        ms.Eps2*float64(tl.FatigueIndex))/
        (ms.A*stdMath.Pow(float64(tl.Reps-1),2)+ms.B),0.5)+1.0;
}
func RepsPrediction(ms *db.ModelState, tl *db.TrainingLog) float64 {
    return stdMath.Pow((ms.D-
        tl.Intensity-
        ms.B*stdMath.Pow(float64(tl.Sets-1),2)-
        ms.Eps*tl.Effort-
        ms.Eps2*float64(tl.FatigueIndex))/
        (ms.A*stdMath.Pow(float64(tl.Sets-1),2)+ms.C),0.5)+1.0;
}
func FatigueIndexPrediction(ms *db.ModelState, tl *db.TrainingLog) float64 {
    return (ms.D-
        tl.Intensity-
        ms.A*stdMath.Pow(float64(tl.Sets-1),2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.B*stdMath.Pow(float64(tl.Sets-1),2)-
        ms.C*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps*tl.Effort)/(ms.Eps2);
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
            return -(stdMath.Pow(s-1,2)*stdMath.Pow(r-1,2)),nil;
        }, func(vals map[string]float64) (float64,error) {
            s,err:=mathUtil.VarAcc(vals,"S");
            if err!=nil {
                return 0, err;
            }
            return -stdMath.Pow(s-1,2),nil;
        }, func(vals map[string]float64) (float64,error) {
            r,err:=mathUtil.VarAcc(vals,"R");
            if err!=nil {
                return 0, err;
            }
            return -stdMath.Pow(r-1,2),nil;
        }, mathUtil.NegatedLinearSummationOp[float64]("E"),
        mathUtil.NegatedLinearSummationOp[float64]("F"),
    },mathUtil.LinearSummationOp[float64]("I");
}

func intensityPredFromLinReg(
        res mathUtil.LinRegResult[float64],
        tl *db.TrainingLog) (float64,error) {
    return res.Predict(map[string]float64{
        "F": float64(tl.FatigueIndex), "I": tl.Intensity, "R": float64(tl.Reps),
        "E": tl.Effort, "S": float64(tl.Sets),
    });
}
