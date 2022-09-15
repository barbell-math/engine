package db;

import (
    "fmt"
    "time"
    "database/sql"
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
)

var testDB CRUD;

func TestMain(m *testing.M){
    //setup();
    m.Run();
    teardown();
}

func setup(){
    var err error=nil;
    testDB,err=NewCRUD("localhost",5432,"carmichaeljr","test");
    if err!=nil && err!=util.DataVersionNotAvailable {
        panic("Could not open database for testing.");
    }
    err=testDB.execSQLScript("../sql/globalInit.sql");
    if util.IsSqlScriptNotFound(err) {
        panic("Could not find 'globalInit.sql' file for testing.");
    } else if err!=nil {
        panic("An error occurred running the 'globalInit.sql' script.");
    }
}

func teardown(){
    //testDB.execSQLScript("../sql/globalInit.sql");
    testDB.Close();
}

func TestVersion(t *testing.T){
    setup();
    val,err:=testDB.getDataVersion();
    testUtil.BasicTest(
        sql.ErrNoRows,err,
        "Attempting to get version before it was added was successful.",t,
    );
    err=testDB.setDataVersion(-1);
    testUtil.BasicTest(nil,err,"Could not set data version.",t);
    val,err=testDB.getDataVersion();
    testUtil.BasicTest(nil,err,"Could not access version.",t);
    testUtil.BasicTest(-1,val,"Version was not set correctly.",t);
    err=testDB.setDataVersion(-2);
    testUtil.BasicTest(nil,err,"Could not set data version.",t);
    val,err=testDB.getDataVersion();
    testUtil.BasicTest(nil,err,"Could not access version.",t);
    testUtil.BasicTest(-2,val,"Version was not set correctly.",t);
    err=testDB.db.QueryRow("SELECT COUNT(*) FROM Version;").Scan(&val);
    testUtil.BasicTest(nil,err,"Could not access version for counting.",t);
    testUtil.BasicTest(1,val,"Version table was not limited to only one entry.",t);
}

func createTestHelper[R DBTable](row1 R, row2 R, t *testing.T){
    var id1, id2, cnt int=0, 0, 0;
    //;_,err:=testDB.ReadExerciseType(id1);
    //;testUtil.BasicTest(
    //;    sql.ErrNoRows,err,
    //;    "Reading exercise type before any were added was successful.",t,
    //;);
    id1,err:=Create(&testDB,row1);
    //row1.Id=id1;
    testUtil.BasicTest(nil,err,"Could not create exercise type.",t);
    testUtil.BasicTest(1 ,id1,"Exercise was not created correctly.",t);
    id2,err=Create(&testDB,row2);
    //row2.Id=id2;
    testUtil.BasicTest(nil,err,"Could not create exercise type.",t);
    testUtil.BasicTest(2 ,id2,"Exercise was not created correctly.",t);
    //val,err:=testDB.ReadExerciseType(id1);
    //testUtil.BasicTest(
    //    nil,err,
    //    "Could not read exercise type that was previously inserted.",t,
    //);
    //testUtil.BasicTest(val,row1,"Did not return the correct value.",t);
    //testUtil.BasicTest(val.T,row1.T,"Type in exercise types doe not match",t);
    //val,err=testDB.ReadExerciseType(id2);
    //testUtil.BasicTest(
    //    nil,err,
    //    "Could not read exercise type that was previously inserted.",t,
    //);
    //testUtil.BasicTest(val,row2,"Did not return the correct value.",t);
    err=testDB.db.QueryRow(
        fmt.Sprintf("SELECT COUNT(*) FROM %s;",getTableName(&row1)),
    ).Scan(&cnt);
    testUtil.BasicTest(nil,err,"Could not access exercise types for counting.",t);
    testUtil.BasicTest(2,cnt,"Wrong number of rows were in exercise types.",t);
}

func TestCreate(t *testing.T){
    setup();
    createTestHelper(
        ExerciseType{Id: -1, T: "TestType", Description: "TestTypeDescription"},
        ExerciseType{Id: -1, T: "TestType1", Description: "TestTypeDescription1"},
        t,
    );
    createTestHelper(
        ExerciseFocus{Focus: "TestFocus"},
        ExerciseFocus{Focus: "TestFocus1"},
        t,
    );
    createTestHelper(
        Exercise{Name: "test", TypeID: 1, FocusID: 1},
        Exercise{Name: "test1", TypeID: 1, FocusID: 1},
        t,
    );
    createTestHelper(
        Client{FirstName: "test", LastName: "test", Email: "test@test.com"},
        Client{FirstName: "test1", LastName: "test1", Email: "test1@test.com"},
        t,
    );
    createTestHelper(
        BodyWeight{UserID: 1, Weight: 1.00, Date: time.Now()},
        BodyWeight{UserID: 1, Weight: 2.00, Date: time.Now()},
        t,
    );
    createTestHelper(
        Rotation{UserID: 1, StartDate: time.Now(), EndDate: time.Now()},
        Rotation{UserID: 2, StartDate: time.Now(), EndDate: time.Now()},
        t,
    );
    createTestHelper(
        TrainingLog{
            UserID: 1, ExerciseID: 1, DatePerformed: time.Now(),
            Weight: 1.00, Sets: 1.00, Reps: 1, Intensity: 0.50,
        },
        TrainingLog{
            UserID: 1, ExerciseID: 1, DatePerformed: time.Now(),
            Weight: 2.00, Sets: 2.00, Reps: 2, Intensity: 0.60,
        },
        t,
    );
}

func TestRead(t *testing.T){
    setup();
    Create(
        &testDB,
        ExerciseType{Id: 1, T: "TestType", Description: "TestTypeDescription"},
    );
    Create(
        &testDB,
        ExerciseType{Id: 1, T: "TestType1", Description: "TestTypeDescription1"},
    );
    Create(
        &testDB,
        ExerciseType{Id: 1, T: "TestType1", Description: "TestTypeDescription2"},
    );
    cntr:=0;
    err:=Read(&testDB,ExerciseType{T: "TestType1"},func(col string) bool {
        return col=="T";
    },
    func(exercise *ExerciseType){
        cntr++;
        testUtil.BasicTest(
            "TestType1",exercise.T,"Read did not select the correct values.",t,
        );
    });
    testUtil.BasicTest(nil,err,"Could not read values from database.",t);
    testUtil.BasicTest(2,cntr,"Read did not select all values.",t);
    fmt.Println(err);
}
