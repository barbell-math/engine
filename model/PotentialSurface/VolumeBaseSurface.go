package potentialSurface

import (
	stdMath "math"

	"github.com/barbell-math/engine/db"
	"github.com/barbell-math/engine/util/dataStruct"
	mathUtil "github.com/barbell-math/engine/util/math/numeric"
)

//The volume base surface follows the following equation:
//  I^2=eps+E_tot/F_tot
// Where:
//  E_tot=eps_1*E
//  F_tot=eps_2*F_w+eps_3*F_e+eps_4*(s-1)^2(r-1)^2+eps_5*(s-1)^2+eps_6*(r-1)^2
//This equation does not take into account latent fatigue, which makes it naive
//because it does not consider the relationship between lifts across time.
var VolumeBaseSurfacePrediction volumeBaseSurfacePrediction;
type volumeBaseSurfacePrediction struct {};

func (v volumeBaseSurfacePrediction)Intensity(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    return stdMath.Sqrt(1/(
        ms.Eps+
        ms.Eps1/tl.Effort+
        ms.Eps2*float64(tl.InterWorkoutFatigue)/tl.Effort+
        ms.Eps3*float64(tl.InterExerciseFatigue)/tl.Effort+
        ms.Eps4*stdMath.Pow(tl.Sets-1,2)*stdMath.Pow(float64(tl.Reps-1),2)/tl.Effort+
        ms.Eps5*stdMath.Pow(tl.Sets-1,2)/tl.Effort+
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2)/tl.Effort));
}

func (v volumeBaseSurfacePrediction)Effort(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    return (-tl.Intensity*tl.Intensity*(
        ms.Eps1+
        ms.Eps2*float64(tl.InterWorkoutFatigue)+
        ms.Eps3*float64(tl.InterExerciseFatigue)+
        ms.Eps4*stdMath.Pow(tl.Sets-1,2)*stdMath.Pow(float64(tl.Reps-1),2)+
        ms.Eps5*stdMath.Pow(tl.Sets-1,2)+
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2)))/(
            1-tl.Intensity*tl.Intensity*ms.Eps);
}

func (v volumeBaseSurfacePrediction)InterWorkoutFatigue(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    return (1/stdMath.Pow(tl.Intensity,2)-
        ms.Eps-
        ms.Eps1/tl.Effort-
        ms.Eps3*float64(tl.InterExerciseFatigue)/tl.Effort-
        ms.Eps4*stdMath.Pow(tl.Sets-1,2)*stdMath.Pow(float64(tl.Reps-1),2)/tl.Effort-
        ms.Eps5*stdMath.Pow(tl.Sets-1,2)/tl.Effort-
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2)/tl.Effort)*tl.Effort/ms.Eps2;
}

func (v volumeBaseSurfacePrediction)InterExerciseFatigue(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    return (1/(tl.Intensity*tl.Intensity)-
        ms.Eps-
        ms.Eps1/tl.Effort-
        ms.Eps2*float64(tl.InterWorkoutFatigue)/tl.Effort-
        ms.Eps4*stdMath.Pow(tl.Sets-1,2)*stdMath.Pow(float64(tl.Reps-1),2)/tl.Effort-
        ms.Eps5*stdMath.Pow(tl.Sets-1,2)/tl.Effort-
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2)/tl.Effort)*tl.Effort/ms.Eps3;
}

func (v volumeBaseSurfacePrediction)Sets(
        ms *db.ModelState,
        tl *db.TrainingLog) float64 {
    return stdMath.Pow(-(ms.Eps+
        ms.Eps1/tl.Effort+
        ms.Eps2*float64(tl.InterWorkoutFatigue)/tl.Effort+
        ms.Eps3*float64(tl.InterExerciseFatigue)/tl.Effort+
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2)/tl.Effort-
        1/(tl.Intensity*tl.Intensity))/(
        ms.Eps4*stdMath.Pow(float64(tl.Reps-1),2)/tl.Effort+
        ms.Eps5/tl.Effort),0.5)+1;
}

func (v volumeBaseSurfacePrediction)Reps(
        ms *db.ModelState,
        tl *db.TrainingLog) float64 {
    return stdMath.Pow(-(ms.Eps+
        ms.Eps1/tl.Effort+
        ms.Eps2*float64(tl.InterWorkoutFatigue)/tl.Effort+
        ms.Eps3*float64(tl.InterExerciseFatigue)/tl.Effort+
        ms.Eps5*stdMath.Pow(tl.Sets-1,2)/tl.Effort-
        1/(tl.Intensity*tl.Intensity))/(
        ms.Eps4*stdMath.Pow(tl.Sets-1,2)/tl.Effort+
        ms.Eps6/tl.Effort),0.5)+1;
}

func (v volumeBaseSurfacePrediction)VolumeSkew(
    ms *db.ModelState,
    tl *db.TrainingLog) float64 {
    return 0;
}

func (v volumeBaseSurfacePrediction)VolumeSkewApprox(
    ms *db.ModelState,
    tl *db.TrainingLog) float64 {
    // Sets/Reps
    return ms.Eps6/ms.Eps5;
}

func (v volumeBaseSurfacePrediction)Stability(ms *db.ModelState) int {
    return 0;
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
            func(vals mathUtil.Vars[float64]) (float64, error) {
                e,err:=vals.Access("E");
                if err!=nil {
                    return 0,err;
                }
                return 1/e,nil;
            }, func(vals mathUtil.Vars[float64]) (float64, error) {
                e,err:=vals.Access("E");
                if err!=nil {
                    return 0,err;
                }
                fw,err:=vals.Access("F_w");
                if err!=nil {
                    return 0,err;
                }
                return fw/e,nil;
            }, func(vals mathUtil.Vars[float64]) (float64, error) {
                e,err:=vals.Access("E");
                if err!=nil {
                    return 0,err;
                }
                fe,err:=vals.Access("F_e");
                if err!=nil {
                    return 0,err;
                }
                return fe/e,nil;
            }, func(vals mathUtil.Vars[float64]) (float64, error) {
                e,err:=vals.Access("E");
                if err!=nil {
                    return 0,err;
                }
                s,err:=vals.Access("S");
                if err!=nil {
                    return 0,err;
                }
                r,err:=vals.Access("R");
                if err!=nil {
                    return 0,err;
                }
                return stdMath.Pow(s-1,2)*stdMath.Pow(r-1,2)/e,nil;
            }, func(vals mathUtil.Vars[float64]) (float64, error) {
                e,err:=vals.Access("E");
                if err!=nil {
                    return 0,err;
                }
                s,err:=vals.Access("S");
                if err!=nil {
                    return 0,err;
                }
                return stdMath.Pow(s-1,2)/e,nil;
            }, func(vals mathUtil.Vars[float64]) (float64, error) {
                e,err:=vals.Access("E");
                if err!=nil {
                    return 0,err;
                }
                r,err:=vals.Access("R");
                if err!=nil {
                    return 0,err;
                }
                return stdMath.Pow(r-1,2)/e,nil;
            }},func(vals mathUtil.Vars[float64]) (float64,error) {
                i,err:=vals.Access("I");
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
func (v *VolumeBaseSurface)Calculations() Calculations { return VolumeBaseSurfacePrediction; }

func (v *VolumeBaseSurface)Update(vals mathUtil.Vars[float64]) error {
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

func (v *VolumeBaseSurface)PredictIntensity(vals mathUtil.Vars[float64]) (float64,error) {
    tmp,err:=v.LinRegResult.Predict(vals);
    if err!=nil {
        return tmp,err;
    }
    return 1/stdMath.Pow(tmp,0.5),err;
}

func (v *VolumeBaseSurface)Stability() int {
    return 0;
}
