package db;

type JoinTable interface {
    ExerciseTypeAndFocus |
    TrainingLogClientAndExerciseAndRotation
};

type ExerciseTypeAndFocus struct {
    Exercise;
    ExerciseType;
    ExerciseFocus;
};

type TrainingLogClientAndExerciseAndRotation struct {
    TrainingLog;
    Client;
    Exercise;
    Rotation;
};
