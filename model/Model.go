package model

import (
	stdMath "math"

	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/util/dataStruct"
	mathUtil "github.com/barbell-math/block/util/math"
)

type FatigueAwareModel struct {
    mathUtil.LinearReg[float64];
    mathUtil.LinRegResult[float64];
};

//Returns non-standard linear regression for the model according to the
//model equation.
//The ordering of the functions makes for this ordering of constants:
//  Eps,Eps1,Eps2,Eps3,Eps4,Eps5,Eps6,Eps7
func NewFatigueAwareModel() FatigueAwareModel {
    return FatigueAwareModel{
        LinearReg: mathUtil.NewLinearReg([]mathUtil.SummationOp[float64]{
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
            }},mathUtil.LinearSummationOp[float64]("I"),
        ),
    };
}

func (f *FatigueAwareModel)Run() (float64,error) {
    res,rcond,err:=f.LinearReg.Run();
    f.LinRegResult=res;
    f.imposeConstraints();
    return rcond,err;
}

func (f *FatigueAwareModel)imposeConstraints() error {
    var constraints=[...]dataStruct.Pair[float64,float64]{
        {A: 0, B: stdMath.Inf(1)}, //Eps: Error
        mathUtil.NoOpConstraint[float64](),     //Eps1: Effort
        mathUtil.NoOpConstraint[float64](),     //Eps2: F_l
        mathUtil.NoOpConstraint[float64](),     //Eps3: F_w
        mathUtil.NoOpConstraint[float64](),     //Eps4: F_e
        {A: 0, B: stdMath.Inf(1)}, //Eps5: s*r
        {A: 0, B: stdMath.Inf(1)}, //Eps6: s
        {A: 0, B: stdMath.Inf(1)}, //Eps7: r
    };
    for i,v:=range(f.V[0]) {
        f.V[0][i]=mathUtil.Constrain(v,constraints[i]);
    }
    return nil;
}

func (f *FatigueAwareModel)PredictFromDataPoint(dp dataPoint) (float64,error) {
    return f.Predict(map[string]float64{
        "I": dp.Intensity,
        "E": dp.Effort,
        "R": float64(dp.Reps),
        "S": float64(dp.Sets),
        "F_w": float64(dp.InterWorkoutFatigue),
        "F_e": float64(dp.InterExerciseFatigue),
        //"F_l": float64(dp.LatentFatigue),
    });
}

func (f *FatigueAwareModel)PredictFromTrainingLog(tl db.TrainingLog) (float64,error) {
    return f.Predict(map[string]float64{
        "I": tl.Intensity,
        "E": tl.Effort,
        "R": float64(tl.Reps),
        "S": float64(tl.Sets),
        "F_w": float64(tl.InterWorkoutFatigue),
        "F_e": float64(tl.InterExerciseFatigue),
        //"F_l": float64(dp.LatentFatigue),
    });
}

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
