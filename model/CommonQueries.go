package model;

import (
    "github.com/barbell-math/engine/db"
)

func nearestModelStateToExerciseQuery(tl *db.TrainingLog) string {
    return `SELECT ModelState.*
        FROM TrainingLog
        JOIN ModelState
        ON TrainingLog.ExerciseID=ModelState.ExerciseID
            AND TrainingLog.ClientID=ModelState.ClientID
            AND TrainingLog.DatePerformed=ModelState.Date
        WHERE TrainingLog.ExerciseID=$1
            AND TrainingLog.DatePerformed<$2
            AND ModelState.StateGeneratorID=$3
            AND ModelState.PotentialSurfaceID=$4
            AND TrainingLog.ClientID=$5
        ORDER BY TrainingLog.DatePerformed DESC
        LIMIT 1;`;
}
