package math;

import (
    "testing"
    "github.com/barbell-math/block/util/test"
)

//write sqerr tests

func TestMeanSqErrIncorrectDimensions(t *testing.T){
    _,err:=MeanSqErr[int]([]int{1,2},[]int{1,2,3});
    if !IsDimensionsDoNotAgree(err) {
        test.FormatError(DimensionsDoNotAgree(""),err,
            "The incorrect error was raised given unequal input lengths.",t,
        );
    }
}

func TestMeanSqErr(t *testing.T){
    res,err:=MeanSqErr[int]([]int{1,1},[]int{3,3});
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    test.BasicTest(4,res,"Incorrect MSE was returned.",t);
    res,err=MeanSqErr[int]([]int{1,3},[]int{1,3});
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    test.BasicTest(0 ,res,"Incorrect MSE was returned.",t);
}

func TestMeanSqErrZeroLength(t *testing.T){
    res,err:=MeanSqErr[int]([]int{},[]int{});
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    test.BasicTest(0 ,res,"Incorrect MSE was returned.",t);
}
