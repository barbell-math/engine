package potentialSurface

import (
	"testing"
    stdMath "math"

	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/util/test"
)

//All seemingly magic numbers come from: https://www.desmos.com/calculator/9zxa3zuum0

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
        Eps: 5, Eps1: 5, Eps2: 1, Eps3: 1,
        Eps4: 2, Eps5: 1, Eps6: 0.5,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 2, Reps: 2, Intensity: 0, Effort: 10,
        InterWorkoutFatigue: 1, InterExerciseFatigue: 1,
    };
    res:=BasicSurfacePrediction.Intensity(&ms,&tl);
    test.BasicTest(float64(49.5),res,
        "Intensity prediction produced incorrect value.",t,
    );
}

func TestBasicSurfaceEffortPrediction(t *testing.T){
    ms:=db.ModelState{
        Eps: 5, Eps1: 5, Eps2: 1, Eps3: 1,
        Eps4: 2, Eps5: 1, Eps6: 0.5,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 2, Reps: 2, Intensity: 0, Effort: 10,
        InterWorkoutFatigue: 1, InterExerciseFatigue: 1,
    };
    res:=BasicSurfacePrediction.Effort(&ms,&tl);
    test.BasicTest(float64(0.1),res,
        "Effort prediction produced incorrect value.",t,
    );
}

func TestBasicSurfaceInterWorkoutFatiguePrediction(t *testing.T){
    ms:=db.ModelState{
        Eps: 5, Eps1: 5, Eps2: 1, Eps3: 1,
        Eps4: 2, Eps5: 1, Eps6: 0.5,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 2, Reps: 2, Intensity: 0, Effort: 10,
        InterWorkoutFatigue: 1, InterExerciseFatigue: 1,
    };
    res:=BasicSurfacePrediction.InterWorkoutFatigue(&ms,&tl);
    test.BasicTest(float64(50.5),res,
        "Inter workout fatigue prediction produced incorrect value.",t,
    );
}

func TestBasicSurfaceInterExerciseFatiguePrediction(t *testing.T){
    ms:=db.ModelState{
        Eps: 5, Eps1: 5, Eps2: 1, Eps3: 1,
        Eps4: 2, Eps5: 1, Eps6: 0.5,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 2, Reps: 2, Intensity: 0, Effort: 10,
        InterWorkoutFatigue: 1, InterExerciseFatigue: 1,
    };
    res:=BasicSurfacePrediction.InterExerciseFatigue(&ms,&tl);
    test.BasicTest(float64(50.5),res,
        "Inter exercise fatigue prediction produced incorrect value.",t,
    );
}

func TestBasicSurfaceSetsPrediction(t *testing.T){
    ms:=db.ModelState{
        Eps: 5, Eps1: 5, Eps2: 1, Eps3: 1,
        Eps4: 2, Eps5: 1, Eps6: 0.5,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 2, Reps: 2, Intensity: 0, Effort: 10,
        InterWorkoutFatigue: 1, InterExerciseFatigue: 1,
    };
    res:=BasicSurfacePrediction.Sets(&ms,&tl);
    test.BasicTest(true,(float64(5.18330013267)-res)<1e-6,
        "Sets prediction produced incorrect value.",t,
    );
}

func TestBasicSurfaceRepsPrediction(t *testing.T){
    ms:=db.ModelState{
        Eps: 5, Eps1: 5, Eps2: 1, Eps3: 1,
        Eps4: 2, Eps5: 1, Eps6: 0.5,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 2, Reps: 2, Intensity: 0, Effort: 10,
        InterWorkoutFatigue: 1, InterExerciseFatigue: 1,
    };
    res:=BasicSurfacePrediction.Reps(&ms,&tl);
    test.BasicTest(true,(float64(5.5607017004)-res)<1e-6,
        "Reps prediction produced incorrect value.",t,
    );
}

func TestBasicSurfaceVolumeSkewSymmetrical(t *testing.T){
    ms:=db.ModelState{
        Eps: 5, Eps1: 5, Eps2: 1, Eps3: 1,
        Eps4: 2, Eps5: 1, Eps6: 1,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 10,
        InterWorkoutFatigue: 1, InterExerciseFatigue: 1,
    };
    test.BasicTest(true,stdMath.Abs(
        float64(3.16149027673)-
        BasicSurfacePrediction.volumeSkewDiagonal(&ms,&tl))<1e-6,
        "The volume skew diagonal computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(418.500546953)-
        BasicSurfacePrediction.volumeSkewIntegral1(&ms,&tl))<1e-6,
        "The first integral value computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(418.500546953)-
        BasicSurfacePrediction.volumeSkewIntegral3(&ms,&tl))<1e-6,
        "The third integral value computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(750.548815622)-
        BasicSurfacePrediction.volumeSkewIntegral2(&ms,&tl))<1e-6,
        "The second integral value computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(750.548815622)-
        BasicSurfacePrediction.volumeSkewIntegral4(&ms,&tl))<1e-6,
        "The fourth integral value computation was not correct.",t,
    );
    test.BasicTest(float64(1),BasicSurfacePrediction.VolumeSkew(&ms,&tl),
        "The volume skew result was not correct.",t,
    );
}

func TestBasicSurfaceVolumeSkewSets(t *testing.T){
    ms:=db.ModelState{
        Eps: 5, Eps1: 5, Eps2: 1, Eps3: 1,
        Eps4: 2, Eps5: 0.5, Eps6: 1,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 10,
        InterWorkoutFatigue: 1, InterExerciseFatigue: 1,
    };
    test.BasicTest(true,stdMath.Abs(
        float64(3.18779699824)-
        BasicSurfacePrediction.volumeSkewDiagonal(&ms,&tl))<1e-6,
        "The volume skew diagonal computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(443.872115854)-
        BasicSurfacePrediction.volumeSkewIntegral1(&ms,&tl))<1e-6,
        "The first integral value computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(435.08379706)-
        BasicSurfacePrediction.volumeSkewIntegral3(&ms,&tl))<1e-6,
        "The third integral value computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(1234.13800052)-
        BasicSurfacePrediction.volumeSkewIntegral2(&ms,&tl))<1e-6,
        "The second integral value computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(751.943070387)-
        BasicSurfacePrediction.volumeSkewIntegral4(&ms,&tl))<1e-6,
        "The fourth integral value computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(1.41362437734)-
        BasicSurfacePrediction.VolumeSkew(&ms,&tl))<1e-6,
        "The volume skew result was not correct.",t,
    );
}

func TestBasicSurfaceVolumeSkewReps(t *testing.T){
    ms:=db.ModelState{
        Eps: 5, Eps1: 5, Eps2: 1, Eps3: 1,
        Eps4: 2, Eps5: 1, Eps6: 0.5,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 10,
        InterWorkoutFatigue: 1, InterExerciseFatigue: 1,
    };
    test.BasicTest(true,stdMath.Abs(
        float64(3.18779699824)-
        BasicSurfacePrediction.volumeSkewDiagonal(&ms,&tl))<1e-6,
        "The volume skew diagonal computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(435.08379706)-
        BasicSurfacePrediction.volumeSkewIntegral1(&ms,&tl))<1e-6,
        "The first integral value computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(443.872115854)-
        BasicSurfacePrediction.volumeSkewIntegral3(&ms,&tl))<1e-6,
        "The third integral value computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(751.943070387)-
        BasicSurfacePrediction.volumeSkewIntegral2(&ms,&tl))<1e-6,
        "The second integral value computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(1234.13800052)-
        BasicSurfacePrediction.volumeSkewIntegral4(&ms,&tl))<1e-6,
        "The fourth integral value computation was not correct.",t,
    );
    test.BasicTest(true,stdMath.Abs(
        float64(0.707401496487)-
        BasicSurfacePrediction.VolumeSkew(&ms,&tl))<1e-6,
        "The volume skew result was not correct.",t,
    );
}

func BenchmarkBasicSurfaceVolumeSkew(b *testing.B) {
    ms:=db.ModelState{
        Eps: 5, Eps1: 5, Eps2: 1, Eps3: 1,
        Eps4: 2, Eps5: 1, Eps6: 0.5,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 0, Effort: 10,
        InterWorkoutFatigue: 1, InterExerciseFatigue: 1,
    };
    for i:=0; i<b.N; i++ {
        BasicSurfacePrediction.VolumeSkew(&ms,&tl);
    }
}
