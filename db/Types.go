package db;

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
