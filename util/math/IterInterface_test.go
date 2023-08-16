package math;

import (
    "testing"
	"github.com/barbell-math/block/util/test"
	"github.com/barbell-math/block/util/algo/iter"
)

func TestIterInterfaceAdd(t *testing.T){
    res,err:=iter.SliceElems([]int{1,2,3,4}).Reduce(0,ReduceAdd[int]);
    test.BasicTest(10,res,"Add did not produce correct value.",t);
    test.BasicTest(nil,err,"Add did not produce correct error.",t);
}

func TestIterInterfaceSub(t *testing.T){
    res,err:=iter.SliceElems([]int{1,2,3,4}).Reduce(0,ReduceSub[int]);
    test.BasicTest(-10,res,"Sub did not produce correct value.",t);
    test.BasicTest(nil,err,"Sub did not produce correct error.",t);
}

func TestIterInterfaceMul(t *testing.T){
    res,err:=iter.SliceElems([]int{1,2,3,4}).Reduce(1,ReduceMul[int]);
    test.BasicTest(24,res,"Mul did not produce correct value.",t);
    test.BasicTest(nil,err,"Mul did not produce correct error.",t);
}

func TestIterInterfaceDiv(t *testing.T){
    res,err:=iter.SliceElems([]float32{1.0,2.0}).Reduce(1,ReduceDiv[float32]);
    test.BasicTest(float32(0.5),res,"Div did not produce correct value.",t);
    test.BasicTest(nil,err,"Div did not produce correct error.",t);
}

func TestIterInterfaceDivByZero(t *testing.T){
    _,err:=iter.SliceElems([]float32{1.0,2.0,0.0}).Reduce(1,ReduceDiv[float32]);
    if !IsDivByZero(err) {
        test.FormatError(DivByZero(""),err,
            "Div returned incorrect error when dividing by zero.",t,
        );
    }
}
