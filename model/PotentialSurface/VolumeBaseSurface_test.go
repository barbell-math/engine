package potentialSurface

import (
	"testing"

	"github.com/barbell-math/engine/db"
	"github.com/barbell-math/engine/util/test"
)

func TestVolumeBaseSurfaceCreation(t *testing.T){
    cntr:=0;
    m:=NewVolumeBaseSurface();
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

func TestVolumeBaseSurfaceIntensityPrediction(t *testing.T){
    ms:=db.ModelState{
        Eps: 1, Eps1: 0, Eps2: 0, Eps3: 0, Eps4: 0,
        Eps5: 0, Eps6: 0,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 1,
        InterWorkoutFatigue: 0, InterExerciseFatigue: 0,
    };
    res:=VolumeBaseSurfacePrediction.Intensity(&ms,&tl);
    test.BasicTest(float64(1.0),res,
        "Intensity prediction produced incorrect value.",t,
    );
}

func TestVolumeBaseSurfaceEffortPrediction(t *testing.T){
    ms:=db.ModelState{
        Eps: 0, Eps1: 1, Eps2: 0, Eps3: 0,
        Eps4: 0, Eps5: 0, Eps6: 0,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 0,
    };
    res:=VolumeBaseSurfacePrediction.Effort(&ms,&tl);
    test.BasicTest(float64(0.0),res,
        "Effort prediction produced incorrect value.",t,
    );
}

func TestVolumeBaseSurfaceInterWorkoutFatiguePrediction(t *testing.T){
    ms:=db.ModelState{
        Eps: 0, Eps1: 1, Eps2: 1, Eps3: 0,
        Eps4: 0, Eps5: 0, Eps6: 0,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 1, Effort: 1,
    };
    res:=VolumeBaseSurfacePrediction.InterWorkoutFatigue(&ms,&tl);
    test.BasicTest(float64(0.0),res,
        "Inter workout fatigue prediction produced incorrect value.",t,
    );
}

func TestVolumeBaseSurfaceInterExerciseFatiguePrediction(t *testing.T){
    ms:=db.ModelState{
        Eps: 0, Eps1: 1, Eps2: 0, Eps3: 1,
        Eps4: 0, Eps5: 0, Eps6: 0,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 1, Effort: 1,
    };
    res:=VolumeBaseSurfacePrediction.InterExerciseFatigue(&ms,&tl);
    test.BasicTest(float64(0.0),res,
        "Inter exercise fatigue prediction produced incorrect value.",t,
    );
}

func TestVolumeBaseSurfaceSetsPrediction(t *testing.T){
    ms:=db.ModelState{
        Eps: 0, Eps1: 1, Eps2: 0, Eps3: 0,
        Eps4: 0, Eps5: 1, Eps6: 0,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 1, Effort: 1,
    };
    res:=VolumeBaseSurfacePrediction.Sets(&ms,&tl);
    test.BasicTest(float64(1.0),res,
        "Sets prediction produced incorrect value.",t,
    );
}

func TestVolumeBaseSurfaceRepsPrediction(t *testing.T){
    ms:=db.ModelState{
        Eps: 0, Eps1: 1, Eps2: 0, Eps3: 0,
        Eps4: 0, Eps5: 0, Eps6: 1,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 1, Effort: 1,
    };
    res:=VolumeBaseSurfacePrediction.Reps(&ms,&tl);
    test.BasicTest(float64(1.0),res,
        "Reps prediction produced incorrect value.",t,
    );
}
