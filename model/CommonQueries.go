package model;

import (
    "github.com/barbell-math/block/db"
)

func timeFrameQuery() string {
    return `SELECT DatePerformed,
            Sets, Reps, Effort, Intensity,
            InterExerciseFatigue, InterWorkoutFatigue
        FROM TrainingLog
        WHERE TrainingLog.DatePerformed<$1
            AND TrainingLog.DatePerformed>$2
            AND TrainingLog.ExerciseID=$3
            AND TrainingLog.ClientID=$4
        ORDER BY DatePerformed DESC;`;
}

func missingModelStatesForGivenStateGenQuery() string {
    return `SELECT newTl.ClientID, newTl.ExerciseID, newTl.DatePerformed
        FROM (SELECT *
            FROM TrainingLog
            WHERE TrainingLog.ClientID=$1
                AND TrainingLog.DatePerformed>$3
        ) newTl
        LEFT JOIN (SELECT *
            FROM ModelState
            WHERE ModelState.ClientID=$1
                AND ModelState.StateGeneratorID=$2
                AND ModelState.Date>$3
        ) newMs
        ON newMs.Date=newTl.DatePerformed
            AND newMs.ExerciseID=newTl.ExerciseID
        JOIN Exercise
        ON Exercise.Id=newTl.ExerciseID
        JOIN ExerciseType
        ON Exercise.TypeID=ExerciseType.ID
        WHERE newMs.Id IS NULL
            AND (ExerciseType.T='Main Compound'
                OR ExerciseType.T='Main Compound Accessory'
        ) GROUP BY newTl.DatePerformed, newTl.ExerciseID, newTl.ClientID;`;
}

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
            AND TrainingLog.ClientID=$4
        ORDER BY TrainingLog.DatePerformed DESC,
            TrainingLog.InterWorkoutFatigue ASC;`;
}
