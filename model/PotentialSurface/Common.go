package potentialSurface

import "github.com/barbell-math/block/db"

type PotentialSurfaceId int;
// SQL serial values default to starting at 1
const (
    BasicSurfaceId PotentialSurfaceId=iota+1
    VolumeBaseSurfaceId
);

type Predictor interface {
    Intensity(ms *db.ModelState, tl *db.TrainingLog) float64;
    Effort(ms *db.ModelState, tl *db.TrainingLog) float64;
    Sets(ms *db.ModelState, tl *db.TrainingLog) float64;
    Reps(ms *db.ModelState, tl *db.TrainingLog) float64;
    InterWorkoutFatigue(ms *db.ModelState, tl *db.TrainingLog) float64;
    InterExerciseFatigue(ms *db.ModelState, tl *db.TrainingLog) float64;
};

type Surface interface {
    Id() PotentialSurfaceId;
    Predictor() Predictor;
    PredictIntensity(vals map[string]float64) (float64,error);
    Run() (float64,error);
    Update(vals map[string]float64) error;
    GetConstant(idx int) float64;
};
