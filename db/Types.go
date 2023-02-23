package db;

import (
    "time"
)

//Note - THE ORDER OF THE STRUCT FIELDS MUST MATCH
//THE ORDER OF THE VALUES IN THE TABLE DEFINITION.
//If they don't, the values returned will be all
//jumbled up if no errors about parsing values arise.

type DBTable  interface {
    ExerciseType |
    ExerciseFocus |
    Exercise |
    Rotation |
    BodyWeight |
    TrainingLog |
    Client |
    ModelState |
    StateGenerator |
    Prediction
};

type ExerciseType struct {
    Id int;
    T string;
    Description string;
};

type ExerciseFocus struct {
    Id int;
    Focus string;
};

type Exercise struct {
    Id int;
    Name string;
    TypeID int;
    FocusID int;
};

//Start and end dates are inclusive
type Rotation struct {
    Id int;
    ClientID int;
    StartDate time.Time;
    EndDate time.Time;
};

type BodyWeight struct {
    Id int;
    ClientID int;
    Weight float32;
    Date time.Time;
};

type TrainingLog struct {
    Id int;
    ClientID int;
    ExerciseID int;
    RotationID int;
    DatePerformed time.Time;
    Weight float32;
    Sets float32;
    Reps int;
    Intensity float64;
    Effort float64;
    Volume float64;
    InterExerciseFatigue int;
    InterWorkoutFatigue int;
};

type Client struct {
    Id int;
    FirstName string;
    LastName string;
    Email string;
};

type StateGenerator struct {
    Id int;
    T string;
    Description string;
};

type ModelState struct {
    Id int;
    ClientID int;
    ExerciseID int;
    StateGeneratorID int;
    Date time.Time;
    Eps,Eps1,Eps2,Eps3 float64;
    Eps4,Eps5,Eps6,Eps7 float64;
    TimeFrame int;
    Win int;
    Rcond float64;
    Mse float64;
};

type Prediction struct {
    Id int;
    StateGeneratorID int;
    TrainingLogID int;
    IntensityPred float64;
};
