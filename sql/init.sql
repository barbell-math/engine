--https://sqliteonline.com/
.headers on
.mode column

DROP TABLE ExerciseTypes;
DROP TABLE Exercises;
DROP TABLE Rotations;

CREATE TABLE ExerciseTypes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	_type TEXT NOT NULL,
	description TEXT NOT NULL
);

CREATE TABLE Exercises (
	id INTEGER PRIMARY KEY AUTOINCREMENT ,
	name TEXT NOT NULL,
	_type INT NOT NULL,
	barbell BOOLEAN DEFAULT FALSE,
	dumbbell BOOLEAN DEFAULT FALSE,
	FOREIGN KEY (_type) REFERENCES ExerciseTypes(id)
);

CREATE TABLE Rotations (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	_start TEXT NOT NULL,
	_end TEXT NOT NULL,
	squatMax FLOAT DEFAULT 0.0,
	benchMax FLOAT DEFAULT 0.0,
	deadliftMax DEFAULT 0.0
);

INSERT INTO ExerciseTypes(_type,description) 
VALUES ('Main Compound','The squat, bench, and deadlift.');
INSERT INTO ExerciseTypes(_type,description) 
VALUES ('Main Compound Accessory','Variations of the squat, bench, and deadlift that do not significantly change the mechanics of the lift itself.');
INSERT INTO ExerciseTypes(_type,description)
VALUES ('Compound Accessory','Multi-joint accessories that are not part of the main compound accessory group.');
INSERT INTO ExerciseTypes(_type,description) 
VALUES ('Accessory','Single joint lifts and core work.');

INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Banded Deadlift',2,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Banded Squat',2,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Bench',1,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Block Press',2,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Bulgarian Split Squat',3,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Cable Crunches',4,FALSE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Cable Row',4,FALSE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Deadbugs',4,FALSE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Deadlift',1,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Deficit Romanian Deadlift',2,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Dumbbell Lateral Side Raise',4,FALSE,TRUE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Dumbbell RDL',3,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Flat Dumbbell Press',3,FALSE,TRUE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Hamstring Curls',4,FALSE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Heeled Squat',2,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Incline Bench',3,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Incline Dumbbell Press',3,FALSE,TRUE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Larson Press',3,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Lat Pulldown',4,FALSE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Narrow Grip Bench',2,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Paused Deadlift',2,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Paused Squat',2,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Plank',4,FALSE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Pullup',3,FALSE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Pushup',3,FALSE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Romanian Deadlift',3,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Saftey Bar Squat',2,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Seated Dumbbell Overhead Press',3,FALSE,TRUE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Skull Crusher',4,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Squat',1,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Tempo Squat',2,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('Tricep Pushdowns',4,FALSE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('V Bar Rows',3,TRUE,FALSE);
INSERT INTO Exercises(name,_type,barbell,dumbbell) VALUES ('YTW''s',4,FALSE,FALSE);

INSERT INTO Rotations(id,_start,_end,squatmax,benchmax,deadliftmax)
VALUES (0,DATETIME('2021-08-25'),DATETIME('2021-08-26'),455,275,500);
INSERT INTO Rotations(_start,_end,squatmax,benchmax,deadliftmax)
VALUES (DATETIME('2021-08-27'),DATETIME('2021-11-15'),425,275,505);

SELECT * FROM ExerciseTypes;
SELECT * FROM Exercises;
SELECT * FROM Rotations;
