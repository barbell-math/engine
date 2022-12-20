package model;

import (
    "fmt"
    "github.com/barbell-math/block/db"
)

//Given a set of values to use when making the prediction, the closest model
//state (in time) that is less than the current time will be used to generate
//a prediction for intensity.
//All values besides Intensity need to be accurate in the training log argument.
//'current time' is defined by the 'DatePerformed' field of the training log arg.
func GeneratePrediction(c *db.DB, tl *db.TrainingLog) (db.Prediction,error) {
    rv:=db.Prediction{ TrainingLogID: tl.Id };
    query:=nearestExerciseQuery(tl);
    err:=db.CustomReadQuery(c,query,[]any{
        tl.ExerciseID,
        tl.DatePerformed,
    },func(m *db.ModelState){
        rv.IntensityPred=IntensityPrediction(m,tl);
        rv.StateGeneratorID=m.StateGeneratorID;
    });
    return rv,err;
}

func nearestExerciseQuery(tl *db.TrainingLog) string {
    return fmt.Sprintf(`SELECT ModelState.*
        FROM TrainingLog
        JOIN ModelState
        ON TrainingLog.ExerciseID=ModelState.ExerciseID
            AND TrainingLog.ClientID=ModelState.ClientID
            AND TrainingLog.DatePerformed=ModelState.Date
        WHERE TrainingLog.ExerciseID=$1
            AND TrainingLog.DatePerformed<$2
        ORDER BY TrainingLog.DatePerformed DESC,
            TrainingLog.FatigueIndex ASC
        LIMIT 1;`,
        //tl.DatePerformed.Format("01/02/2006"),
    );
}
