package model;

import (
    "time"
    "github.com/carmichaeljr/powerlifting-engine/db"
)

type StateGenerator interface {
    UpdateModelStates(c *db.CRUD, clientID int, rv chan<- []error);
    GenerateModelState(c *db.CRUD, tl db.TrainingLog, ch chan<- StateGenerationRes);
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
