package db;

import (
    "fmt"
    "time"
    "database/sql"
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
)

func createTestHelper[R DBTable](
        row1 R,
        row2 R,
        row3 R,
        row4 R,
        row5 R,
        t *testing.T){
    var cnt int=0;
    var id1, id2, id3 []int;
    _,err:=Create[R](&testDB);
    testUtil.BasicTest(sql.ErrNoRows,err,
        "Not creating any rows did not result in appropriate error.",t,
    );
    id1,err=Create(&testDB,row1);
    testUtil.BasicTest(nil,err,"Could not create value in database.",t);
    testUtil.BasicTest(1 ,id1[0],"Value was not created correctly.",t);
    testUtil.BasicTest(
        1,len(id1),"More values were created than should have been.",t,
    );
    id2,err=Create(&testDB,row2);
    testUtil.BasicTest(nil,err,"Could not create value in database.",t);
    testUtil.BasicTest(2 ,id2[0],"Value was not created correctly.",t);
    testUtil.BasicTest(
        1,len(id2),"More values were created than should have been.",t,
    );
    err=testDB.db.QueryRow(
        fmt.Sprintf("SELECT COUNT(*) FROM %s;",getTableName(&row1)),
    ).Scan(&cnt);
    testUtil.BasicTest(nil,err,"Could not access table for counting.",t);
    testUtil.BasicTest(2,cnt,"Wrong number of rows were in table.",t);
    id3,err=Create(&testDB,row3,row4,row5);
    testUtil.BasicTest(nil,err,"Could not create value in database.",t);
    testUtil.BasicTest(3,id3[0],"Value was not created correctly.",t);
    testUtil.BasicTest(4,id3[1],"Value was not created correctly.",t);
    testUtil.BasicTest(5,id3[2],"Value was not created correctly.",t);
    err=testDB.db.QueryRow(
        fmt.Sprintf("SELECT COUNT(*) FROM %s;",getTableName(&row1)),
    ).Scan(&cnt);
    testUtil.BasicTest(nil,err,"Could not access table for counting.",t);
    testUtil.BasicTest(5,cnt,"Wrong number of rows were in table.",t);
}

func TestCreate(t *testing.T){
    setup();
    createTestHelper(
        ExerciseType{Id: -1, T: "TestType", Description: "TestTypeDescription"},
        ExerciseType{Id: -1, T: "TestType1", Description: "TestTypeDescription1"},
        ExerciseType{Id: -1, T: "TestType2", Description: "TestTypeDescription1"},
        ExerciseType{Id: -1, T: "TestType3", Description: "TestTypeDescription1"},
        ExerciseType{Id: -1, T: "TestType4", Description: "TestTypeDescription1"},
        t,
    );
    createTestHelper(
        ExerciseFocus{Focus: "TestFocus"},
        ExerciseFocus{Focus: "TestFocus1"},
        ExerciseFocus{Focus: "TestFocus2"},
        ExerciseFocus{Focus: "TestFocus3"},
        ExerciseFocus{Focus: "TestFocus4"},
        t,
    );
    createTestHelper(
        Exercise{Name: "test", TypeID: 1, FocusID: 1},
        Exercise{Name: "test1", TypeID: 1, FocusID: 1},
        Exercise{Name: "test2", TypeID: 1, FocusID: 1},
        Exercise{Name: "test3", TypeID: 1, FocusID: 1},
        Exercise{Name: "test4", TypeID: 1, FocusID: 1},
        t,
    );
    createTestHelper(
        Client{FirstName: "test", LastName: "test", Email: "test@test.com"},
        Client{FirstName: "test1", LastName: "test1", Email: "test1@test.com"},
        Client{FirstName: "test1", LastName: "test1", Email: "test2@test.com"},
        Client{FirstName: "test1", LastName: "test1", Email: "test3@test.com"},
        Client{FirstName: "test1", LastName: "test1", Email: "test4@test.com"},
        t,
    );
    createTestHelper(
        BodyWeight{ClientID: 1, Weight: 1.00, Date: time.Now()},
        BodyWeight{ClientID: 1, Weight: 2.00, Date: time.Now()},
        BodyWeight{ClientID: 1, Weight: 2.00, Date: time.Now()},
        BodyWeight{ClientID: 1, Weight: 2.00, Date: time.Now()},
        BodyWeight{ClientID: 1, Weight: 2.00, Date: time.Now()},
        t,
    );
    createTestHelper(
        Rotation{ClientID: 1, StartDate: time.Now(), EndDate: time.Now()},
        Rotation{ClientID: 2, StartDate: time.Now(), EndDate: time.Now()},
        Rotation{ClientID: 2, StartDate: time.Now(), EndDate: time.Now()},
        Rotation{ClientID: 2, StartDate: time.Now(), EndDate: time.Now()},
        Rotation{ClientID: 2, StartDate: time.Now(), EndDate: time.Now()},
        t,
    );
    createTestHelper(
        TrainingLog{
            ClientID: 1, ExerciseID: 1, DatePerformed: time.Now(),
            Weight: 1.00, Sets: 1.00, Reps: 1, Intensity: 0.50, RotationID: 1,
        },
        TrainingLog{
            ClientID: 1, ExerciseID: 1, DatePerformed: time.Now(),
            Weight: 2.00, Sets: 2.00, Reps: 2, Intensity: 0.60, RotationID: 1,
        },
        TrainingLog{
            ClientID: 1, ExerciseID: 1, DatePerformed: time.Now(),
            Weight: 1.00, Sets: 1.00, Reps: 1, Intensity: 0.50, RotationID: 1,
        },
        TrainingLog{
            ClientID: 1, ExerciseID: 1, DatePerformed: time.Now(),
            Weight: 1.00, Sets: 1.00, Reps: 1, Intensity: 0.50, RotationID: 1,
        },
        TrainingLog{
            ClientID: 1, ExerciseID: 1, DatePerformed: time.Now(),
            Weight: 1.00, Sets: 1.00, Reps: 1, Intensity: 0.50, RotationID: 1,
        },t,
    );
}

func TestRead(t *testing.T){
    setup();
    vals:=[]ExerciseType{
        ExerciseType{T: "TestType", Description: "TestTypeDescription"},
        ExerciseType{T: "TestType1", Description: "TestTypeDescription1"},
        ExerciseType{T: "TestType2", Description: "TestTypeDescription1"},
    };
    readFilter:=func (col string) bool { return col=="Description"; };
    for _,val:=range(vals) {
        Create(&testDB,val);
    }
    cntr:=0;
    err:=Read(&testDB,vals[0],func(col string) bool {
            return col=="NonExistantCol";
        }, func(exercise *ExerciseType){ cntr++; },
    );
    if !util.IsFilterRemovedAllColumns(err) {
        testUtil.FormatError(
            util.FilterRemovedAllColumns(""),err,
            "Filtering all columns did not result in the appropriate error.",t,
        );
    }
    testUtil.BasicTest(0, cntr,"Read selected values it was not supposed to.",t);
    cntr=0;
    err=Read(&testDB,vals[0],readFilter,func(e *ExerciseType){
        cntr++;
        testUtil.BasicTest(
            "TestType",e.T,"Exercise type selected was not correct.",t,
        );
        testUtil.BasicTest(
            "TestTypeDescription",e.Description,"Exercise type selected was not correct.",t,
        );
    });
    testUtil.BasicTest(nil,err,"Read returned an error it was not supposed to.",t);
    testUtil.BasicTest(1,cntr,"Read selected values it was not supposed to.",t);
    cntr=0;
    err=Read(&testDB,vals[1],readFilter,func(e *ExerciseType){
        cntr++;
        testUtil.BasicTest(
            "TestTypeDescription1",e.Description,"Exercise type selected was not correct.",t,
        );
    });
    testUtil.BasicTest(nil,err,"Read returned an error it was not supposed to.",t);
    testUtil.BasicTest(2,cntr,"Read selected values it was not supposed to.",t);
}

func TestUpdate(t *testing.T){
    setup();
    numRows,err:=Update(
        &testDB,ExerciseType{},NoFilter,ExerciseType{},AllButIDFilter,
    );
    testUtil.BasicTest(int64(0), numRows,"Update created rows.",t);
    testUtil.BasicTest(nil,err,"Updating 0 rows resulted in an error.",t);
    Create(&testDB,ExerciseType{T: "test", Description: "testing"});
    Create(&testDB,ExerciseType{T: "test1", Description: "testing"});
    Create(&testDB,ExerciseType{T: "test2", Description: "testing"});
    numRows,err=Update(
        &testDB,
        ExerciseType{},
        GenColFilter(false),
        ExerciseType{},
        AllButIDFilter,
    );
    testUtil.BasicTest(
        int64(0), numRows,"Update updated rows it wasn't supposed to.",t,
    );
    if !util.IsFilterRemovedAllColumns(err) {
        testUtil.FormatError(
            util.FilterRemovedAllColumns(""),err,
            "Filtering all columns did not result in the appropriate error.",t,
        );
    }
    numRows,err=Update(
        &testDB,
        ExerciseType{},
        AllButIDFilter,
        ExerciseType{},
        GenColFilter(false),
    );
    testUtil.BasicTest(
        int64(0), numRows,"Update updated rows it wasn't supposed to.",t,
    );
    if !util.IsFilterRemovedAllColumns(err) {
        testUtil.FormatError(
            util.FilterRemovedAllColumns(""),err,
            "Filtering all columns did not result in the appropriate error.",t,
        );
    }
    numRows,err=Update(
        &testDB,
        ExerciseType{T: "test"},
        GenColFilter(false,"T"),
        ExerciseType{T: "updatedTest", Description: "updatedTesting"},
        AllButIDFilter,
    );
    testUtil.BasicTest(
        int64(1),numRows,"Update did not update the correct number of rows.",t,
    );
    testUtil.BasicTest(nil,err,"Updating rows resulted in an error.",t);
    numRows,err=Update(
        &testDB,
        ExerciseType{T: "test1", Description: "testing"},
        GenColFilter(false,"Description"),
        ExerciseType{Description: "updatedDescription"},
        GenColFilter(false,"Description"),
    );
    testUtil.BasicTest(
        int64(2),numRows,"Update did not update the correct number of rows.",t,
    );
    testUtil.BasicTest(nil,err,"Updating rows resulted in an error.",t);
}

func TestDelete(t *testing.T){
    setup();
    Create(&testDB,
        ExerciseType{T: "Test",Description: "testing"},
        ExerciseType{T: "Test1",Description: "testing1"},
        ExerciseType{T: "Test2",Description: "testing1"},
        ExerciseType{T: "Test3",Description: "testing1"},
    );
    res,err:=Delete(&testDB,ExerciseType{},GenColFilter(false));
    if !util.IsFilterRemovedAllColumns(err) {
        testUtil.FormatError(
            util.FilterRemovedAllColumns(""),err,
            "Filtering all columns did not result in the appropriate error.",t,
        );
    }
    res,err=Delete(
        &testDB,
        ExerciseType{T: "Test1",Description: "testing1"},
        AllButIDFilter,
    );
    testUtil.BasicTest(nil,err,"Delete was unsuccessful.",t);
    testUtil.BasicTest(int64(1),res,"Delete removed to many rows.",t);
    res,err=Delete(&testDB,
        ExerciseType{Description: "testing1"},GenColFilter(false,"Description"),
    );
    testUtil.BasicTest(nil,err,"Delete was unsuccessful.",t);
    testUtil.BasicTest(int64(2),res,"Delete removed to many rows.",t);
    res,err=Delete(&testDB,ExerciseType{T: "Test"},GenColFilter(false,"T"));
    testUtil.BasicTest(nil,err,"Delete was unsuccessful.",t);
    testUtil.BasicTest(int64(1),res,"Delete removed to many rows.",t);
    err=testDB.db.QueryRow(
        fmt.Sprintf("SELECT COUNT(*) FROM ExerciseType;"),
    ).Scan(&res);
    testUtil.BasicTest(nil,err,"Could not access table for counting.",t);
    testUtil.BasicTest(int64(0) ,res,"Wrong number of rows were in table.",t);
}

func TestReadAll(t *testing.T){
    setup();
    var cntr int=0;
    err:=ReadAll(&testDB,func(e *ExerciseType){ cntr++ });
    testUtil.BasicTest(nil,err,"ReadAll operations was unsuccessful.",t);
    testUtil.BasicTest(0 ,cntr,"ReadAll did not select all rows.",t);
    for i:=0; i<10; i++ {
        Create(&testDB,
            ExerciseType{T: fmt.Sprintf("test%d",i),Description: "testing"},
        );
    }
    cntr=0;
    err=ReadAll(&testDB,func(e *ExerciseType){ cntr++ });
    testUtil.BasicTest(nil,err,"ReadAll operations was unsuccessful.",t);
    testUtil.BasicTest(10,cntr,"ReadAll did not select all rows.",t);
}

func TestUpdateAll(t *testing.T){
    setup();
    res,err:=UpdateAll(&testDB,ExerciseType{},GenColFilter(false));
    if !util.IsFilterRemovedAllColumns(err) {
        testUtil.FormatError(
            util.FilterRemovedAllColumns(""),err,
            "Filtering all columns did not result in the appropriate error.",t,
        );
    }
    testUtil.BasicTest(int64(0),res,"Update updated rows it was not supposed to.",t);
    res,err=UpdateAll(&testDB,ExerciseType{},GenColFilter(false,"Description"));
    testUtil.BasicTest(nil,err,"UpdateAll operation was unsuccessful.",t);
    testUtil.BasicTest(int64(0),res,"UpdateAll did not update all rows.",t);
    Create(&testDB,ExerciseType{T:"test",Description:"testingDiff"});
    for i:=0; i<10; i++ {
        Create(&testDB,
            ExerciseType{T: fmt.Sprintf("testing%d",i),Description: "testing"},
        );
    }
    res,err=UpdateAll(&testDB,
        ExerciseType{Description: "newDesc"},
        GenColFilter(false,"Description"),
    );
    testUtil.BasicTest(nil,err,"UpdateAll operation was unsuccessful.",t);
    testUtil.BasicTest(int64(11),res,"UpdateAll did not update all rows.",t);
    ReadAll(&testDB,func(e *ExerciseType){
        testUtil.BasicTest("newDesc",e.Description,
            "Description value was not updated properly",t,
        );
    });
}

func TestDeleteAll(t *testing.T){
    setup();
    cntr,err:=DeleteAll[ExerciseType](&testDB);
    testUtil.BasicTest(nil,err,"DeleteAll operations was unsuccessful.",t);
    testUtil.BasicTest(int64(0) ,cntr,"DeleteAll did not delete all rows.",t);
    for i:=0; i<10; i++ {
        Create(&testDB,
            ExerciseType{T: fmt.Sprintf("test%d",i),Description: "testing"},
        );
    }
    cntr,err=DeleteAll[ExerciseType](&testDB);
    testUtil.BasicTest(nil,err,"DeleteAll operations was unsuccessful.",t);
    testUtil.BasicTest(int64(10),cntr,"DeleteAll did not delete all rows.",t);
    err=testDB.db.QueryRow(
        fmt.Sprintf("SELECT COUNT(*) FROM ExerciseType;"),
    ).Scan(&cntr);
    testUtil.BasicTest(nil,err,"Could not access table for counting.",t);
    testUtil.BasicTest(int64(0) ,cntr,"Table was not empty.",t);
}
