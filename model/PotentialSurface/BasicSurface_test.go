package potentialSurface

import (
	"testing"

	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/util/test"
)

func TestBasicSurfaceCreation(t *testing.T){
    cntr:=0;
    m:=NewBasicSurface();
    m.IterLHS(func(r int, c int, v float64){
        cntr++;
    });
    test.BasicTest(49,cntr,"LHS Lin reg wrong size for model.",t);
    cntr=0;
    m.IterRHS(func(r int, c int, v float64){
        cntr++;
    });
    test.BasicTest(7,cntr,"RHS Lin reg wrong size for model.",t);
}

func TestBasicSurfaceIntensityPrediction(t *testing.T){
    ms:=db.ModelState{
        Eps: 0, Eps1: 0, Eps2: 0, Eps3: 0, Eps4: 0,
        Eps5: 0, Eps6: 0,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0,
        InterWorkoutFatigue: 0, InterExerciseFatigue: 0,
    };
    res:=BasicSurfacePrediction.Intensity(&ms,&tl);
    test.BasicTest(float64(0.0),res,
        "Intensity prediction produced incorrect value.",t,
    );
}

func TestBasicSurfaceEffortPrediction(t *testing.T){
    //Eps1 has to be 1 to avoid div by 0 error
    ms:=db.ModelState{
        Eps: 0, Eps1: 1, Eps2: 0, Eps3: 0,
        Eps4: 0, Eps5: 0, Eps6: 0,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0,
    };
    res:=BasicSurfacePrediction.Effort(&ms,&tl);
    test.BasicTest(float64(0.0),res,
        "Effort prediction produced incorrect value.",t,
    );
}

func TestBasicSurfaceInterWorkoutFatiguePrediction(t *testing.T){
    //Eps3 has to be 1 to avoid div by 0 error
    ms:=db.ModelState{
        Eps: 0, Eps1: 0, Eps2: 1, Eps3: 0,
        Eps4: 0, Eps5: 0, Eps6: 0,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0,
    };
    res:=BasicSurfacePrediction.InterWorkoutFatigue(&ms,&tl);
    test.BasicTest(float64(0.0),res,
        "Inter workout fatigue prediction produced incorrect value.",t,
    );
}

func TestBasicSurfaceInterExerciseFatiguePrediction(t *testing.T){
    //Eps4 has to be 1 to avoid div by 0 error
    ms:=db.ModelState{
        Eps: 0, Eps1: 0, Eps2: 0, Eps3: 1,
        Eps4: 0, Eps5: 0, Eps6: 0,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0,
    };
    res:=BasicSurfacePrediction.InterExerciseFatigue(&ms,&tl);
    test.BasicTest(float64(0.0),res,
        "Inter exercise fatigue prediction produced incorrect value.",t,
    );
}

func TestBasicSurfaceSetsPrediction(t *testing.T){
    //Eps5 has to be 1 to avoid div by 0 error
    ms:=db.ModelState{
        Eps: 0, Eps1: 0, Eps2: 0, Eps3: 0,
        Eps4: 0, Eps5: 1, Eps6: 0,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0,
    };
    res:=BasicSurfacePrediction.Sets(&ms,&tl);
    test.BasicTest(float64(1.0),res,
        "Sets prediction produced incorrect value.",t,
    );
}

func TestBasicSurfaceRepsPrediction(t *testing.T){
    //Eps7 has to be 1 to avoid div by 0 error
    ms:=db.ModelState{
        Eps: 0, Eps1: 0, Eps2: 0, Eps3: 0,
        Eps4: 0, Eps5: 0, Eps6: 1,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0,
    };
    res:=BasicSurfacePrediction.Reps(&ms,&tl);
    test.BasicTest(float64(1.0),res,
        "Reps prediction produced incorrect value.",t,
    );
}
