package math;

import (
    "testing"
    "github.com/barbell-math/block/util/test"
    customerr "github.com/barbell-math/block/util/err"
)

func TestSqErrIncorrectDimensions(t *testing.T){
    _,err:=SqErr([]int{1,2,3},[]int{1,2});
    if !customerr.IsDimensionsDoNotAgree(err) {
        test.FormatError(customerr.DimensionsDoNotAgree(""),err,
            "The incorrect error was raised given unequal input lengths.",t,
        );
    }
}

func TestSqErr(t *testing.T){
    res,err:=SqErr([]int{1,1},[]int{3,3});
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    test.BasicTest(4,res[0],"Incorrect SE was returned.",t);
    test.BasicTest(4,res[1],"Incorrect SE was returned.",t);
    res,err=SqErr([]int{1,3},[]int{1,3});
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    test.BasicTest(0,res[0],"Incorrect SE was returned.",t);
    test.BasicTest(0,res[1],"Incorrect SE was returned.",t);
}

func TestSqErrZeroLength(t *testing.T){
    res,err:=SqErr([]int{},[]int{});
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    test.BasicTest(0 ,len(res),"Incorrect MSE was returned.",t);
}

func TestMeanSqErrIncorrectDimensions(t *testing.T){
    _,err:=MeanSqErr([]int{1,2},[]int{1,2,3});
    if !customerr.IsDimensionsDoNotAgree(err) {
        test.FormatError(customerr.DimensionsDoNotAgree(""),err,
            "The incorrect error was raised given unequal input lengths.",t,
        );
    }
}

func TestMeanSqErr(t *testing.T){
    res,err:=MeanSqErr([]int{1,1},[]int{3,3});
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    test.BasicTest(4,res,"Incorrect MSE was returned.",t);
    res,err=MeanSqErr([]int{1,3},[]int{1,3});
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    test.BasicTest(0 ,res,"Incorrect MSE was returned.",t);
}

func TestMeanSqErrZeroLength(t *testing.T){
    res,err:=MeanSqErr([]int{},[]int{});
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    test.BasicTest(0 ,res,"Incorrect MSE was returned.",t);
}
