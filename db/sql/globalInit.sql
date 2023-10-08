DROP TABLE IF EXISTS Version CASCADE;
DROP TABLE IF EXISTS TrainingLog CASCADE;
DROP TABLE IF EXISTS Rotation CASCADE;
DROP TABLE IF EXISTS Exercise CASCADE;
DROP TABLE IF EXISTS BodyWeight CASCADE;
DROP TABLE IF EXISTS Client CASCADE;
DROP TABLE IF EXISTS ExerciseType CASCADE;
DROP TABLE IF EXISTS ExerciseFocus CASCADE;
DROP TABLE IF EXISTS ModelState CASCADE;
DROP TABLE IF EXISTS Prediction CASCADE;
DROP TABLE IF EXISTS StateGenerator CASCADE;
DROP TABLE IF EXISTS PotentialSurface CASCADE;

CREATE TABLE IF NOT EXISTS Version (
    Num INT NOT NULL
);

CREATE TABLE Client (
	Id SERIAL PRIMARY KEY,
	FirstName TEXT NOT NULL,
	LastName TEXT NOT NULL,
	Email TEXT NOT NULL UNIQUE
);

CREATE TABLE ExerciseType (
	Id SERIAL PRIMARY KEY,
	T TEXT NOT NULL UNIQUE,
	Description TEXT NOT NULL
);

CREATE TABLE ExerciseFocus (
	Id SERIAL PRIMARY KEY,
	Focus TEXT NOT NULL UNIQUE
);

CREATE TABLE Exercise (
	Id SERIAL PRIMARY KEY,
	Name TEXT NOT NULL UNIQUE,
	TypeID INT NOT NULL,
	FocusID INT NOT NULL,
    FOREIGN KEY (typeID) REFERENCES ExerciseType(Id),
    FOREIGN KEY (focusID) REFERENCES ExerciseFocus(Id)
);

CREATE TABLE Rotation (
	Id SERIAL PRIMARY KEY,
	ClientID INTEGER NOT NULL,
	StartDate DATE NOT NULL,
	EndDate DATE NOT NULL,
	FOREIGN KEY (ClientID) REFERENCES Client(Id)
);

CREATE TABLE BodyWeight (
    Id SERIAL PRIMARY KEY,
	ClientID INTEGER NOT NULL,
	Weight FLOAT NOT NULL,
    Date DATE NOT NULL,
	FOREIGN KEY (ClientID) REFERENCES Client(Id)
);

CREATE TABLE TrainingLog (
    Id SERIAL PRIMARY KEY,
	ClientID INTEGER NOT NULL,
	ExerciseID INTEGER NOT NULL,
    RotationID INTEGER NOT NULL,
	DatePerformed DATE NOT NULL DEFAULT CURRENT_DATE,
	Weight FLOAT NOT NULL,
	Sets FLOAT NOT NULL,
	Reps SMALLINT NOT NULL,
	Intensity FLOAT,
    Effort FLOAT,
    Volume FLOAT NOT NULL,
    InterExerciseFatigue INT NOT NULL,
    InterWorkoutFatigue INT NOT NULL,
	FOREIGN KEY (ClientID) REFERENCES Client(ID),
	FOREIGN KEY (ExerciseID) REFERENCES Exercise(ID),
	FOREIGN KEY (RotationID) REFERENCES Rotation(ID)
);

CREATE TABLE PotentialSurface (
    Id SERIAL PRIMARY KEY,
	T TEXT NOT NULL UNIQUE,
	Description TEXT NOT NULL
);

CREATE TABLE StateGenerator (
    Id SERIAL PRIMARY KEY,
	T TEXT NOT NULL UNIQUE,
	Description TEXT NOT NULL
);

CREATE TABLE ModelState (
    Id SERIAL PRIMARY KEY,
    ClientID INTEGER NOT NULL,
    ExerciseID INTEGER NOT NULL,
    PotentialSurfaceID INTEGER NOT NULL,
    StateGeneratorID INTEGER NOT NULL,
    Date DATE NOT NULL,
    Eps FLOAT NOT NULL,
    Eps1 FLOAT NOT NULL,
    Eps2 FLOAT NOT NULL,
    Eps3 FLOAT NOT NULL,
    Eps4 FLOAT NOT NULL,
    Eps5 FLOAT NOT NULL,
    Eps6 FLOAT NOT NULL,
    Eps7 FLOAT NOT NULL,
    TimeFrame INTEGER NOT NULL,
    Win INTEGER NOT NULL,
    Rcond FLOAT NOT NULL,
    Mse FLOAT NOT NULL,
    FOREIGN KEY (ClientID) REFERENCES Client(Id),
    FOREIGN KEY (ExerciseID) REFERENCES Exercise(Id),
    FOREIGN KEY (StateGeneratorID) REFERENCES StateGenerator(Id),
    FOREIGN KEY (PotentialSurfaceID) REFERENCES PotentialSurface(Id)
);

CREATE TABLE Prediction (
    Id SERIAL PRIMARY KEY,
    PotentialSurfaceID INTEGER NOT NULL,
    StateGeneratorID INTEGER NOT NULL,
    TrainingLogID INTEGER NOT NULL,
    IntensityPred FLOAT NOT NULL,
    FOREIGN KEY (TrainingLogID) REFERENCES TrainingLog(Id),
    FOREIGN KEY (StateGeneratorID) REFERENCES StateGenerator(Id)
);

ALTER TABLE ModelState
ADD CONSTRAINT uniqueDayExerciseClientState
UNIQUE(ClientID,ExerciseID,StateGeneratorID,PotentialSurfaceID,Date);

ALTER TABLE Prediction
ADD CONSTRAINT uniqueGeneratorTrainingLogID
UNIQUE(StateGeneratorID,PotentialSurfaceID,TrainingLogID);

INSERT INTO Version(num) VALUES (0);
