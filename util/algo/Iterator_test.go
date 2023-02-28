package algo;

import (
    "fmt"
    "testing"
    "github.com/barbell-math/block/util/test"
)

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

func TestMap(t *testing.T){
    s:=[]int{1,2,3,4};
    mapped:=Map(SliceElems(s),func(val int) (string,error) {
        return fmt.Sprintf("%d",val),nil;
    }).Collect();
    for i,v:=range(s) {
        test.BasicTest(fmt.Sprintf("%d",v),mapped[i],
            "Mapping did not mutate values as expected.",t,
        );
    }
}

func TestForEach(t *testing.T) {
    
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

func TestIndex(t *testing.T) {
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
