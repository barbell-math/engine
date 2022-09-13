INSERT INTO Client(firstName,lastName,email)
VALUES ('Jack','Carmichael','test@outlook.com');

--A users initial rotation must always start with id 0
--This forms the original base.
INSERT INTO Rotations(id,userID,startDate,endDate)
VALUES (0,1,DATETIME('2021-08-25'),DATETIME('2021-08-26'));
INSERT INTO Rotations(userID,startDate,endDate)
VALUES (1,DATETIME('2021-08-27'),DATETIME('2021-11-15'));
INSERT INTO Rotations(userID,startDate,endDate)
VALUES (1,DATETIME('2021-11-16'),DATETIME('now','localtime'));

INSERT INTO Maxes (userID,roationID,exerciseID,weight)
VALUES (
	SELECT id FROM Users WHERE firstName='Jack',
	,0,
	SELECT exerciseID FROM Exercises WHERE name='Squat',
	455
);

SELECT * FROM USERS;
SELECT * FROM ExerciseTypes;
SELECT * FROM Exercises;
SELECT * FROM Rotations;
SELECT * FROM Maxes;
SELECT * FROM BodyWeight;
SELECT * FROM TrainingLog;
