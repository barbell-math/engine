package iter

import (
    "fmt"
	"errors"
	"testing"

	"github.com/barbell-math/block/util/test"
)

func collectIterHelper[T any](vals []T, t *testing.T){
    collected,err:=SliceElems(vals).Collect();
    test.SlicesMatch(vals,collected,t);
    test.BasicTest(nil,err,
        "Collect returned and error when it was not supposed to.",t,
    );
}
func TestCollect(t *testing.T){
    collectIterHelper([]int{1,2,3,4},t);
    collectIterHelper([]int{1},t);
    collectIterHelper([]int{},t);
}

func collectIntoIterHelper[T any](
        vals []T,
        buff []T,
        expectedNumChange int,
        t *testing.T){
    origLen:=len(buff);
    rv,err:=SliceElems(vals).CollectInto(buff);
    test.BasicTest(origLen,len(buff),
        "Buffer length changed when it shouldn't have.",t,
    );
    test.BasicTest(expectedNumChange,rv,
        "Total number of elements changed is not correct",t,
    );
    test.BasicTest(nil,err,
        "CollectInto returned and error when it was not supposed to.",t,
    );
    min:=len(buff);
    if len(vals)<len(buff) {
        min=len(vals);
    }
    for i:=0; i<min; i++ {
        test.BasicTest(vals[i],buff[i],fmt.Sprintf(
            "Values do not match | Index: %d",i,
        ),t);
    }
}
func TestCollectInto(t *testing.T){
    collectIntoIterHelper([]int{1,2,3,4},make([]int,5),4,t);
    collectIntoIterHelper([]int{1,2,3,4},make([]int,4),4,t);
    collectIntoIterHelper([]int{1,2,3,4},make([]int,3),3,t);
    collectIntoIterHelper([]int{1},make([]int,1),1,t);
    collectIntoIterHelper([]int{},make([]int,0),0,t);
}

func appendIterHelper[T any](orig []T, vals []T, t *testing.T){
    origLen:=len(orig);
    expLen:=len(orig)+len(vals);
    tmp,err:=SliceElems(vals).AppendTo(&orig);
    test.BasicTest(expLen,len(orig),
        "Append did not append correct number of elements.",t,
    );
    test.BasicTest(tmp,len(vals),
        "Append did not append correct number of elements.",t,
    );
    test.BasicTest(nil,err,
        "AppendTo returned an error when it was not supposed to.",t,
    );
    for i:=origLen; i<len(orig); i++ {
        test.BasicTest(vals[i-origLen],orig[i],
            "Append did not append correct values.",t,
        );
    }
}
func TestAppendTo(t *testing.T){
    appendIterHelper([]int{1,2,3,4},[]int{1,2,3,4},t);
    appendIterHelper([]int{1,2,3,4},[]int{1},t);
    appendIterHelper([]int{1,2,3,4},[]int{},t);
    appendIterHelper([]int{1},[]int{1,2,3,4},t);
    appendIterHelper([]int{},[]int{1,2,3,4},t);
    appendIterHelper([]int{},[]int{},t);
}

func TestAll(t *testing.T){
    res,err:=SliceElems([]int{1,2,3,4}).All(func(val int) (bool,error) {
        return val>0,nil;
    });
    test.BasicTest(true,res,"All did not return correct result.",t);
    test.BasicTest(nil,err,
        "All returned an error when it was not supposed to.",t,
    );
    res,err=SliceElems([]int{1,2,3,4}).All(func(val int) (bool,error) {
        return val<0,nil;
    });
    test.BasicTest(false,res,"All did not return correct result.",t);
    test.BasicTest(nil,err,
        "All returned an error when it was not supposed to.",t,
    );
    res,err=SliceElems([]int{1,2,3,4}).All(func(val int) (bool,error) {
        return val<2,nil;
    });
    test.BasicTest(false,res,"All did not return correct result.",t);
    test.BasicTest(nil,err,
        "All returned an error when it was not supposed to.",t,
    );
}

func TestAny(t *testing.T){
    found,err:=SliceElems([]int{1,2,3,4}).Any(func(val int) (bool,error) {
        return val>0,nil;
    });
    test.BasicTest(true,found,"Any did not return the correct result.",t);
    test.BasicTest(nil,err,
        "Any returned an error when it was not supposed to.",t,
    );
    found,err=SliceElems([]int{1,2,3,4}).Any(func(val int) (bool,error) {
        return val<2,nil;
    });
    test.BasicTest(true,found,"Any did not return the correct result.",t);
    test.BasicTest(nil,err,
        "Any returned an error when it was not supposed to.",t,
    );
    found,err=SliceElems([]int{1,2,3,4}).Any(func(val int) (bool,error) {
        return val<0,nil;
    });
    test.BasicTest(false,found,"Any did not return the correct result.",t);
    test.BasicTest(nil,err,
        "Any returned an error when it was not supposed to.",t,
    );
}

func findIterHelperFound[T comparable](elems []T, lookingFor T, t *testing.T){
    v,err,ok:=SliceElems(elems).Find(func(val T) (bool,error) {
        return val==lookingFor,nil;  
    });
    test.BasicTest(lookingFor,v,
        "Find did not find the element when it was supposed to.",t,
    );
    test.BasicTest(true,ok,
        "Find did not find the element when it was supposed to.",t,
    );
    test.BasicTest(nil,err,
        "Find returned an error when it was not supposed to.",t,
    );
}
func findIterHelperNotFound[T comparable](elems []T, lookingFor T, t *testing.T){
    var tmp T;
    v,err,ok:=SliceElems(elems).Find(func(val T) (bool,error) {
        return val==lookingFor,nil;  
    });
    test.BasicTest(tmp,v,
        "Find found the element when it was not supposed to.",t,
    );
    test.BasicTest(false,ok,
        "Find found the element when it was not supposed to.",t,
    );
    test.BasicTest(nil,err,
        "Find returned an error when it was not supposed to.",t,
    );
}
func TestFind(t *testing.T){
    findIterHelperFound([]int{1,2,3,4},1,t);
    findIterHelperFound([]int{1,2,3,4},4,t);
    findIterHelperNotFound([]int{1,2,3,4},5,t);
    findIterHelperFound([]int{1},1,t);
    findIterHelperNotFound([]int{1},5,t);
}

func indexIterHelper[T comparable](elems []T,
        lookingFor T,
        expectedIndex int,
        t *testing.T) {
    v,err:=SliceElems(elems).Index(func(val T) (bool,error) {
        return val==lookingFor,nil;
    });
    test.BasicTest(expectedIndex,v,"Index returned incorrect value.",t);
    test.BasicTest(nil,err,
        "Index returned an error when it was not supposed to.",t,
    );
}
func TestIndex(t *testing.T){
    indexIterHelper([]int{1,2,3,4},1,0,t);
    indexIterHelper([]int{1,2,3,4},3,2,t);
    indexIterHelper([]int{1,2,3,4},4,3,t);
    indexIterHelper([]int{1,2,3,4},5,-1,t);
    indexIterHelper([]int{1},1,0,t);
    indexIterHelper([]int{1},2,-1,t);
    indexIterHelper([]int{},1,-1,t);
}

func TestIndexErrorFound(t *testing.T){
    cntr:=0;
    v,err:=SliceElems([]int{1,2,3,4}).Index(func(val int) (bool,error) {
        cntr++;
        if val==3 {
            return true,errors.New("");
        }
        return false,nil;
    });
    test.BasicTest(3,cntr,"Index did not respect early exit.",t);
    test.BasicTest(2,v,"Index returned incorrect value.",t);
    if err==nil {
        test.FormatError("!nil",err,
            "Error was nil when it was not supposed to be.",t,
        );
    }
}

func TestIndexErrorNotFound(t *testing.T){
    cntr:=0;
    v,err:=SliceElems([]int{1,2,3,4}).Index(func(val int) (bool,error) {
        cntr++;
        if val==3 {
            return false,errors.New("");
        }
        return false,nil;
    });
    test.BasicTest(3,cntr,"Index did not respect early exit.",t);
    test.BasicTest(-1,v,"Index returned incorrect value.",t);
    if err==nil {
        test.FormatError("!nil",err,
            "Error was nil when it was not supposed to be.",t,
        );
    }
}

func nthIterHelper[T any](
        vals []T,
        index int,
        expectedVal T,
        expectedError bool,
        t *testing.T){
    val,err,ok:=SliceElems(vals).Nth(index);
    test.BasicTest(expectedVal,val,"Nth returned wrong value.",t);
    test.BasicTest(expectedError,ok,"Nth returned wrong error flag.",t);
    test.BasicTest(nil,err,
        "Nth returned an error when it was not suppoed to.",t,
    );
}
func TestNth(t *testing.T) {
    nthIterHelper([]int{1,2,3,4},0,1,true,t);
    nthIterHelper([]int{1,2,3,4},2,3,true,t);
    nthIterHelper([]int{1,2,3,4},3,4,true,t);
    nthIterHelper([]int{1,2,3,4},4,0,false,t);
    nthIterHelper([]int{1},0,1,true,t);
    nthIterHelper([]int{1},1,0,false,t);
    nthIterHelper([]int{},0,0,false,t);
}

func TestCount(t *testing.T){
    c,err:=SliceElems([]int{1,2,3,4}).Count();
    test.BasicTest(4,c,"Count did not return correct value.",t);
    test.BasicTest(nil,err,
        "Count returned an error when it was not supposed to.",t,
    );
    c,err=SliceElems([]int{1}).Count();
    test.BasicTest(1,c,"Count did not return correct value.",t);
    test.BasicTest(nil,err,
        "Count returned an error when it was not supposed to.",t,
    );
    c,err=SliceElems([]int{}).Count();
    test.BasicTest(0,c,"Count did not return correct value.",t);
    test.BasicTest(nil,err,
        "Count returned an error when it was not supposed to.",t,
    );
}

func toChanIterHelper[T any](vals []T, t *testing.T) {
    cntr:=0;
    c:=make(chan T);
    go func(){
        for val:=range(c) {
            test.BasicTest(cntr,val,"ToChan returned incorrect values.",t);
            cntr++;
        }
    }()
    SliceElems(vals).ToChan(c);
    close(c);
}
func TestToChan(t *testing.T) {
    vals:=make([]int,200);
    for i:=0; i<200; i++ {
        vals[i]=i;
    }
    toChanIterHelper([]int{},t);
    toChanIterHelper([]int{0},t);
    toChanIterHelper(vals,t);
}

func TestReduce(t *testing.T){
    i:=0;
    newErr:=fmt.Errorf("NEW ERROR");
    tmp,err:=SliceElems([]int{1,2,3,4}).Reduce(0,func(accum *int, iter int) error {
        *accum=*accum+iter;
        return nil;
    });
    test.BasicTest(10,tmp,"Reduce did not return appropriate value.",t);
    test.BasicTest(nil,err,"Reduce did not return appropriate error.",t);
    tmp,err=SliceElems([]int{1,2,3,4}).Reduce(0,func(accum *int, iter int) error {
        i++;
        return newErr;
    });
    test.BasicTest(1,i,"Reduce did not short circuit properly.",t);
    test.BasicTest(newErr,err,"Reduce did not return appropriate error.",t);
}

