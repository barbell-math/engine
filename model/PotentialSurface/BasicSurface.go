package potentialSurface

import (
	stdMath "math"

	"github.com/barbell-math/engine/db"
	"github.com/barbell-math/engine/util/dataStruct"
	mathUtil "github.com/barbell-math/engine/util/math/numeric"
)

//The basic surface follows the following equation:
//  I=eps+eps_1*E-( eps_2*F_w+eps_3*F_e+eps_4*(s-1)^2(r-1)^2+eps_5*(s-1)^2+eps_6*(r-1)^2 )
//This equation does not take into account latent fatigue, which makes it naive
//because it does not consider the relationship between lifts across time.
var BasicSurfaceCalculation basicSurfaceCalculation;
type basicSurfaceCalculation struct { };

func (b basicSurfaceCalculation)Intensity(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    return (ms.Eps+
        ms.Eps1*float64(tl.Effort)-
        ms.Eps2*float64(tl.InterWorkoutFatigue)-
        ms.Eps3*float64(tl.InterExerciseFatigue)-
        ms.Eps4*stdMath.Pow(tl.Sets-1,2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps5*stdMath.Pow(tl.Sets-1,2)-
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2));
}

func (b basicSurfaceCalculation)OneRMEstimation(ms *db.ModelState) float64 {
    tl:=db.TrainingLog{
        Sets: 1, Reps: 1,
        Effort: 10,
        InterWorkoutFatigue: 0,
        InterExerciseFatigue: 0,
    };
    return b.Intensity(ms,&tl);
}

func (b basicSurfaceCalculation)Effort(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    return (tl.Intensity-ms.Eps+
        ms.Eps2*float64(tl.InterWorkoutFatigue)+
        ms.Eps3*float64(tl.InterExerciseFatigue)+
        ms.Eps4*stdMath.Pow(tl.Sets-1,2)*stdMath.Pow(float64(tl.Reps-1),2)+
        ms.Eps5*stdMath.Pow(tl.Sets-1,2)+
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2))/ms.Eps1;
}

func (b basicSurfaceCalculation)InterWorkoutFatigue(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    return (ms.Eps+
        ms.Eps1*tl.Effort-
        ms.Eps3*float64(tl.InterExerciseFatigue)-
        ms.Eps4*stdMath.Pow(tl.Sets-1,2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps5*stdMath.Pow(tl.Sets-1,2)-
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2)-
        tl.Intensity)/ms.Eps2;
}

func (b basicSurfaceCalculation)InterExerciseFatigue(
        ms *db.ModelState, 
        tl *db.TrainingLog) float64 {
    return (ms.Eps+
        ms.Eps1*tl.Effort-
        ms.Eps2*float64(tl.InterWorkoutFatigue)-
        ms.Eps4*stdMath.Pow(tl.Sets-1,2)*stdMath.Pow(float64(tl.Reps-1),2)-
        ms.Eps5*stdMath.Pow(tl.Sets-1,2)-
        ms.Eps6*stdMath.Pow(float64(tl.Reps-1),2)-
        tl.Intensity)/ms.Eps3;
}

func (b basicSurfaceCalculation)Sets(
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

func (b basicSurfaceCalculation)Reps(
        ms *db.ModelState,
        tl *db.TrainingLog) float64 {
    return stdMath.Pow((
        ms.Eps+
        ms.Eps1*tl.Effort-
        ms.Eps2*float64(tl.InterWorkoutFatigue)-
        ms.Eps3*float64(tl.InterExerciseFatigue)-
        ms.Eps5*stdMath.Pow(tl.Sets-1,2)-
        tl.Intensity)/(
        ms.Eps4*stdMath.Pow(tl.Sets-1,2)+
        ms.Eps6),0.5)+1.0;
}

func (b basicSurfaceCalculation)VolumeSkew(
    ms *db.ModelState,
    tl *db.TrainingLog) float64 {
    return (b.volumeSkewIntegral1(ms,tl)+
        b.volumeSkewIntegral2(ms,tl))/(
        b.volumeSkewIntegral3(ms,tl)+
        b.volumeSkewIntegral4(ms,tl));
}

func (b basicSurfaceCalculation)volumeSkewDiagonal(
    ms *db.ModelState,
    tl *db.TrainingLog) float64 {
    return stdMath.Pow(stdMath.Max((-ms.Eps5-ms.Eps6+
        stdMath.Pow(stdMath.Max(stdMath.Pow(ms.Eps5+ms.Eps6,2)+
            4*ms.Eps4*(
                ms.Eps+
                ms.Eps1*tl.Effort-
                ms.Eps2*float64(tl.InterWorkoutFatigue)-
                ms.Eps3*float64(tl.InterExerciseFatigue)),0),0.5))/(2*ms.Eps4),0),0.5)+1;
}

func (b basicSurfaceCalculation)setsWhenRepsEquals1(
    ms *db.ModelState,
    tl *db.TrainingLog) float64 {
    return stdMath.Pow(stdMath.Max((ms.Eps+ms.Eps1*tl.Effort-
        ms.Eps2*float64(tl.InterWorkoutFatigue)-
        ms.Eps3*float64(tl.InterExerciseFatigue))/ms.Eps5,0),0.5)+1;
}

func (b basicSurfaceCalculation)repsWhenSetsEquals1(
    ms *db.ModelState,
    tl *db.TrainingLog) float64 {
    return stdMath.Pow(stdMath.Max((ms.Eps+ms.Eps1*tl.Effort-
        ms.Eps2*float64(tl.InterWorkoutFatigue)-
        ms.Eps3*float64(tl.InterExerciseFatigue))/ms.Eps6,0),0.5)+1;
}

func (b basicSurfaceCalculation)volumeSkewIntegral1(
    ms *db.ModelState,
    tl *db.TrainingLog) float64 {
    // Note - the order of the params is VERY important. It needs to correlate
    // with the order of the params for the double integral. The only difference
    // between this and integral 3 is the ordering of these values!
    f:=func(s float64, r float64) float64 { 
        tl.Sets=s;
        tl.Reps=r;
        return s*r*b.Intensity(ms,tl); 
    }
    rv,_:=mathUtil.DoubleIntegral(f)(
        1,
        b.volumeSkewDiagonal(ms,tl),
        mathUtil.ConstIntegralBound[float64](1),
        func(s float64) float64 { return s; },
        1201,
    );
    return rv;
}

func (b basicSurfaceCalculation)volumeSkewIntegral2(
    ms *db.ModelState,
    tl *db.TrainingLog) float64 {
    // Note - the order of the params is VERY important. It needs to correlate
    // with the order of the params for the double integral. The only difference
    // between this and integral 4 is the ordering of these values as well as
    // the switching of the Eps5 and Eps6 constants.
    f:=func(s float64, r float64) float64 { 
        tl.Sets=s;
        tl.Reps=r;
        return s*r*b.Intensity(ms,tl); 
    }
    rv,_:=mathUtil.DoubleIntegral(f)(
        b.volumeSkewDiagonal(ms,tl),
        b.setsWhenRepsEquals1(ms,tl),
        mathUtil.ConstIntegralBound[float64](1),
        func(s float64) float64 {
            return stdMath.Pow(stdMath.Max((
                ms.Eps+ms.Eps1*tl.Effort-
                ms.Eps2*float64(tl.InterWorkoutFatigue)-
                ms.Eps3*float64(tl.InterExerciseFatigue)-
                ms.Eps5*stdMath.Pow(s-1,2))/(
                ms.Eps4*stdMath.Pow(s-1,2)+
                ms.Eps6),0),0.5)+1;
        },
        1201,
    );
    return rv;
}

func (b basicSurfaceCalculation)volumeSkewIntegral3(
    ms *db.ModelState,
    tl *db.TrainingLog) float64 {
    // Note - the order of the params is VERY important. It needs to correlate
    // with the order of the params for the double integral. The only difference
    // between this and integral 1 is the ordering of these values!
    f:=func(r float64, s float64) float64 { 
        tl.Sets=s;
        tl.Reps=r;
        return s*r*b.Intensity(ms,tl); 
    }
    rv,_:=mathUtil.DoubleIntegral(f)(
        1,
        b.volumeSkewDiagonal(ms,tl),
        mathUtil.ConstIntegralBound[float64](1),
        func(r float64) float64 { return r; },
        1201,
    );
    return rv;
}

func (b basicSurfaceCalculation)volumeSkewIntegral4(
    ms *db.ModelState,
    tl *db.TrainingLog) float64 {
    // Note - the order of the params is VERY important. It needs to correlate
    // with the order of the params for the double integral. The only difference
    // between this and integral 4 is the ordering of these values as well as
    // the switching of the Eps5 and Eps6 constants.
    f:=func(r float64, s float64) float64 { 
        tl.Sets=s;
        tl.Reps=r;
        return s*r*b.Intensity(ms,tl); 
    }
    rv,_:=mathUtil.DoubleIntegral(f)(
        b.volumeSkewDiagonal(ms,tl),
        b.repsWhenSetsEquals1(ms,tl),
        mathUtil.ConstIntegralBound[float64](1),
        func(r float64) float64 {
            return stdMath.Pow(stdMath.Max((
                ms.Eps+ms.Eps1*tl.Effort-
                ms.Eps2*float64(tl.InterWorkoutFatigue)-
                ms.Eps3*float64(tl.InterExerciseFatigue)-
                ms.Eps6*stdMath.Pow(r-1,2))/(
                ms.Eps4*stdMath.Pow(r-1,2)+
                ms.Eps5),0),0.5)+1;
        },
        1201,
    );
    return rv;
}

func (b basicSurfaceCalculation)VolumeSkewApprox(
    ms *db.ModelState,
    tl *db.TrainingLog) float64 {
    return (stdMath.Pow(ms.Eps6,0.5)*(
        stdMath.Pow(
            ms.Eps+
            ms.Eps1*tl.Effort-
            ms.Eps2*float64(tl.InterWorkoutFatigue)-
            ms.Eps3*float64(tl.InterExerciseFatigue),
            0.5,
        )+stdMath.Pow(ms.Eps5,0.5)))/(stdMath.Pow(ms.Eps5,0.5)*(
        stdMath.Pow(
            ms.Eps+
            ms.Eps1*tl.Effort-
            ms.Eps2*float64(tl.InterWorkoutFatigue)-
            ms.Eps3*float64(tl.InterExerciseFatigue),
            0.5,
        )+stdMath.Pow(ms.Eps6,0.5)));
}

func (b basicSurfaceCalculation)Stability(ms *db.ModelState) int {
    rv:=0;
    if ms.Eps>0 {
        rv++;
    }
    if ms.Eps1>0 {
        rv++;
    }
    if ms.Eps2>0 {
        rv++;
    }
    if ms.Eps3>0 {
        rv++;
    }
    if ms.Eps4>0 {
        rv++;
    }
    if ms.Eps5>0 {
        rv++;
    }
    if ms.Eps6>0 {
        rv++;
    }
    if ms.Eps7>0 {
        rv++;
    }
    return rv;
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

func (b BasicSurface)ToGenericSurf() Surface { return &b; }

func (b *BasicSurface)Id() PotentialSurfaceId { return BasicSurfaceId; }
func (b *BasicSurface)Calculations() Calculations { return BasicSurfaceCalculation; }

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
        {A: 0, B: stdMath.Inf(1)}, //Eps1: Effort
        {A: 0, B: stdMath.Inf(1)}, //Eps2: F_w
        {A: 0, B: stdMath.Inf(1)}, //Eps3: F_e
        {A: 0, B: stdMath.Inf(1)}, //Eps4: s*r
        {A: 0, B: stdMath.Inf(1)}, //Eps5: s
        {A: 0, B: stdMath.Inf(1)}, //Eps6: r
    };
    for i,v:=range(b.LinRegResult.V) {
        b.LinRegResult.V[i][0]=mathUtil.Constrain(v[0],constraints[i]);
    }
}

func (b *BasicSurface)PredictIntensity(vals map[string]float64) (float64,error) {
    return b.LinRegResult.Predict(vals);
}

func (b *BasicSurface)Stability() int {
    rv:=0;
    for _,v:=range(b.LinRegResult.V) {
        if v[0]>0 {
            rv++;
        }
    }
    return rv;
}
