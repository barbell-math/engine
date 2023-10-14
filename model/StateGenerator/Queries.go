package stateGenerator

func timeFrameQuery() string {
    return `SELECT DatePerformed,
            Sets, Reps, Effort, Intensity,
            InterExerciseFatigue, InterWorkoutFatigue
        FROM TrainingLog
        WHERE TrainingLog.DatePerformed<=$1
            AND TrainingLog.DatePerformed>$2
            AND TrainingLog.ExerciseID=$3
            AND TrainingLog.ClientID=$4
        ORDER BY 
            DatePerformed DESC,
            Id ASC;`;
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

