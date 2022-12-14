package db;

import (
    "testing"
    "database/sql"
    "github.com/barbell-math/block/util/test"
)

func TestCustomReadQueryWrongQueryType(t *testing.T){
    setup();
    cntr:=0;
    err:=CustomReadQuery(&testDB,"UPDATE Exercise SET Name WHERE Id=$1;",[]any{0},
    func(e *Exercise){
        cntr++;
    });
    test.BasicTest(0, cntr,
        "Custom read query read values it was not supposed to.",t,
    );
    if !IsUnsupportedQueryType(err) {
        test.BasicTest(UnsupportedQueryType(""),err,
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
            test.BasicTest(testOrder[cntr],e.Name,
                "Custom query was not run correctly.",t,
            );
            cntr++;
    });
    test.BasicTest(nil,err,
        "An error was raised when it shouldn't have been.",t,
    );
    test.BasicTest(3, cntr,
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
    test.BasicTest(sql.ErrNoRows,err,
        "Custom read query returned incorrect error.",t,
    );
    test.BasicTest(0, cntr,
        "Custom read query read values it was not supposed to.",t,
    );
}

func TestCustomReadQueryTypes(t *testing.T){
    setup();
    cntr:=0;
    currentId:=1;
    Create(&testDB,ExerciseType{T: "TestType"});
    Create(&testDB,ExerciseFocus{Focus: "TestFocus"});
    Create(&testDB,Exercise{Name: "Squat", TypeID: 1, FocusID: 1});
    err:=CustomReadQuery(&testDB,"SELECT * FROM Exercise WHERE Id=$1;",
        []any{"1"},func(e *Exercise){
            test.BasicTest(currentId,e.Id,"Select on Id was not found properly.",t);
            cntr++;
    });
    test.BasicTest(nil,err,
        "An error was raised when it shouldn't have been.",t,
    );
    test.BasicTest(1, cntr,
        "Custom read query read values it was not supposed to.",t,
    );
}
