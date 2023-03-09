package math;

import (
    "testing"
	"github.com/barbell-math/block/util/algo"
	"github.com/barbell-math/block/util/test"
)

func TestIterInterfaceAdd(t *testing.T){
    res,err:=algo.SliceElems([]int{1,2,3,4}).Reduce(Add[int],0);
    test.BasicTest(10,res,"Add did not produce correct value.",t);
    test.BasicTest(nil,err,"Add did not produce correct error.",t);
}

func TestIterInterfaceSub(t *testing.T){
    res,err:=algo.SliceElems([]int{1,2,3,4}).Reduce(Sub[int],0);
    test.BasicTest(-10,res,"Sub did not produce correct value.",t);
    test.BasicTest(nil,err,"Sub did not produce correct error.",t);
}

func TestIterInterfaceMul(t *testing.T){
    res,err:=algo.SliceElems([]int{1,2,3,4}).Reduce(Mul[int],1);
    test.BasicTest(24,res,"Mul did not produce correct value.",t);
    test.BasicTest(nil,err,"Mul did not produce correct error.",t);
}

func TestIterInterfaceDiv(t *testing.T){
    res,err:=algo.SliceElems([]float32{1.0,2.0}).Reduce(Div[float32],1);
    test.BasicTest(float32(0.5),res,"Div did not produce correct value.",t);
    test.BasicTest(nil,err,"Div did not produce correct error.",t);
}

func TestIterInterfaceDivByZero(t *testing.T){
    _,err:=algo.SliceElems([]float32{1.0,2.0,0.0}).Reduce(Div[float32],1);
    if !IsDivByZero(err) {
        test.FormatError(DivByZero(""),err,
            "Div returned incorrect error when dividing by zero.",t,
        );
    }
}
