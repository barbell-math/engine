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

//func TestMap(t *testing.T){
//    sIter:=SliceElems([]int{1,2,3,4});
//    mapper:=Map(sIter,func(val int) (string,error) {
//        return fmt.Sprintf("%d",val),nil;
//    });
//    fmt.Printf("%+v\n",mapper.Collect());
//    fmt.Printf("%T\n",mapper.Collect());
//    tmp:=SliceElems([]int{1,2,3,4,5}).Map(func(val int) (string,error){
//        return fmt.Sprintf("%d",val),nil;
//    }).Collect();
//    fmt.Printf("%+v\n",tmp);
//    fmt.Printf("%T\n",tmp);
//}
