package db;

import (
    "time"
    "testing"
    "database/sql"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
    "github.com/carmichaeljr/powerlifting-engine/settings"
)

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

func TestGetClientID(t *testing.T){
    setup();
    Create(&testDB,
        Client{FirstName: "testF", LastName: "testL", Email: "test@test.com"},
        Client{FirstName: "testF1", LastName: "testL1", Email: "test1@test.com"},
        Client{FirstName: "testF2", LastName: "testL2", Email: "test2@test.com"},
        Client{FirstName: "testF2", LastName: "testL3", Email: "test3@test.com"},
    );
    id,err:=GetClientID(&testDB,"test@test.com");
    testUtil.BasicTest(nil,err,
        "Client was not found when it should have been.",t,
    );
    if id<=0 {
        testUtil.FormatError(">0",id,"ID was not set appropriately.",t);
    }
    id,err=GetClientID(&testDB,"test1@test.com");
    testUtil.BasicTest(nil,err,
        "Client was not found when it should have been.",t,
    );
    if id<=0 {
        testUtil.FormatError(">0",id,"ID was not set appropriately.",t);
    }
    id,err=GetClientID(&testDB,"testing@test.com");
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
