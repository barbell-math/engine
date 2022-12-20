package db;

import (
    "testing"
    "database/sql"
    "github.com/barbell-math/block/util/test"
    customReflect "github.com/barbell-math/block/util/reflect"
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
    createExerciseTestData();
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
    createExerciseTestData();
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

func TestCustomReadQueryNonStructVal(t *testing.T){
    setup();
    createExerciseTestData();
    cntr:=0;
    err:=CustomReadQuery(&testDB,"SELECT Id FROM Exercise ORDER BY Id;",
        []any{},func(v *int){
            cntr++;
    });
    test.BasicTest(0, cntr,
        "Custom read query read values it was not supposed to.",t,
    );
    if !customReflect.IsNonStructValue(err) {
        test.FormatError(customReflect.NonStructValue(""),err,
            "Custom read query did not raise the correct error.",t,
        );
    }
}

func TestCustomDeleteQueryWrongQueryType(t *testing.T){
    setup();
    cntr,err:=CustomDeleteQuery(
        &testDB,"UPDATE Exercise SET Name WHERE Id=$1;",[]any{0},
    );
    test.BasicTest(int64(0), cntr,
        "Custom delete query deleted values it was not supposed to.",t,
    );
    if !IsUnsupportedQueryType(err) {
        test.BasicTest(UnsupportedQueryType(""),err,
            "Custom delete query did not return error on non-delete stmt.",t,
        );
    }
}

func TestCustomDeleteQuery(t *testing.T){
    setup();
    createExerciseTestData();
    cntr,err:=CustomDeleteQuery(&testDB,"DELETE FROM Exercise;",[]any{});
    test.BasicTest(int64(3),cntr,
        "Custom delete query deleted values it was not supposed to.",t,
    );
    test.BasicTest(nil,err,
        "Custom delete query created an error it was not supposed to.",t,
    );
}

func TestCustomDeleteQueryTypes(t *testing.T){
    setup();
    createExerciseTestData();
    cntr,err:=CustomDeleteQuery(&testDB,
        `DELETE FROM Exercise
         WHERE Id IN (
             SELECT Exercise.Id
             FROM Exercise
             JOIN ExerciseType
             ON ExerciseType.Id=Exercise.TypeID
             WHERE ExerciseType.T='Compound Accessory'
         );`,[]any{});
    test.BasicTest(int64(0),cntr,
        "Custom delete query deleted values it was not supposed to.",t,
    );
    test.BasicTest(nil,err,
        "Custom delete query created an error it was not supposed to.",t,
    );
    cntr,err=CustomDeleteQuery(&testDB,
        `DELETE FROM Exercise
         WHERE Id IN (
             SELECT Exercise.Id
             FROM Exercise
             JOIN ExerciseType
             ON ExerciseType.Id=Exercise.TypeID
             WHERE ExerciseType.T='Accessory'
         );`,[]any{});
    test.BasicTest(int64(3),cntr,
        "Custom delete query deleted values it was not supposed to.",t,
    );
    test.BasicTest(nil,err,
        "Custom delete query created an error it was not supposed to.",t,
    );
}
