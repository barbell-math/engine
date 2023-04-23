package iter;

import (
    "fmt"
    "testing"
    "github.com/barbell-math/block/util/test"
)

func valElemIterHelper[T any](val T, err error, r int, t *testing.T){
    var tmp T;
    iter:=ValElem(val,err,r);
    for i:=0; i<3; i++ {
        vIter,eIter,contIter:=iter(Continue);
        if i<r {
            test.BasicTest(val,vIter,"ValElem did not return correct value.",t);
            test.BasicTest(err,eIter,"ValElem did not return correct error.",t);
            test.BasicTest(true,contIter,
                "ValElem did not return correct continue status.",t,
            );
        } else {
            test.BasicTest(tmp,vIter,"ValElem did not return correct value.",t);
            test.BasicTest(nil,eIter,"ValElem did not return correct error.",t);
            test.BasicTest(false,contIter,
                "ValElem did not return correct continue status.",t,
            );
        }
    }
}
func TestValElem(t *testing.T){
    valElemIterHelper(1,nil,1,t);
    valElemIterHelper(2,fmt.Errorf("NEW ERROR"),1,t);
    valElemIterHelper(1,nil,2,t);
    valElemIterHelper(2,fmt.Errorf("NEW ERROR"),2,t);
    valElemIterHelper(1,nil,5,t);
    valElemIterHelper(2,fmt.Errorf("NEW ERROR"),5,t);
}

func sliceElemsIterHelper[T any](vals []T, t *testing.T){
    sIter:=SliceElems(vals);
    for i:=0; i<len(vals); i++ {
        sV,sErr,sBool:=sIter(Continue);
        test.BasicTest(vals[i],sV,
            "SliceElems iteration does not match actual values.",t,
        );
        test.BasicTest(nil,sErr,
            "SliceElems iteration produced an error when it shouldn't have.",t,
        );
        test.BasicTest(true,sBool,
            "SliceElems iteration did not stop when it should have.",t,
        );
    }
    var tmp T;
    sV,sErr,sBool:=sIter(Continue);
    test.BasicTest(tmp,sV,
        "SliceElems iteration does not match actual values.",t,
    );
    test.BasicTest(nil,sErr,
        "SliceElems iteration produced an error when it shouldn't have.",t,
    );
    test.BasicTest(false,sBool,
        "SliceElems iteration did not stop when it should have.",t,
    );
}
func TestSliceElems(t *testing.T){
    sliceElemsIterHelper([]string{"one","two","three"},t);
    sliceElemsIterHelper([]int{1,2,3},t);
    sliceElemsIterHelper([]int{1},t);
    sliceElemsIterHelper([]int{},t);
}

func stringElemsIterHelper(vals string, t *testing.T){
    sIter:=StrElems(vals);
    for i:=0; i<len(vals); i++ {
        sV,sErr,sBool:=sIter(Continue);
        test.BasicTest(vals[i],sV,
            "SliceElems iteration does not match actual values.",t,
        );
        test.BasicTest(nil,sErr,
            "SliceElems iteration produced an error when it shouldn't have.",t,
        );
        test.BasicTest(true,sBool,
            "SliceElems iteration did not stop when it should have.",t,
        );
    }
    var tmp string;
    sV,sErr,sBool:=sIter(Continue);
    test.BasicTest(tmp,sV,
        "SliceElems iteration does not match actual values.",t,
    );
    test.BasicTest(nil,sErr,
        "SliceElems iteration produced an error when it shouldn't have.",t,
    );
    test.BasicTest(false,sBool,
        "SliceElems iteration did not stop when it should have.",t,
    );
}
func TestStringElems(t *testing.T){
    sliceElemsIterHelper([]string{"one","two","three"},t);
    sliceElemsIterHelper([]int{1,2,3},t);
    sliceElemsIterHelper([]int{1},t);
    sliceElemsIterHelper([]int{},t);
}

//func TestMapElems(t *testing.T){
//
//}

func testChanIterHelper(chanNum int, t *testing.T){
    c:=make(chan int);
    go func(c chan int, numElems int){
        for i:=0; i<numElems; i++ {
            c <- i;
        }
        close(c);
    }(c,chanNum);
    cnt,err:=ChanElems(c).Count();
    test.BasicTest(chanNum,cnt,
        "ChanElems did not get proper numner of values",t,
    );
    test.BasicTest(nil,err,
        "ChanElems returned an error when it was not supposed to.",t,
    );
}
func TestChanElems(t *testing.T){
    testChanIterHelper(0,t);
    testChanIterHelper(1,t);
    testChanIterHelper(5,t);
    testChanIterHelper(20,t);
}

func testFileLinesHelper(numLines int, path string, t *testing.T){
    fIter:=FileLines(fmt.Sprintf("./testData/%s",path));
    for i:=0; i<numLines; i++ {
        fV,fErr,fBool:=fIter(Continue);
        test.BasicTest(fmt.Sprintf("%d",i+1),fV,
            "SliceElems iteration does not match actual values.",t,
        );
        test.BasicTest(nil,fErr,
            "SliceElems iteration produced an error when it shouldn't have.",t,
        );
        test.BasicTest(true,fBool,
            "SliceElems iteration did not stop when it should have.",t,
        );
    }
    fV,fErr,fBool:=fIter(Continue);
    test.BasicTest("",fV,
        "SliceElems iteration does not match actual values.",t,
    );
    test.BasicTest(nil,fErr,
        "SliceElems iteration produced an error when it shouldn't have.",t,
    );
    test.BasicTest(false,fBool,
        "SliceElems iteration did not stop when it should have.",t,
    );
}
func TestFileLines(t *testing.T){
    testFileLinesHelper(0,"emptyFile.txt",t);
    testFileLinesHelper(1,"oneLine.txt",t);
    testFileLinesHelper(3,"threeLines.txt",t);
}
