package db;

import (
    "time"
    "testing"
    "database/sql"
    "github.com/barbell-math/block/util/test"
    customerr "github.com/barbell-math/block/util/err"
)

func TestGetExerciseByName(t *testing.T){
    setup();
    createExerciseTestData();
    e,err:=GetExerciseByName(&testDB,"Squat");
    test.BasicTest(nil,err,
        "Exercise was not found when it should have been.",t,
    );
    if e.Id<=0 {
        test.FormatError(">0",e.Id,"ID was not set appropriately.",t);
    }
    e,err=GetExerciseByName(&testDB,"NotAnExercise");
    if err!=sql.ErrNoRows {
        test.FormatError(nil,err,
            "No error was generated when getting non-existent exercise.",t,
        );
    }
}

func TestGetClientByEmail(t *testing.T){
    setup();
    Create(&testDB,
        Client{FirstName: "testF", LastName: "testL", Email: "test@test.com"},
        Client{FirstName: "testF1", LastName: "testL1", Email: "test1@test.com"},
        Client{FirstName: "testF2", LastName: "testL2", Email: "test2@test.com"},
        Client{FirstName: "testF2", LastName: "testL3", Email: "test3@test.com"},
    );
    c,err:=GetClientByEmail(&testDB,"test@test.com");
    test.BasicTest(nil,err,
        "Client was not found when it should have been.",t,
    );
    if c.Id<=0 {
        test.FormatError(">0",c.Id,"ID was not set appropriately.",t);
    }
    c,err=GetClientByEmail(&testDB,"test1@test.com");
    test.BasicTest(nil,err,
        "Client was not found when it should have been.",t,
    );
    if c.Id<=0 {
        test.FormatError(">0",c.Id,"ID was not set appropriately.",t);
    }
    c,err=GetClientByEmail(&testDB,"testing@test.com");
    if err!=sql.ErrNoRows {
        test.FormatError(nil,err,
            "No error was generated when getting non-existent exercise.",t,
        );
    }
}

func TestGetExerciseTypeByName(t *testing.T){
    setup();
    createExerciseTestData();
    c,err:=GetExerciseTypeByName(&testDB,"Accessory");
    test.BasicTest(nil,err,
        "Accessory was not found when it should have been.",t,
    );
    if c.Id<=0 {
        test.FormatError(">0",c.Id,"ID was not set appropriately.",t);
    }
    c,err=GetExerciseTypeByName(&testDB,"NotAnExerciseType");
    if err!=sql.ErrNoRows {
        test.FormatError(nil,err,
            "No error was generated when getting non-existent exercise.",t,
        );
    }
}

func TestGetExerciseFocusByName(t *testing.T){
    setup();
    createExerciseTestData();
    c,err:=GetExerciseFocusByName(&testDB,"Squat");
    test.BasicTest(nil,err,
        "Accessory was not found when it should have been.",t,
    );
    if c.Id<=0 {
        test.FormatError(">0",c.Id,"ID was not set appropriately.",t);
    }
    c,err=GetExerciseFocusByName(&testDB,"NotAnExerciseFocus");
    if err!=sql.ErrNoRows {
        test.FormatError(nil,err,
            "No error was generated when getting non-existent exercise.",t,
        );
    }
}

func TestGetStateGeneratorByName(t *testing.T){
    setup();
    Create(&testDB,StateGenerator{T: "TestingStateGenerator"});
    gen,err:=GetStateGeneratorByName(&testDB,"TestingStateGenerator");
    test.BasicTest(nil,err,
        "State Generator was not found when it should have been.",t,
    );
    if gen.Id<=0 {
        test.FormatError(">0",gen.Id,"ID was not set appropriately.",t);
    }
    gen,err=GetStateGeneratorByName(&testDB,"NotAStateGenerator");
    if err!=sql.ErrNoRows {
        test.FormatError(nil,err,
            "No error was generated when getting non-existent state generator.",t,
        );
    }
}

func TestInitClient(t *testing.T){
    setup();
    createExerciseTestData();
    c:=Client{
        FirstName: "test",
        LastName: "testl",
        Email: "test@test.com",
    };
    err:=InitClient(&testDB,&c,446,286,545);
    test.BasicTest(nil,err,
        "Init Client returned an error when it shouldn't have.",t,
    );
    sId,_:=GetExerciseByName(&testDB,"Squat");
    bId,_:=GetExerciseByName(&testDB,"Bench");
    dId,_:=GetExerciseByName(&testDB,"Deadlift");
    err=Read(&testDB,TrainingLog{Id: 1},OnlyIDFilter, func(v *TrainingLog) bool {
        y,m,d:=time.Now().Date();
        y1,m1,d1:=v.DatePerformed.Date();
        test.BasicTest(y,y1,"Year is not set correctly in training log.",t);
        test.BasicTest(m,m1,"Month is not set correctly in training log.",t);
        test.BasicTest(d-1,d1,"Day is not set correctly in training log.",t);
        test.BasicTest(1,v.ClientID,
            "Client ID was not set correctly in the zero rotation.",t,
        );
        switch (v.ExerciseID){
            case sId.Id:
                test.BasicTest(float32(446),v.Weight,
                    "Squat 1RM not set correctly in zero rotation.",t,
                );
            case bId.Id:
                test.BasicTest(float32(286),v.Weight,
                    "Squat 1RM not set correctly in zero rotation.",t,
                );
            case dId.Id:
                test.BasicTest(float32(545),v.Weight,
                    "Squat 1RM not set correctly in zero rotation.",t,
                );
            default: t.Errorf(
                "Non SBD exercise max was made from user init | ID: %d.",v.Id,
            );
        }
        return true;
    });
    test.BasicTest(nil,err,"An error occurred reading the training log.",t);
    err=Read(&testDB,Rotation{Id: 1},OnlyIDFilter,func(v *Rotation) bool {
        y,m,d:=time.Now().Date();
        y1,m1,d1:=v.StartDate.Date();
        test.BasicTest(y,y1,"Year is not set correctly in rotation.",t);
        test.BasicTest(m,m1,"Month is not set correctly in rotation.",t);
        test.BasicTest(d-1,d1,"Day is not set correctly in rotation.",t);
        y1,m1,d1=v.EndDate.Date();
        test.BasicTest(y,y1,"Year is not set correctly in rotation.",t);
        test.BasicTest(m,m1,"Month is not set correctly in rotation.",t);
        test.BasicTest(d,d1,"Day is not set correctly in rotation.",t);
        return true;
    });
    test.BasicTest(nil,err,"An error occurred reading the training log.",t);
    err=Read(&testDB,Client{Id: 1},OnlyIDFilter,func(v *Client) bool {
        test.BasicTest("test",v.FirstName,"Client f-name not set correctly.",t);
        test.BasicTest("testl",v.LastName,"Client l-name not set correctly.",t);
        test.BasicTest("test@test.com",v.Email,
            "Client email not set correctly.",t,
        );
        return true;
    });
    test.BasicTest(nil,err,"An error occurred reading the training log.",t);
}

func TestRmClient(t *testing.T){
    setup();
    Create(&testDB,StateGenerator{T: "State Generator"});
    createExerciseTestData();
    c:=Client{
        Id: 1,
        FirstName: "test",
        LastName: "testl",
        Email: "test@test.com",
    };
    err:=customerr.ChainedErrorOps(
        func(r ...any) (any,error) {
            return nil,InitClient(&testDB,&c,455,286,545);
        }, func(r ...any) (any,error) {
            return Create(&testDB,BodyWeight{
                ClientID: 1, Weight: 148.0, Date: time.Now(),
            });
        }, func (r ...any) (any,error) {
            return Create(&testDB,ModelState{
                ClientID: 1, ExerciseID: 1, StateGeneratorID: 1,
                Date: time.Now(), TimeFrame: 100,
                Eps5: 1, Eps6: 1, Eps7: 1, Eps: 1, Eps2: 1, Eps3: 1,
            });
        }, func (r ...any) (any,error) {
            return Create(&testDB,StateGenerator{
                T: "TestPredictor", Description: "TestDescription",
            });
        }, func (r ...any) (any,error) {
            return Create(&testDB,Prediction{
                StateGeneratorID: 1, TrainingLogID: 1, IntensityPred: 0,
            });
        },
    );
    test.BasicTest(nil,err,"Database was not setup correctly to run test.",t);
    val,err:=RmClient(&testDB,&c);
    test.BasicTest(nil,err,"RmClient created an error when it shouldn't have.",t);
    test.BasicTest(int64(8),val,"RmClient did not delete all client data.",t);
}
