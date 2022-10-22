INSERT INTO Client(Id,FirstName,LastName,Email) VALUES (1,'testF','testL','test@test.com');

\copy ExerciseType(Id,T,Description) FROM '../testData/ExerciseTypeTestData.csv' DELIMITER ',' CSV HEADER;

\copy ExerciseFocus(Id,Focus) FROM '../testData/ExerciseFocusTestData.csv' DELIMITER ',' CSV HEADER;

\copy Exercise(Id,Name,TypeID,FocusID) FROM '../testData/ExerciseTestData.csv' DELIMITER ',' CSV HEADER;

\copy Rotation(Id,ClientID,StartDate,EndDate) FROM '../testData/RotationTestData.csv' DELIMITER ',' CSV HEADER;

\copy TrainingLog(Id,ClientID,ExerciseID,RotationID,DatePerformed,Weight,Sets,Reps,Intensity,Effort,Volume) FROM '../testData/TrainingLogTestData.csv' DELIMITER ',' CSV HEADER;
