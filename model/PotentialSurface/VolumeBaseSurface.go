package potentialSurface

import (
	//"fmt"
	stdMath "math"

	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/util/dataStruct"
	mathUtil "github.com/barbell-math/block/util/math"
)

//The basic surface follows the following equation:
//  I^2(F+eps*E)=E
// Where:
//  E=eps_1*E
//  F=1+eps_2*F_w+eps_3*F_e+eps_4*(s-1)^2(r-1)^2+eps_5*(s-1)^2+eps_6*(r-1)^2
//This equation does not take into account latent fatigue, which makes it naive
//because it does not consider the relationship between lifts across time.
var VolumeBaseSurfacePrediction volumeBaseSurfacePrediction;
type volumeBaseSurfacePrediction struct {};

func (v volumeBaseSurfacePrediction)e(
        ms *db.ModelState,
        tl *db.TrainingLog) float64 {
    return ms.Eps1*float64(tl.Effort);
}

func (v volumeBaseSurfacePrediction)f(
        ms *db.ModelState,
        tl *db.TrainingLog) float64 {
    return (1+ms.Eps2*float64(tl.InterWorkoutFatigue)+
        ms.Eps3*float64(tl.InterExerciseFatigue)+
        ms.Eps4*stdMath.Pow(float64(tl.Sets-1),2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps5*stdMath.Pow(float64(tl.Sets-1),2)-
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2));
}

func (v volumeBaseSurfacePrediction)Intensity(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    E:=v.e(ms,tl);
    F:=v.f(ms,tl);
    return stdMath.Pow(E/(F+ms.Eps*E),0.5);
}

func (v volumeBaseSurfacePrediction)Effort(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    F:=v.f(ms,tl);
    pSq:=tl.Intensity*tl.Intensity;
    return (pSq*F)/(ms.Eps1*(1-ms.Eps*pSq));
}
func (v volumeBaseSurfacePrediction)InterWorkoutFatigue(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    E:=v.e(ms,tl);
    pSq:=tl.Intensity*tl.Intensity;
    return ((E*(1/pSq-ms.Eps)-1-ms.Eps3*float64(tl.InterExerciseFatigue)-
        ms.Eps4*stdMath.Pow(float64(tl.Sets)-1,2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps5*stdMath.Pow(float64(tl.Sets)-1,2)-
        ms.Eps6*stdMath.Pow(float64(tl.Reps)-1,2))/ms.Eps2);
}

func (v volumeBaseSurfacePrediction)InterExerciseFatigue(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    E:=v.e(ms,tl);
    pSq:=tl.Intensity*tl.Intensity;
    return ((E*(1/pSq-ms.Eps)-1-ms.Eps2*float64(tl.InterWorkoutFatigue)-
        ms.Eps4*stdMath.Pow(float64(tl.Sets)-1,2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps5*stdMath.Pow(float64(tl.Sets)-1,2)-
        ms.Eps6*stdMath.Pow(float64(tl.Reps)-1,2))/ms.Eps3);
}

func (v volumeBaseSurfacePrediction)Sets(
        ms *db.ModelState,
        tl *db.TrainingLog) float64 {
    E:=v.e(ms,tl);
    pSq:=tl.Intensity*tl.Intensity;
    return stdMath.Pow((E*(1/pSq-ms.Eps)-1-
        ms.Eps2*float64(tl.InterExerciseFatigue)-
        ms.Eps3*float64(tl.InterWorkoutFatigue)-
        ms.Eps6*stdMath.Pow(float64(tl.Reps)-1,2))/
    (ms.Eps4*stdMath.Pow(float64(tl.Reps)-1,2)+ms.Eps5),0.5)+1
}

func (v volumeBaseSurfacePrediction)Reps(
        ms *db.ModelState,
        tl *db.TrainingLog) float64 {
    E:=v.e(ms,tl);
    pSq:=tl.Intensity*tl.Intensity;
    return stdMath.Pow((E*(1/pSq-ms.Eps)-1-
        ms.Eps2*float64(tl.InterExerciseFatigue)-
        ms.Eps3*float64(tl.InterWorkoutFatigue)-
        ms.Eps5*stdMath.Pow(float64(tl.Sets)-1,2))/
    (ms.Eps4*stdMath.Pow(float64(tl.Sets)-1,2)+ms.Eps6),0.5)+1
}

type VolumeBaseSurface struct {
    mathUtil.LinearReg[float64];
    mathUtil.LinRegResult[float64];
};

//The ordering of the functions makes for this ordering of constants:
//  Eps,Eps1,Eps2,Eps3,Eps4,Eps5,Eps6
func NewVolumeBaseSurface() VolumeBaseSurface {
    return VolumeBaseSurface{
        LinearReg: mathUtil.NewLinearReg([]mathUtil.SummationOp[float64]{
            mathUtil.ConstSummationOp[float64](1),
            func(vals map[string]float64) (float64, error) {
                e,err:=mathUtil.VarAcc(vals,"E");
                if err!=nil {
                    return 0,err;
                }
                return 1/e,nil;
            }, func(vals map[string]float64) (float64, error) {
                e,err:=mathUtil.VarAcc(vals,"E");
                if err!=nil {
                    return 0,err;
                }
                fw,err:=mathUtil.VarAcc(vals,"F_w");
                if err!=nil {
                    return 0,err;
                }
                return fw/e,nil;
            }, func(vals map[string]float64) (float64, error) {
                e,err:=mathUtil.VarAcc(vals,"E");
                if err!=nil {
                    return 0,err;
                }
                fe,err:=mathUtil.VarAcc(vals,"F_e");
                if err!=nil {
                    return 0,err;
                }
                return fe/e,nil;
            }, func(vals map[string]float64) (float64, error) {
                e,err:=mathUtil.VarAcc(vals,"E");
                if err!=nil {
                    return 0,err;
                }
                s,err:=mathUtil.VarAcc(vals,"S");
                if err!=nil {
                    return 0,err;
                }
                r,err:=mathUtil.VarAcc(vals,"R");
                if err!=nil {
                    return 0,err;
                }
                return stdMath.Pow(s-1,2)*stdMath.Pow(r-1,2)/e,nil;
            }, func(vals map[string]float64) (float64, error) {
                e,err:=mathUtil.VarAcc(vals,"E");
                if err!=nil {
                    return 0,err;
                }
                s,err:=mathUtil.VarAcc(vals,"S");
                if err!=nil {
                    return 0,err;
                }
                return stdMath.Pow(s-1,2)/e,nil;
            }, func(vals map[string]float64) (float64, error) {
                e,err:=mathUtil.VarAcc(vals,"E");
                if err!=nil {
                    return 0,err;
                }
                r,err:=mathUtil.VarAcc(vals,"R");
                if err!=nil {
                    return 0,err;
                }
                return stdMath.Pow(r-1,2)/e,nil;
            }},func(vals map[string]float64) (float64,error) {
                i,err:=mathUtil.VarAcc(vals,"I");
                if err!=nil {
                    return 0,err;
                }
                return 1/(i*i),err;
            },
        ),
    };
}

func (v VolumeBaseSurface)ToGenericSurf() Surface { return &v; }

func (v *VolumeBaseSurface)Id() PotentialSurfaceId { return VolumeBaseSurfaceId; }
func (v *VolumeBaseSurface)Predictor() Predictor { return VolumeBaseSurfacePrediction; }

func (v *VolumeBaseSurface)Update(vals map[string]float64) error {
    return v.UpdateSummations(vals);
}

func (v *VolumeBaseSurface)Run() (float64,error) {
    res,rcond,err:=v.LinearReg.Run();
    v.LinRegResult=res;
    v.imposeConstraints();
    return rcond,err;
}

func (v *VolumeBaseSurface)imposeConstraints() {
    var constraints=[...]dataStruct.Pair[float64,float64]{
        {A: 0, B: stdMath.Inf(1)}, //Eps: Error
        {A: 0, B: stdMath.Inf(1)}, //Eps1: Effort
        //mathUtil.NoOpConstraint[float64](),     //Eps1: Effort
        mathUtil.NoOpConstraint[float64](),     //Eps2: F_w
        mathUtil.NoOpConstraint[float64](),     //Eps3: F_e
        {A: 0, B: stdMath.Inf(1)}, //Eps4: s*r
        {A: 0, B: stdMath.Inf(1)}, //Eps5: s
        {A: 0, B: stdMath.Inf(1)}, //Eps6: r
    };
    for i,iterV:=range(v.LinRegResult.V) {
        v.LinRegResult.V[i][0]=mathUtil.Constrain(iterV[0],constraints[i]);
    }
}

// func (v *VolumeBaseSurface)GetConstant(i int) float64 {
//     // Requested eps_1, scale and return
//     if i==1 {
//         return 1/v.LinRegResult.V[1][0];
//     } else if i==0 {
//         // Requested value is eps, just return
//         return v.LinRegResult.V[0][0];
//     }
//     tmp:=v.LinRegResult.GetConstant(i);
//     // Requested eps_2...eps_6, scale the value and return
//     rv:=tmp/v.LinRegResult.V[1][0];
//     if stdMath.Abs(rv)==stdMath.Inf(1) || rv==stdMath.NaN() {
//         fmt.Println("Working precision has been reached!");
//         fmt.Println("Const ",i,": ",tmp);
//         fmt.Println("Eps: ",v.LinRegResult.V[1][0]);
//         fmt.Println("rv: ",rv);
//     }
//     return tmp/v.LinRegResult.V[1][0];
// }

func (v *VolumeBaseSurface)PredictIntensity(vals map[string]float64) (float64,error) {
    tmp,err:=v.LinRegResult.Predict(vals);
    if err!=nil {
        return tmp,err;
    }
    return 1/stdMath.Pow(tmp,0.5),err;
}
