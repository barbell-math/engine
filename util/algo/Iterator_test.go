package algo;

import (
    "fmt"
    "testing"
    "github.com/barbell-math/block/util/test"
)

func valElemIterHelper[T any](val T, err error, t *testing.T){
    var tmp T;
    iter:=ValElem(val,err);
    for i:=0; i<3; i++ {
        vIter,eIter,contIter:=iter();
        if i==0 {
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
    valElemIterHelper(1,nil,t);
    valElemIterHelper(2,fmt.Errorf("NEW ERROR"),t);
}

func sliceElemsIterHelper[T any](vals []T, t *testing.T){
    sIter:=SliceElems(vals);
    for i:=0; i<len(vals); i++ {
        sV,sErr,sBool:=sIter();
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
    sV,sErr,sBool:=sIter();
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
    test.BasicTest(chanNum,ChanElems(c).Count(),
        "ChanElems did not get proper numner of values",t,
    );
}
func TestChanElems(t *testing.T){
    testChanIterHelper(0,t);
    testChanIterHelper(1,t);
    testChanIterHelper(5,t);
    testChanIterHelper(20,t);
}

func TestMap(t *testing.T){
    s:=[]int{1,2,3,4};
    mapped:=Map(SliceElems(s),func(val int) (string,error) {
        return fmt.Sprintf("%d",val),nil;
    }).Collect();
    for i,v:=range(s) {
        test.BasicTest(fmt.Sprintf("%d",v),mapped[i],
            "Mapping did not mutate elements as expected.",t,
        );
    }
}

func forEachIterHelper[T any](
        vals []T,
        op func(index int, val T) T,
        t *testing.T){
    i:=0;
    cpy:=append([]T{},vals...);
    err:=SliceElems(vals).ForEach(func(index int, val T) error {
        vals[i]=op(index,val);
        i++;
        return nil;
    });
    test.BasicTest(nil,err,"ForEach returned an error when it shouldn't have.",t);
    test.BasicTest(len(cpy),len(vals),"ForEach changed size of slice.",t);
    for i,v:=range(cpy) {
        test.BasicTest(op(i,v),vals[i],"ForEach did not mutate elements properly.",t);
    }
}
func TestForEach(t *testing.T){
    forEachIterHelper([]int{1,2,3,4},func(index int, val int) int {
        return val+1;
    },t);
    forEachIterHelper([]int{1},func(index int, val int) int {
        return val+1;
    },t);
    forEachIterHelper([]int{},func(index int, val int) int {
        return val+1;
    },t);
    forEachIterHelper([]int{5,5,5,5},func(index int, val int) int {
        return index+1;
    },t);
}

func TestFilterVal(t *testing.T){
    test.BasicTest(2,SliceElems([]int{1,2,3,4}).
    Filter(func(val int) bool {
        return val<3;
    }).Count(),"Filter did not work appropriately.",t);
    test.BasicTest(4,SliceElems([]int{1,2,3,4}).
    Filter(func(val int) bool {
        return val<5;
    }).Count(),"Filter did not work appropriately.",t);
    test.BasicTest(1,SliceElems([]int{1,2,3,4}).
    Filter(func(val int) bool {
        return val<2;
    }).Count(),"Filter did not work appropriately.",t);
    test.BasicTest(0,SliceElems([]int{1,2,3,4}).
    Filter(func(val int) bool {
        return val<1;
    }).Count(),"Filter did not work appropriately.",t);
}

func TestReduce(t *testing.T){
    i:=0;
    newErr:=fmt.Errorf("NEW ERROR");
    tmp,err:=SliceElems([]int{1,2,3,4}).Reduce(func(accum *int, iter int) error {
        *accum=*accum+iter;
        return nil;
    },0);
    test.BasicTest(10,tmp,"Reduce did not return appropriate value.",t);
    test.BasicTest(nil,err,"Reduce did not return appropriate error.",t);
    tmp,err=SliceElems([]int{1,2,3,4}).Reduce(func(accum *int, iter int) error {
        i++;
        return newErr;
    },0);
    test.BasicTest(1,i,"Reduce did not short circuit properly.",t);
    test.BasicTest(newErr,err,"Reduce did not return appropriate error.",t);
}

func TestTake(t *testing.T){
    test.BasicTest(2,SliceElems([]int{1,2,3,4}).Take(2).Count(),
        "Take took more items than it should have.",t,
    );
    test.BasicTest(4,SliceElems([]int{1,2,3,4}).Take(4).Count(),
        "Take took more items than it should have.",t,
    );
    test.BasicTest(4,SliceElems([]int{1,2,3,4}).Take(5).Count(),
        "Take took more items than it should have.",t,
    );
    test.BasicTest(0,SliceElems([]int{1,2,3,4}).Take(0).Count(),
        "Take took more items than it should have.",t,
    );
    test.BasicTest(1,SliceElems([]int{1,2,3,4}).Take(1).Count(),
        "Take took more items than it should have.",t,
    );
}

func TestTakeWhile(t *testing.T){
    test.BasicTest(2,SliceElems([]int{1,2,3,4}).
    TakeWhile(func(val int, err error) bool {
        return val<3;
    }).Count(),"TakeWhile did not take correct number of elements.",t);
    test.BasicTest(0,SliceElems([]int{1,2,3,4}).
    TakeWhile(func(val int, err error) bool {
        return val<1;
    }).Count(),"TakeWhile did not take correct number of elements.",t);
    test.BasicTest(1,SliceElems([]int{1,2,3,4}).
    TakeWhile(func(val int, err error) bool {
        return val<2;
    }).Count(),"TakeWhile did not take correct number of elements.",t);
    test.BasicTest(4,SliceElems([]int{1,2,3,4}).
    TakeWhile(func(val int, err error) bool {
        return val<5;
    }).Count(),"TakeWhile did not take correct number of elements.",t);
}

func collectIterHelper[T any](vals []T, t *testing.T){
    collected:=SliceElems(vals).Collect();
    test.SlicesMatch(vals,collected,t);
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
    rv:=SliceElems(vals).CollectInto(buff);
    test.BasicTest(origLen,len(buff),
        "Buffer length changed when it shouldn't have.",t,
    );
    test.BasicTest(expectedNumChange,rv,
        "Total number of elements changed is not correct",t,
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
    tmp:=SliceElems(vals).AppendTo(&orig);
    test.BasicTest(expLen,len(orig),
        "Append did not append correct number of elements.",t,
    );
    test.BasicTest(tmp,len(vals),
        "Append did not append correct number of elements.",t,
    );
    for i:=origLen; i<len(orig); i++ {
        test.BasicTest(vals[i-origLen],orig[i],
            "Append did not append correct values.",t,
        );
    }
}
func TestAppend(t *testing.T){
    appendIterHelper([]int{1,2,3,4},[]int{1,2,3,4},t);
    appendIterHelper([]int{1,2,3,4},[]int{1},t);
    appendIterHelper([]int{1,2,3,4},[]int{},t);
    appendIterHelper([]int{1},[]int{1,2,3,4},t);
    appendIterHelper([]int{},[]int{1,2,3,4},t);
    appendIterHelper([]int{},[]int{},t);
}

func TestAll(t *testing.T){
    test.BasicTest(true,SliceElems([]int{1,2,3,4}).All(func(val int) bool {
        return val>0;
    }),"All did not return correct result.",t);
    test.BasicTest(false,SliceElems([]int{1,2,3,4}).All(func(val int) bool {
        return val<0;
    }),"All did not return correct result.",t);
    test.BasicTest(false,SliceElems([]int{1,2,3,4}).All(func(val int) bool {
        return val<2;
    }),"All did not return correct result.",t);
}

func TestAny(t *testing.T){
    test.BasicTest(true,SliceElems([]int{1,2,3,4}).Any(func(val int) bool {
        return val>0;
    }),"All did not return correct result.",t);
    test.BasicTest(false,SliceElems([]int{1,2,3,4}).Any(func(val int) bool {
        return val<0;
    }),"All did not return correct result.",t);
    test.BasicTest(true,SliceElems([]int{1,2,3,4}).Any(func(val int) bool {
        return val<2;
    }),"All did not return correct result.",t);
}

func TestCount(t *testing.T){
    test.BasicTest(4,SliceElems([]int{1,2,3,4}).Count(),
        "Count did not return correct value.",t,
    );
    test.BasicTest(1,SliceElems([]int{1}).Count(),
        "Count did not return correct value.",t,
    );
    test.BasicTest(0,SliceElems([]int{}).Count(),
        "Count did not return correct value.",t,
    );
}

func TestFind(t *testing.T){
    v,ok:=SliceElems([]int{1,2,3,4}).Find(func(val int) bool {
        return val==1;  
    });
    test.BasicTest(1,v,"Find did not find the element when it was supposed to.",t);
    test.BasicTest(true,ok,
        "Find did not find the element when it was supposed to.",t,
    );
    v,ok=SliceElems([]int{1,2,3,4}).Find(func(val int) bool {
        return val==4;  
    });
    test.BasicTest(4,v,"Find did not find the element when it was supposed to.",t);
    test.BasicTest(true,ok,
        "Find did not find the element when it was supposed to.",t,
    );
    v,ok=SliceElems([]int{1,2,3,4}).Find(func(val int) bool {
        return val==5;  
    });
    test.BasicTest(0,v,"Find found the element when it was not supposed to.",t);
    test.BasicTest(false,ok,
        "Find found the element when it was not supposed to.",t,
    );
    v,ok=SliceElems([]int{1}).Find(func(val int) bool {
        return val==1;  
    });
    test.BasicTest(1,v,"Find did not find the element when it was supposed to.",t);
    test.BasicTest(true,ok,
        "Find did not find the element when it was supposed to.",t,
    );
    v,ok=SliceElems([]int{}).Find(func(val int) bool {
        return val==5;  
    });
    test.BasicTest(0,v,"Find found the element when it was not supposed to.",t);
    test.BasicTest(false,ok,
        "Find found the element when it was not supposed to.",t,
    );
}

func TestIndex(t *testing.T){
    test.BasicTest(0,SliceElems([]int{1,2,3,4}).Index(func(val int) bool {
        return val==1;
    }),"Index returned incorrect value.",t);
    test.BasicTest(2,SliceElems([]int{1,2,3,4}).Index(func(val int) bool {
        return val==3;
    }),"Index returned incorrect value.",t);
    test.BasicTest(3,SliceElems([]int{1,2,3,4}).Index(func(val int) bool {
        return val==4;
    }),"Index returned incorrect value.",t);
    test.BasicTest(-1,SliceElems([]int{1,2,3,4}).Index(func(val int) bool {
        return val==5;
    }),"Index returned incorrect value.",t);
    test.BasicTest(0,SliceElems([]int{1}).Index(func(val int) bool {
        return val==1;
    }),"Index returned incorrect value.",t);
    test.BasicTest(-1,SliceElems([]int{1}).Index(func(val int) bool {
        return val==2;
    }),"Index returned incorrect value.",t);
    test.BasicTest(-1,SliceElems([]int{}).Index(func(val int) bool {
        return val==1;
    }),"Index returned incorrect value.",t);
}

func nthIterHelper[T any](
        vals []T,
        index int,
        expectedVal T,
        expectedError bool,
        t *testing.T){
    val,ok:=SliceElems(vals).Nth(index);
    test.BasicTest(expectedVal,val,"Nth returned wrong value.",t);
    test.BasicTest(expectedError,ok,"Nth returned wrong error flag.",t);
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
