package model;

import (
    "time"
    "github.com/barbell-math/block/db"
)

type StateGenerator interface {
    GenerateModelState(c *db.DB, tl db.TrainingLog, ch chan<- StateGenerationRes);
};

type StateGenerationRes struct {
    Ms db.ModelState;
    Err error;
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
    FatigueIndex float64;
};
