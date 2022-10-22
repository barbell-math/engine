package db;

import (
    "testing"
)

func TestCSVInterface(t *testing.T){
    setup();
    err:=CSVToDBTable("../testData/ExerciseTypeTestData.csv",',',"",
    func(e *ExerciseType){
        fmt.Println(*e);
    });
    testUtil.BasicTest(
    fmt.Println(err);
}
