package db;

import (
    "testing"
    "database/sql"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
)

func TestCustomReadQueryWrongQueryType(t *testing.T){
    setup();
    cntr:=0;
    err:=CustomReadQuery(&testDB,"UPDATE Exercise SET Name WHERE Id=$1;",[]any{0},
    func(e *Exercise){
        cntr++;
    });
    testUtil.BasicTest(0, cntr,
        "Custom read query read values it was not supposed to.",t,
    );
    if !util.IsUnsupportedQueryType(err) {
        testUtil.BasicTest(util.UnsupportedQueryType(""),err,
            "Custom read query did not return error on non-select stmt.",t,
        );
    }
}

func TestCustomReadQuery(t *testing.T){
    setup();
    cntr:=0;
    testOrder:=[]string{"Deadlift","Bench","Squat"};
    Create(&testDB,ExerciseType{T: "TestType"});
    Create(&testDB,ExerciseFocus{Focus: "TestFocus"});
    Create(&testDB,Exercise{Name: "Squat", TypeID: 1, FocusID: 1});
    Create(&testDB,Exercise{Name: "Bench", TypeID: 1, FocusID: 1});
    Create(&testDB,Exercise{Name: "Deadlift", TypeID: 1, FocusID: 1});
    err:=CustomReadQuery(&testDB,"SELECT * FROM Exercise ORDER BY Id DESC;",
        []any{},func(e *Exercise){
            testUtil.BasicTest(testOrder[cntr],e.Name,
                "Custom query was not run correctly.",t,
            );
            cntr++;
    });
    testUtil.BasicTest(nil,err,
        "An error was raised when it shouldn't have been.",t,
    );
    testUtil.BasicTest(3, cntr,
        "Custom read query read values it was not supposed to.",t,
    );
}

func TestCustomReadQueryEmpty(t *testing.T){
    setup();
    cntr:=0;
    err:=CustomReadQuery(&testDB,"SELECT * FROM Exercise ORDER BY Id DESC;",
        []any{},func(e *Exercise){
            cntr++;
    });
    testUtil.BasicTest(sql.ErrNoRows,err,
        "Custom read query returned incorrect error.",t,
    );
    testUtil.BasicTest(0, cntr,
        "Custom read query read values it was not supposed to.",t,
    );
}
