package db;

import (
    "fmt"
    "time"
    "database/sql"
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
    "github.com/carmichaeljr/powerlifting-engine/settings"
)

var testDB CRUD;

func TestMain(m *testing.M){
    //setup();
    settings.ReadSettings("../testData/dbTestSettings.json");
    m.Run();
    teardown();
}

func setup(){
    var err error=nil;
    testDB,err=NewCRUD(settings.DBHost(),settings.DBPort(),settings.DBName());
    if err!=nil && err!=util.DataVersionNotAvailable {
        panic("Could not open database for testing.");
    }
    if err=testDB.ResetDB(); err!=nil {
        panic(fmt.Sprintf("Could not reset DB for testing. Check location of global init SQL file relative to the ./testData/dbTestSettings.json file. \n  | Given err: %v",err));
    }
}

func teardown(){
    //testDB.ResetDB();
    testDB.Close();
}

func TestVersion(t *testing.T){
    setup();
    //val,err:=testDB.getDataVersion();
    //testUtil.BasicTest(
    //    sql.ErrNoRows,err,
    //    "Attempting to get version before it was added was successful.",t,
    //);
    err:=testDB.setDataVersion(-1);
    testUtil.BasicTest(nil,err,"Could not set data version.",t);
    val,err:=testDB.getDataVersion();
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
    id3,err=Create(&testDB,row1,row1,row2);
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
        BodyWeight{ClientID: 1, Weight: 1.00, Date: time.Now()},
        BodyWeight{ClientID: 1, Weight: 2.00, Date: time.Now()},
        t,
    );
    createTestHelper(
        Rotation{ClientID: 1, StartDate: time.Now(), EndDate: time.Now()},
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
        },t,
    );
}

func readTestHelper[R DBTable](
        vals []R,
        readFilter func(col string) bool,
        readFunc1 func(val *R),
        readFunc2 func(val *R),
        t *testing.T){
    for _,val:=range(vals) {
        Create(&testDB,val);
    }
    cntr:=0;
    err:=Read(&testDB,vals[0],func(col string) bool {
        return col=="NonExistantCol";
    },
    func(exercise *R){ cntr++; });
    if !util.IsFilterRemovedAllColumns(err) {
        testUtil.FormatError(
            util.FilterRemovedAllColumns(""),err,
            "Filtering all columns did not result in the appropriate error.",t,
        );
    }
    testUtil.BasicTest(0, cntr,"Read selected values it was not supposed to.",t);
    cntr=0;
    err=Read(&testDB,vals[0],readFilter,func(val *R){
        cntr++;
        readFunc1(val);
    });
    testUtil.BasicTest(nil,err,"Read returned an error it was not supposed to.",t);
    testUtil.BasicTest(1,cntr,"Read selected values it was not supposed to.",t);
    cntr=0;
    err=Read(&testDB,vals[1],readFilter,func(val *R){
        cntr++;
        readFunc2(val);
    });
    testUtil.BasicTest(nil,err,"Read returned an error it was not supposed to.",t);
    testUtil.BasicTest(2,cntr,"Read selected values it was not supposed to.",t);
}

func TestRead(t *testing.T){
    setup();
    readTestHelper([]ExerciseType{
        ExerciseType{T: "TestType", Description: "TestTypeDescription"},
        ExerciseType{T: "TestType1", Description: "TestTypeDescription1"},
        ExerciseType{T: "TestType1", Description: "TestTypeDescription1"},
    },
    func (col string) bool { return col=="T"; },
    func (e *ExerciseType) {
        testUtil.BasicTest(
            "TestType",e.T,"Exercise type selected was not correct.",t,
        );
        testUtil.BasicTest(
            "TestTypeDescription",e.Description,"Exercise type selected was not correct.",t,
        );
    },
    func (e *ExerciseType) {
        testUtil.BasicTest(
            "TestType1",e.T,"Exercise type selected was not correct.",t,
        );
        testUtil.BasicTest(
            "TestTypeDescription1",e.Description,"Exercise type selected was not correct.",t,
        );
    },t);
    readTestHelper([]ExerciseFocus{
        ExerciseFocus{Focus: "Focus"},
        ExerciseFocus{Focus: "Focus1"},
        ExerciseFocus{Focus: "Focus1"},
    },
    func (col string) bool { return col=="Focus"; },
    func (focus *ExerciseFocus) {
        testUtil.BasicTest(
            "Focus",focus.Focus,"Exercise focus selected was not correct.",t,
        );
    },
    func (focus *ExerciseFocus) {
        testUtil.BasicTest(
            "Focus1",focus.Focus,"Exercise focus selected was not correct.",t,
        );
    },t);
    readTestHelper([]Exercise{
        Exercise{Name: "Exercise", TypeID: 1, FocusID: 1},
        Exercise{Name: "Exercise1", TypeID: 2, FocusID: 2},
        Exercise{Name: "Exercise1", TypeID: 2, FocusID: 2},
    },
    func (col string) bool { return col=="Name"; },
    func (e *Exercise) {
        testUtil.BasicTest(
            "Exercise",e.Name,"Exercise selected was not correct.",t,
        );
        testUtil.BasicTest(1,e.TypeID,"Exercise selected was not correct.",t);
        testUtil.BasicTest(1,e.FocusID,"Exercise selected was not correct.",t);
    },
    func (e *Exercise) {
        testUtil.BasicTest(
            "Exercise1",e.Name,"Exercise focus selected was not correct.",t,
        );
        testUtil.BasicTest(2,e.TypeID,"Exercise selected was not correct.",t);
        testUtil.BasicTest(2,e.FocusID,"Exercise selected was not correct.",t);
    },t);
    readTestHelper([]Client{
        Client{FirstName: "test", LastName: "test", Email: "test@test.com"},
        Client{FirstName: "test1", LastName: "test1", Email: "test1@test.com"},
        Client{FirstName: "test1", LastName: "test1", Email: "test1@test.com"},
    },
    func (col string) bool { return col=="FirstName"; },
    func (c *Client) {
        testUtil.BasicTest(
            "test",c.FirstName,"Client selected was not correct.",t,
        );
        testUtil.BasicTest(
            "test",c.LastName,"Client selected was not correct.",t,
        );
        testUtil.BasicTest(
            "test@test.com",c.Email,"Client selected was not correct.",t,
        );
    },
    func (c *Client) {
        testUtil.BasicTest(
            "test1",c.FirstName,"Client selected was not correct.",t,
        );
        testUtil.BasicTest(
            "test1",c.LastName,"Client selected was not correct.",t,
        );
        testUtil.BasicTest(
            "test1@test.com",c.Email,"Client selected was not correct.",t,
        );
    },t);
    readTestHelper([]Rotation{
        Rotation{ClientID: 1, StartDate: time.Now(), EndDate: time.Now()},
        Rotation{ClientID: 2, StartDate: time.Now(), EndDate: time.Now()},
        Rotation{ClientID: 2, StartDate: time.Now(), EndDate: time.Now()},
    },
    func (col string) bool { return col=="ClientID"; },
    func (r *Rotation) {
        testUtil.BasicTest(
            1,r.ClientID,"Rotation selected was not correct.",t,
        );
        testUtil.BasicTest(
            time.Now().Format("00-00-0000"),r.StartDate.Format("00-00-0000"),
            "Rotation selected was not correct.",t,
        );
        testUtil.BasicTest(
            time.Now().Format("00-00-0000"),r.EndDate.Format("00-00-0000"),
            "Rotation selected was not correct.",t,
        );
    },
    func (r *Rotation) {
        testUtil.BasicTest(
            2,r.ClientID,"Rotation focus selected was not correct.",t,
        );
        testUtil.BasicTest(
            time.Now().Format("00-00-0000"),r.StartDate.Format("00-00-0000"),
            "Rotation selected was not correct.",t,
        );
        testUtil.BasicTest(
            time.Now().Format("00-00-0000"),r.EndDate.Format("00-00-0000"),
            "Rotation selected was not correct.",t,
        );
    },t);
    readTestHelper([]BodyWeight{
        BodyWeight{ClientID: 1, Weight: 1.0, Date: time.Now()},
        BodyWeight{ClientID: 2, Weight: 2.0, Date: time.Now()},
        BodyWeight{ClientID: 2, Weight: 2.0, Date: time.Now()},
    },
    func (col string) bool { return col=="ClientID"; },
    func (b *BodyWeight) {
        testUtil.BasicTest(
            1,b.ClientID,"Bodyweight selected was not correct.",t,
        );
        testUtil.BasicTest(
            float32(1),b.Weight,"Bodyweight selected was not correct.",t,
        );
        testUtil.BasicTest(
            time.Now().Format("00-00-0000"),b.Date.Format("00-00-0000"),
            "BodyWeight selected was not correct.",t,
        );
    },
    func (b *BodyWeight) {
        testUtil.BasicTest(
            2,b.ClientID,"Bodyweight selected was not correct.",t,
        );
        testUtil.BasicTest(
            float32(2),b.Weight,"Bodyweight selected was not correct.",t,
        );
        testUtil.BasicTest(
            time.Now().Format("00-00-0000"),b.Date.Format("00-00-0000"),
            "BodyWeight selected was not correct.",t,
        );
    },t);
    readTestHelper([]TrainingLog{
        TrainingLog{
            ClientID: 1, ExerciseID: 1, DatePerformed: time.Now(),
            Weight: 1.0, Sets: 1.0, Reps: 1, Intensity: 0.5, RotationID: 1,
        },
        TrainingLog{
            ClientID: 2, ExerciseID: 1, DatePerformed: time.Now(),
            Weight: 1.0, Sets: 1.0, Reps: 1, Intensity: 0.5, RotationID: 1,
        },
        TrainingLog{
            ClientID: 2, ExerciseID: 1, DatePerformed: time.Now(),
            Weight: 1.0, Sets: 1.0, Reps: 1, Intensity: 0.5, RotationID: 1,
        },
    },
    func (col string) bool { return col=="ClientID"; },
    func (b *TrainingLog) {
        testUtil.BasicTest(
            1,b.ClientID,"Training log selected was not correct.",t,
        );
        testUtil.BasicTest(
            1,b.ExerciseID,"Training log selected was not correct.",t,
        );
        testUtil.BasicTest(
            time.Now().Format("00-00-0000"),b.DatePerformed.Format("00-00-0000"),
            "Training log selected was not correct.",t,
        );
        testUtil.BasicTest(
            1,b.Reps,"Training log selected was not correct.",t,
        );
    },
    func (b *TrainingLog) {
        testUtil.BasicTest(
            2,b.ClientID,"Training log selected was not correct.",t,
        );
        testUtil.BasicTest(
            1,b.ExerciseID,"Training log selected was not correct.",t,
        );
        testUtil.BasicTest(
            time.Now().Format("00-00-0000"),b.DatePerformed.Format("00-00-0000"),
            "Training log selected was not correct.",t,
        );
        testUtil.BasicTest(
            1,b.Reps,"Training log selected was not correct.",t,
        );
    },t);
    cntr:=0;
    err:=Read(&testDB,TrainingLog{ClientID: 2, Weight: 1.0},
    func (col string) bool {
        return col=="ClientID" || col=="Weight";
    },
    func (l *TrainingLog){
        cntr++;
        testUtil.BasicTest(2,l.ClientID,"Training log selected was not correct.",t);
    });
    testUtil.BasicTest(nil,err,"Read was unsuccessful.",t);
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
        ExerciseType{T: "Test1",Description: "testing"},
        ExerciseType{T: "Test1",Description: "testing"},
        ExerciseType{T: "Test1",Description: "testing1"},
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
    res,err=Delete(&testDB,ExerciseType{T: "Test1"},GenColFilter(false,"T"));
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
        Create(&testDB,ExerciseType{T: "test",Description: "testing"});
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
        Create(&testDB,ExerciseType{T: "testing",Description: "testing"});
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
        Create(&testDB,ExerciseType{T: "test",Description: "testing"});
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

func TestGetExerciseID(t *testing.T){
    setup();
    settings.Modify(func(s *settings.Settings){
        s.DBInfo.DataVersion=1;
    });
    testDB.RunDataConversion();
    id,err:=GetExerciseID(&testDB,"Squat");
    testUtil.BasicTest(nil,err,
        "Exercise was not found when it should have been.",t,
    );
    if id<=0 {
        testUtil.FormatError(">0",id,"ID was not set appropriately.",t);
    }
    id,err=GetExerciseID(&testDB,"Bench");
    testUtil.BasicTest(nil,err,
        "Exercise was not found when it should have been.",t,
    );
    if id<=0 {
        testUtil.FormatError(">0",id,"ID was not set appropriately.",t);
    }
    id,err=GetExerciseID(&testDB,"NotAnExercise");
    if err!=sql.ErrNoRows {
        testUtil.FormatError(nil,err,
            "No error was generated when getting non-existent exercise.",t,
        );
    }
}

func TestInitClient(t *testing.T){
    setup();
    settings.Modify(func(s *settings.Settings){ s.DBInfo.DataVersion=1; });
    testDB.RunDataConversion();
    c:=Client{
        FirstName: "test",
        LastName: "testl",
        Email: "test@test.com",
    };
    err:=InitClient(&testDB,&c,446,286,545);
    testUtil.BasicTest(nil,err,
        "Init Client returned an error when it shouldn't have.",t,
    );
    sId,_:=GetExerciseID(&testDB,"Squat");
    bId,_:=GetExerciseID(&testDB,"Bench");
    dId,_:=GetExerciseID(&testDB,"Deadlift");
    err=Read(&testDB,TrainingLog{Id: 1},OnlyIDFilter, func(v *TrainingLog){
        y,m,d:=time.Now().Date();
        y1,m1,d1:=v.DatePerformed.Date();
        testUtil.BasicTest(y,y1,"Year is not set correctly in training log.",t);
        testUtil.BasicTest(m,m1,"Month is not set correctly in training log.",t);
        testUtil.BasicTest(d-1,d1,"Day is not set correctly in training log.",t);
        testUtil.BasicTest(1,v.ClientID,
            "Client ID was not set correctly in the zero rotation.",t,
        );
        switch (v.ExerciseID){
            case sId:
                testUtil.BasicTest(float32(446),v.Weight,
                    "Squat 1RM not set correctly in zero rotation.",t,
                );
            case bId:
                testUtil.BasicTest(float32(286),v.Weight,
                    "Squat 1RM not set correctly in zero rotation.",t,
                );
            case dId:
                testUtil.BasicTest(float32(545),v.Weight,
                    "Squat 1RM not set correctly in zero rotation.",t,
                );
            default: t.Errorf(
                "Non SBD exercise max was made from user init | ID: %d.",v.Id,
            );
        }
    });
    testUtil.BasicTest(nil,err,"An error occurred reading the training log.",t);
    err=Read(&testDB,Rotation{Id: 1},OnlyIDFilter,func(v *Rotation){
        y,m,d:=time.Now().Date();
        y1,m1,d1:=v.StartDate.Date();
        testUtil.BasicTest(y,y1,"Year is not set correctly in rotation.",t);
        testUtil.BasicTest(m,m1,"Month is not set correctly in rotation.",t);
        testUtil.BasicTest(d-1,d1,"Day is not set correctly in rotation.",t);
        y1,m1,d1=v.EndDate.Date();
        testUtil.BasicTest(y,y1,"Year is not set correctly in rotation.",t);
        testUtil.BasicTest(m,m1,"Month is not set correctly in rotation.",t);
        testUtil.BasicTest(d,d1,"Day is not set correctly in rotation.",t);
    });
    testUtil.BasicTest(nil,err,"An error occurred reading the training log.",t);
    err=Read(&testDB,Client{Id: 1},OnlyIDFilter,func(v *Client){
        testUtil.BasicTest("test",v.FirstName,"Client f-name not set correctly.",t);
        testUtil.BasicTest("testl",v.LastName,"Client l-name not set correctly.",t);
        testUtil.BasicTest("test@test.com",v.Email,
            "Client email not set correctly.",t,
        );
    });
    testUtil.BasicTest(nil,err,"An error occurred reading the training log.",t);
}

func TestRmClient(t *testing.T){
    setup();
    settings.Modify(func(s *settings.Settings){ s.DBInfo.DataVersion=1; });
    testDB.RunDataConversion();
    c:=Client{
        Id: 1,
        FirstName: "test",
        LastName: "testl",
        Email: "test@test.com",
    };
    err:=util.ChainedErrorOps(
        func(r ...any) (any,error) {
            return nil,InitClient(&testDB,&c,455,286,545);
        }, func(r ...any) (any,error) {
            return Create(&testDB,BodyWeight{
                ClientID: 1, Weight: 148.0, Date: time.Now(),
            });
        }, func (r ...any) (any,error) {
            return Create(&testDB,ModelState{
                ClientID: 1, Date: time.Now(), TimeFrame: 100,
                A: 1, B: 1, C: 1, D: 1, Eps: 1, Eps2: 1,
            });
        },
    );
    testUtil.BasicTest(nil,err,"Database was not setup correctly to run test.",t);
    val,err:=RmClient(&testDB,&c);
    testUtil.BasicTest(nil,err,"RmClient created an error when it shouldn't have.",t);
    testUtil.BasicTest(int64(7),val,"RmClient did not delete all client data.",t);
}
