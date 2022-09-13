package db;

import (
    "time"
)

type ExerciseType struct {
    id int;
    _type string;
    description string;
};

type ExerciseFocus struct {
    id int;
    focus string;
};

type Exercise struct {
    id int;
    name string;
    typeID int;
    focusID int;
};

type Rotation struct {
    id int;
    userID int;
    startDate time.Time;
    endData time.Time;
};

type BodyWeight struct {
    userID int;
    weight float64;
};

type TrainingLog struct {
    userID int;
    exerciseID int;
    datePerformed time.Time;
    weight float64;
    sets int;
    reps int;
    intensity float64;
};

type Client struct {
    id int;
    fName string;
    lName string;
    email string;
};
