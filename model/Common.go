package model;

import (
    "time"
    "github.com/barbell-math/block/db"
    logUtil "github.com/barbell-math/block/util/io/log"
)

type stateGenerator interface {
    GenerateClientModelStates(d *db.DB, c db.Client, ch chan<- []error);
    GenerateModelState(
        d *db.DB,
        missingData missingModelStateData,
        ch chan<- StateGeneratorRes,
    );
};

type StateGeneratorRes struct {
    Ms db.ModelState;
    Err error;
};

//The struct that holds values when searching for missing model states.
//Note - THE ORDER OF THE STRUCT FIELDS MUST MATCH THE ORDER OF THE VALUES
//IN THE QUERY. Otherwise the values returned will be all jumbled up.
type missingModelStateData struct {
    ClientID int;
    ExerciseID int;
    Date time.Time;
};

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

var SLIDING_WINDOW_DP_DEBUG=logUtil.NewBlankLog[*dataPoint]();
var SLIDING_WINDOW_MS_DEBUG=logUtil.NewBlankLog[db.ModelState]();

//func msMissingQuery(sg db.StateGenerator) string {
//    return `SELECT TrainingLog.DatePerformed,
//        TrainingLog.ExerciseID
//    FROM TrainingLog
//    LEFT JOIN ModelState
//    ON TrainingLog.ExerciseID=ModelState.ExerciseID
//        AND ModelState.ClientID=TrainingLog.ClientID
//        AND TrainingLog.DatePerformed=ModelState.Date
//    JOIN Exercise
//    ON Exercise.Id=TrainingLog.ExerciseID
//    JOIN ExerciseType
//    ON ExerciseType.Id=Exercise.TypeID
//    JOIN
//    WHERE TrainingLog.ClientID=$1
//        AND ModelState.Id IS NULL
//        AND (ExerciseType.T='Main Compound'
//        OR ExerciseType.T='Main Compound Accessory')
//    GROUP BY TrainingLog.DatePerformed,
//        TrainingLog.ExerciseID;`;
//}
