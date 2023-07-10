package stateGenerator

import (
    "time"
    "github.com/barbell-math/block/db"
    logUtil "github.com/barbell-math/block/util/io/log"
)

type StateGeneratorId int;
// SQL serial values default to starting at 1
const (
    SlidingWindowStateGenId StateGeneratorId=iota+1
);

//The struct that holds values when linear regression is performed.
//Note - THE ORDER OF THE STRUCT FIELDS MUST MATCH THE ORDER OF THE VALUES
//IN THE QUERY. Otherwise the values returned will be all jumbled up.
type dataPoint struct {
    DatePerformed time.Time;
    Sets float64;
    Reps float64;
    Effort float64;
    Intensity float64;
    InterExerciseFatigue float64;
    InterWorkoutFatigue float64;
};

//The struct that holds values when searching for missing model states.
//Note - THE ORDER OF THE STRUCT FIELDS MUST MATCH THE ORDER OF THE VALUES
//IN THE QUERY. Otherwise the values returned will be all jumbled up.
type missingModelStateData struct {
    ClientID int;
    ExerciseID int;
    Date time.Time;
};

var SLIDING_WINDOW_DP_DEBUG=logUtil.NewBlankLog[*dataPoint]();
var SLIDING_WINDOW_MS_DEBUG=logUtil.NewBlankLog[db.ModelState]();
var SLIDING_WINDOW_MS_PARALLEL_RESULT_DEBUG=logUtil.NewBlankLog[db.ModelState]();
