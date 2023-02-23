package model;

import (
    stdMath "math"
    "github.com/barbell-math/block/db"
    mathUtil "github.com/barbell-math/block/util/math"
)

type ModelPredictor func(ms *db.ModelState, tl *db.TrainingLog) float64;
func IntensityPrediction(ms *db.ModelState, tl *db.TrainingLog) float64 {
    return (ms.Eps+
        ms.Eps1*float64(tl.Effort)-
        //ms.Eps2*float64(tl.LatentFatigueindex)-
        ms.Eps3*float64(tl.InterWorkoutFatigue)-
        ms.Eps4*float64(tl.InterExerciseFatigue)-
        ms.Eps5*stdMath.Pow(float64(tl.Sets-1),2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps6*stdMath.Pow(float64(tl.Sets-1),2)-
        ms.Eps7*stdMath.Pow(float64(tl.Reps-1),2));
}
func EffortPrediction(ms *db.ModelState, tl *db.TrainingLog) float64 {
    return (tl.Intensity-ms.Eps+
        //ms.Eps2*float64(tl.LatentFatigue)+
        ms.Eps3*float64(tl.InterWorkoutFatigue)+
        ms.Eps4*float64(tl.InterExerciseFatigue)+
        ms.Eps5*stdMath.Pow(float64(tl.Sets-1),2)*stdMath.Pow(float64(tl.Reps-1),2)+
        ms.Eps6*stdMath.Pow(float64(tl.Sets-1),2)+
        ms.Eps7*stdMath.Pow(float64(tl.Reps-1),2))/ms.Eps1;
}
func LatentFatiguePrediction(ms *db.ModelState, tl *db.TrainingLog) float64 {
    return (ms.Eps+
        ms.Eps1*tl.Effort-
        ms.Eps3*float64(tl.InterWorkoutFatigue)-
        ms.Eps4*float64(tl.InterExerciseFatigue)-
        ms.Eps5*stdMath.Pow(float64(tl.Sets-1),2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps6*stdMath.Pow(float64(tl.Sets-1),2)-
        ms.Eps7*stdMath.Pow(float64(tl.Reps-1),2)-
        tl.Intensity)/ms.Eps2;
}
func InterWorkoutFatiguePrediciton(ms *db.ModelState, tl *db.TrainingLog) float64 {
    return (ms.Eps+
        ms.Eps1*tl.Effort-
        //ms.Eps2*float64(tl.LatentFatigue)-
        ms.Eps4*float64(tl.InterExerciseFatigue)-
        ms.Eps5*stdMath.Pow(float64(tl.Sets-1),2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps6*stdMath.Pow(float64(tl.Sets-1),2)-
        ms.Eps7*stdMath.Pow(float64(tl.Reps-1),2)-
        tl.Intensity)/ms.Eps3;
}
func InterExerciseFatiguePrediction(ms *db.ModelState, tl *db.TrainingLog) float64 {
    return (ms.Eps+
        ms.Eps1*tl.Effort-
        //ms.Eps2*float64(tl.LatentFatigue)-
        ms.Eps3*float64(tl.InterWorkoutFatigue)-
        ms.Eps5*stdMath.Pow(float64(tl.Sets-1),2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps6*stdMath.Pow(float64(tl.Sets-1),2)-
        ms.Eps7*stdMath.Pow(float64(tl.Reps-1),2)-
        tl.Intensity)/ms.Eps4;
}
func SetsPrediction(ms *db.ModelState, tl *db.TrainingLog) float64 {
    return stdMath.Pow((
        ms.Eps+
        ms.Eps1*tl.Effort-
        //ms.Eps2*tl.LatentFatigue-
        ms.Eps3*float64(tl.InterWorkoutFatigue)-
        ms.Eps4*float64(tl.InterExerciseFatigue)-
        ms.Eps7*stdMath.Pow(float64(tl.Reps-1),2)-
        tl.Intensity)/(
        ms.Eps5*stdMath.Pow(float64(tl.Reps-1),2)+
        ms.Eps6),0.5)+1.0;
}

func RepsPrediction(ms *db.ModelState, tl *db.TrainingLog) float64 {
    return stdMath.Pow((
        ms.Eps+
        ms.Eps1*tl.Effort-
        //ms.Eps2*tl.LatentFatigue-
        ms.Eps3*float64(tl.InterWorkoutFatigue)-
        ms.Eps4*float64(tl.InterExerciseFatigue)-
        ms.Eps6*stdMath.Pow(float64(tl.Sets-1),2)-
        tl.Intensity)/(
        ms.Eps5*stdMath.Pow(float64(tl.Sets-1),2)+
        ms.Eps7),0.5)+1.0;
}

//Returns non-standard linear regression for the model according to the
//model equation.
func fatigueAwareModel() mathUtil.LinearReg[float64] {
    return mathUtil.NewLinearReg[float64](fatigueAwareSumOpGen());
}
//The ordering of the functions makes for this ordering of constants:
//  Eps,Eps1,Eps2,Eps3,Eps4,Eps5,Eps6,Eps7
func fatigueAwareSumOpGen() ([]mathUtil.SummationOp[float64],
        mathUtil.SummationOp[float64]) {
    return []mathUtil.SummationOp[float64]{
        mathUtil.ConstSummationOp[float64](1),
        mathUtil.LinearSummationOp[float64]("E"),
        //mathUtil.NegatedLinearSummationOp[float64]("F_l"),
        mathUtil.NegatedLinearSummationOp[float64]("F_w"),
        mathUtil.NegatedLinearSummationOp[float64]("F_e"),
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
    }},mathUtil.LinearSummationOp[float64]("I");
}

func intensityPredFromLinReg[V *db.TrainingLog | *dataPoint](
        res mathUtil.LinRegResult[float64],
        v V) (float64,error) {
    var dp dataPoint;
    switch any(v).(type) {
        case db.TrainingLog:
            tmp:=any(v).(db.TrainingLog);
            dp=dataPoint{
                Intensity: tmp.Intensity,
                Sets: float64(tmp.Sets),
                Reps: float64(tmp.Reps),
                Effort: tmp.Effort,
                InterWorkoutFatigue: float64(tmp.InterWorkoutFatigue),
                InterExerciseFatigue: float64(tmp.InterExerciseFatigue),
            };
        case dataPoint:
            dp=any(v).(dataPoint);
    }
    return res.Predict(map[string]float64{
        "I": dp.Intensity,
        "E": dp.Effort,
        "R": float64(dp.Reps),
        "S": float64(dp.Sets),
        "F_w": float64(dp.InterWorkoutFatigue),
        "F_e": float64(dp.InterExerciseFatigue),
        //"F_l": float64(dp.LatentFatigue),
    });
}
