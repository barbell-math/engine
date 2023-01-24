package model;

import (
    "testing"
    "github.com/barbell-math/block/util/test"
    "github.com/barbell-math/block/util/dataStruct"
)

func TestNewSlidingWindowStateGen(t *testing.T){
    _,err:=NewSlidingWindowStateGen(
        dataStruct.Pair[int]{1,0},dataStruct.Pair[int]{0, 1},1,
    );
    if !IsInvalidPredictionState(err) {
        test.FormatError(InvalidPredictionState(""),err,
            "The wrong error was raised when creating an invalid prediction generator.",t,
        );
    }
}
