package potentialSurface

import (
	stdMath "math"

	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/util/dataStruct"
	mathUtil "github.com/barbell-math/block/util/math"
)

//The basic surface follows the following equation:
//  I=eps+eps_1*E-( eps_2*F_w+eps_3*F_e+eps_4*(s-1)^2(r-1)^2+eps_5*(s-1)^2+eps_6*(r-1)^2 )
//This equation does not take into account latent fatigue, which makes it naive
//because it does not consider the relationship between lifts across time.
var BasicSurfacePrediction basicSurfacePrediction;
type basicSurfacePrediction struct { };

func (b basicSurfacePrediction)Intensity(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    return (ms.Eps+
        ms.Eps1*float64(tl.Effort)-
        ms.Eps2*float64(tl.InterWorkoutFatigue)-
        ms.Eps3*float64(tl.InterExerciseFatigue)-
        ms.Eps4*stdMath.Pow(float64(tl.Sets-1),2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps5*stdMath.Pow(float64(tl.Sets-1),2)-
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2));
}

func (b basicSurfacePrediction)Effort(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    return (tl.Intensity-ms.Eps+
        ms.Eps2*float64(tl.InterWorkoutFatigue)+
        ms.Eps3*float64(tl.InterExerciseFatigue)+
        ms.Eps4*stdMath.Pow(float64(tl.Sets-1),2)*stdMath.Pow(float64(tl.Reps-1),2)+
        ms.Eps5*stdMath.Pow(float64(tl.Sets-1),2)+
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2))/ms.Eps1;
}

func (b basicSurfacePrediction)InterWorkoutFatigue(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    return (ms.Eps+
        ms.Eps1*tl.Effort-
        ms.Eps3*float64(tl.InterExerciseFatigue)-
        ms.Eps4*stdMath.Pow(float64(tl.Sets-1),2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps5*stdMath.Pow(float64(tl.Sets-1),2)-
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2)-
        tl.Intensity)/ms.Eps2;
}

func (b basicSurfacePrediction)InterExerciseFatigue(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    return (ms.Eps+
        ms.Eps1*tl.Effort-
        ms.Eps2*float64(tl.InterWorkoutFatigue)-
        ms.Eps4*stdMath.Pow(float64(tl.Sets-1),2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps5*stdMath.Pow(float64(tl.Sets-1),2)-
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2)-
        tl.Intensity)/ms.Eps3;
}

func (b basicSurfacePrediction)Sets(
        ms *db.ModelState,
        tl *db.TrainingLog) float64 {
    return stdMath.Pow((
        ms.Eps+
        ms.Eps1*tl.Effort-
        ms.Eps2*float64(tl.InterWorkoutFatigue)-
        ms.Eps3*float64(tl.InterExerciseFatigue)-
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2)-
        tl.Intensity)/(
        ms.Eps4*stdMath.Pow(float64(tl.Reps-1),2)+
        ms.Eps5),0.5)+1.0;
}

func (b basicSurfacePrediction)Reps(
        ms *db.ModelState,
        tl *db.TrainingLog) float64 {
    return stdMath.Pow((
        ms.Eps+
        ms.Eps1*tl.Effort-
        ms.Eps2*float64(tl.InterWorkoutFatigue)-
        ms.Eps3*float64(tl.InterExerciseFatigue)-
        ms.Eps5*stdMath.Pow(float64(tl.Sets-1),2)-
        tl.Intensity)/(
        ms.Eps4*stdMath.Pow(float64(tl.Sets-1),2)+
        ms.Eps6),0.5)+1.0;
}

type BasicSurface struct {
    mathUtil.LinearReg[float64];
    mathUtil.LinRegResult[float64];
};

//The ordering of the functions makes for this ordering of constants:
//  Eps,Eps1,Eps2,Eps3,Eps4,Eps5,Eps6
func NewBasicSurface() BasicSurface {
    return BasicSurface{
        LinearReg: mathUtil.NewLinearReg([]mathUtil.SummationOp[float64]{
            mathUtil.ConstSummationOp[float64](1),
            mathUtil.LinearSummationOp[float64]("E"),
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

func (b *BasicSurface)Id() PotentialSurfaceId { return BasicSurfaceId; }
func (b *BasicSurface)Predictor() Predictor { return BasicSurfacePrediction; }

func (b *BasicSurface)Update(vals map[string]float64) error {
    return b.UpdateSummations(vals);
}

func (b *BasicSurface)Run() (float64,error) {
    res,rcond,err:=b.LinearReg.Run();
    b.LinRegResult=res;
    b.imposeConstraints();
    return rcond,err;
}

func (b *BasicSurface)imposeConstraints() {
    var constraints=[...]dataStruct.Pair[float64,float64]{
        {A: 0, B: stdMath.Inf(1)}, //Eps: Error
        mathUtil.NoOpConstraint[float64](),     //Eps1: Effort
        mathUtil.NoOpConstraint[float64](),     //Eps2: F_w
        mathUtil.NoOpConstraint[float64](),     //Eps3: F_e
        {A: 0, B: stdMath.Inf(1)}, //Eps4: s*r
        {A: 0, B: stdMath.Inf(1)}, //Eps5: s
        {A: 0, B: stdMath.Inf(1)}, //Eps6: r
    };
    for i,v:=range(b.LinRegResult.V) {
        b.LinRegResult.V[i][0]=mathUtil.Constrain(v[0],constraints[i]);
    }
}

// func (f *FatigueAwareModel)PredictIntensityFromDataPoint(dp dataPoint) (float64,error) {
//     return f.Predict(map[string]float64{
//         "I": dp.Intensity,
//         "E": dp.Effort,
//         "R": float64(dp.Reps),
//         "S": float64(dp.Sets),
//         "F_w": float64(dp.InterWorkoutFatigue),
//         "F_e": float64(dp.InterExerciseFatigue),
//         //"F_l": float64(dp.LatentFatigue),
//     });
// }

func (b *BasicSurface)PredictIntensity(tl *db.TrainingLog) (float64,error) {
    return b.Predict(map[string]float64{
        "I": tl.Intensity,
        "E": tl.Effort,
        "R": float64(tl.Reps),
        "S": float64(tl.Sets),
        "F_w": float64(tl.InterWorkoutFatigue),
        "F_e": float64(tl.InterExerciseFatigue),
    });
}
