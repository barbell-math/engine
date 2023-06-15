package model;

import (
    "github.com/barbell-math/block/db"
    "github.com/barbell-math/block/util/algo/iter"
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
    if cntr,err:=db.CustomReadQuery[db.ModelState](c,
        nearestModelStateToExerciseQuery(tl),[]any{
            tl.ExerciseID,
            tl.DatePerformed,
            sg.Id,
            tl.ClientID,
    }).Next(func(index int,
        val *db.ModelState,
        status iter.IteratorFeedback,
    ) (iter.IteratorFeedback, *db.ModelState, error) {
        if status!=iter.Break {
            rv.IntensityPred=IntensityPrediction(val,tl);
            rv.StateGeneratorID=val.StateGeneratorID;
        }
        return iter.Continue,val,nil;
    }).Count(); err==nil && cntr>1 {
        return rv,ManyPredictions(
            "Multiple predictions exist for the given training log and state generator.",
        );
    } else {
        return rv,err;
    }
}
