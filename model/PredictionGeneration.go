package model;

import (
    "github.com/barbell-math/block/db"
)

//Given a set of values to use when making the prediction, the closest model
//state (in time) that is less than the current time and has the appropriate
//state generator will be used to generate a prediction for intensity.
//All values besides Intensity need to be accurate in the training log argument.
//'current time' is defined by the 'DatePerformed' field of the training log arg.
func GeneratePrediction(
        c *db.DB,
        tl *db.TrainingLog,
        sg *db.StateGenerator) (db.Prediction,error) {
    rv:=db.Prediction{ TrainingLogID: tl.Id };
    cntr:=0;
    if err:=db.CustomReadQuery(c,nearestModelStateToExerciseQuery(tl),[]any{
        tl.ExerciseID,
        tl.DatePerformed,
        sg.Id,
        tl.ClientID,
    }, func(m *db.ModelState){
        rv.IntensityPred=IntensityPrediction(m,tl);
        rv.StateGeneratorID=m.StateGeneratorID;
        cntr++;
    }); err==nil && cntr>1 {
        return rv,ManyPredictions(
            "Multiple predictions exist for the given training log and state generator.",
        );
    } else {
        return rv,err;
    }
}
