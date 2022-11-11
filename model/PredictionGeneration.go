package model;

import (
    "fmt"
    "github.com/carmichaeljr/powerlifting-engine/db"
)

//Given a set of values to use when making the prediction, the closest model
//state (in time) that is less than the current time will be used to generate
//a prediction for intensity.
//All values besides Intensity need to be accurate in the training log argument.
//'current time' is defined by the 'DatePerformed' field of the training log arg.
func GeneratePrediction(c *db.CRUD, tl *db.TrainingLog) (db.Prediction,error) {
    rv:=db.Prediction{ TrainingLogID: tl.Id };
    query:=nearestExerciseQuery(tl);
    err:=db.CustomReadQuery(c,query,[]any{tl.ExerciseID},func(m *db.ModelState){
        //fmt.Printf("MS in pred: %+v\n",m);
        //fmt.Printf("TL in pred: %+v\n",tl);
        rv.IntensityPred=IntensityPrediction(m,tl);
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
            AND TrainingLog.DatePerformed<'%s'::date
        ORDER BY TrainingLog.DatePerformed DESC,
            TrainingLog.FatigueIndex ASC
        LIMIT 1;`,
        tl.DatePerformed.Format("01/02/2006"),
    );
}
