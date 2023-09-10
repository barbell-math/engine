package potentialSurface

import "github.com/barbell-math/block/db"

type PotentialSurfaceId int;
// SQL serial values default to starting at 1
const (
    BasicSurfaceId PotentialSurfaceId=iota+1
    VolumeBaseSurfaceId
);

func CalculationsFromSurfaceId(id PotentialSurfaceId) Calculations {
    switch id {
        case BasicSurfaceId: return BasicSurfaceCalculation;
        case VolumeBaseSurfaceId: return VolumeBaseSurfacePrediction;
        default: return nil;
    }
}

type Calculations interface {
    Intensity(ms *db.ModelState, tl *db.TrainingLog) float64;
    OneRMEstimation(ms *db.ModelState) float64;
    Effort(ms *db.ModelState, tl *db.TrainingLog) float64;
    Sets(ms *db.ModelState, tl *db.TrainingLog) float64;
    Reps(ms *db.ModelState, tl *db.TrainingLog) float64;
    InterWorkoutFatigue(ms *db.ModelState, tl *db.TrainingLog) float64;
    InterExerciseFatigue(ms *db.ModelState, tl *db.TrainingLog) float64;
    VolumeSkew(ms *db.ModelState, tl *db.TrainingLog) float64;
    VolumeSkewApprox(ms *db.ModelState, tl *db.TrainingLog) float64;
    Stability(ms *db.ModelState) int;
};

type Surface interface {
    Id() PotentialSurfaceId;
    Calculations() Calculations;
    PredictIntensity(vals map[string]float64) (float64,error);
    Run() (float64,error);
    Update(vals map[string]float64) error;
    GetConstant(idx int) float64;
    Stability() int;
    ToGenericSurf() Surface;
};
