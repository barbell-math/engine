package potentialSurface

import (
	"fmt"

	"github.com/barbell-math/engine/db"
	mathUtil "github.com/barbell-math/engine/util/math/numeric"
)

type PotentialSurfaceId int;
// SQL serial values default to starting at 1
const (
    BasicSurfaceId PotentialSurfaceId=iota+1
    VolumeBaseSurfaceId
);

func CalculationsFromSurfaceId(id PotentialSurfaceId) (Calculations,error) {
    switch id {
        case BasicSurfaceId: return BasicSurfaceCalculation,nil;
        case VolumeBaseSurfaceId: return VolumeBaseSurfacePrediction,nil;
        default: return nil,InvalidPotentialSurfaceId(
            fmt.Sprintf("Id: %d",id),
        );
    }
}

type Calculations interface {
    Intensity(ms *db.ModelState, tl *db.TrainingLog) float64;
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
    PredictIntensity(vals mathUtil.Vars[float64]) (float64,error);
    Run() (float64,error);
    Update(vals mathUtil.Vars[float64]) error;
    GetConstant(idx int) float64;
    Stability() int;
    ToGenericSurf() Surface;
};
