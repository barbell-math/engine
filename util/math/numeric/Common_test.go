package numeric;

import (
    "testing"
    "github.com/barbell-math/engine/util/test"
    customerr "github.com/barbell-math/engine/util/err"
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

func TestRangeNoValues(t *testing.T){
    sequence,err:=Range(0,0,0).Collect();
    test.BasicTest(0,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    sequence,err=Range(0,0,1).Collect();
    test.BasicTest(0,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    sequence,err=Range(5,0,1).Collect();
    test.BasicTest(0,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
}

func TestRangeSomeValues(t *testing.T){
    sequence,err:=Range(0,5,1).Collect();
    test.BasicTest(5,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    for i:=0; i<5; i++ {
        test.BasicTest(i,sequence[i],"Range produced incorrect values.",t);
    }
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    sequence,err=Range(0,5,2).Collect();
    test.BasicTest(3,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    for i:=0; i<3; i++ {
        test.BasicTest(i*2,sequence[i],"Range produced incorrect values.",t);
    }
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    sequence,err=Range(0,5,4).Collect();
    test.BasicTest(2,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    for i:=0; i<2; i++ {
        test.BasicTest(i*4,sequence[i],"Range produced incorrect values.",t);
    }
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    sequence,err=Range(0,5,5).Collect();
    test.BasicTest(1,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    test.BasicTest(0,sequence[0],"Range produced incorrect values.",t);
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
}

func TestRangeStartValue(t *testing.T){
    sequence,err:=Range(5,10,1).Collect();
    test.BasicTest(5,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    for i:=0; i<5; i++ {
        test.BasicTest(i+5,sequence[i],"Range produced incorrect values.",t);
    }
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    sequence,err=Range(5,10,2).Collect();
    test.BasicTest(3,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    for i:=0; i<3; i++ {
        test.BasicTest(i*2+5,sequence[i],"Range produced incorrect values.",t);
    }
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    sequence,err=Range(5,10,5).Collect();
    test.BasicTest(1,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    test.BasicTest(5,sequence[0],"Range produced incorrect values.",t);
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
}

func TestRangeNegativeStep(t *testing.T){
    sequence,err:=Range(5,0,-1).Collect();
    test.BasicTest(5,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    for i:=0; i<5; i++ {
        test.BasicTest(5-i,sequence[i],"Range produced incorrect values.",t);
    }
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    sequence,err=Range(5,0,-2).Collect();
    test.BasicTest(3,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    for i:=0; i<3; i++ {
        test.BasicTest(5-i*2,sequence[i],"Range produced incorrect values.",t);
    }
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    sequence,err=Range(5,0,-4).Collect();
    test.BasicTest(2,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    for i:=0; i<2; i++ {
        test.BasicTest(5-i*4,sequence[i],"Range produced incorrect values.",t);
    }
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
    sequence,err=Range(5,0,-5).Collect();
    test.BasicTest(1,len(sequence),
        "Range when start=stop produced more values than it should have.",t,
    );
    test.BasicTest(5,sequence[0],"Range produced incorrect values.",t);
    test.BasicTest(nil,err,"Error was raised when it shouldn't have been.",t);
}
