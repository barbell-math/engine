package util;

import (
    "testing"
)

func TestCSVInterface(t *testing.T){
    setup();
    err:=CSVToDBTable("../testData/ExerciseTypeTestData.csv",',',"",
    func(e *ExerciseType){
        fmt.Println(*e);
    });
    fmt.Println(err);
}
