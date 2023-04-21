package iter

import (
    "fmt"
    "testing"
    "github.com/barbell-math/block/util/test"
    //"github.com/barbell-math/block/util/dataStruct"
    //"github.com/barbell-math/block/util/dataStruct/types"
)

func TestTake(t *testing.T) {
    cnt,err:=SliceElems([]int{1, 2, 3, 4}).Take(0).Count()
    test.BasicTest(0,cnt,"Take took the wrong number of items.",t);
    test.BasicTest(nil,err,"Take took the wrong number of items.",t);
    cnt,err=SliceElems([]int{1, 2, 3, 4}).Take(1).Count()
    test.BasicTest(1,cnt,"Take took the wrong number of items.",t);
    test.BasicTest(nil,err,"Take took the wrong number of items.",t);
	cnt,err=SliceElems([]int{1, 2, 3, 4}).Take(2).Count()
    test.BasicTest(2,cnt,"Take took the wrong number of items.",t);
    test.BasicTest(nil,err,"Take took the wrong number of items.",t);
	cnt,err=SliceElems([]int{1, 2, 3, 4}).Take(4).Count()
    test.BasicTest(4,cnt,"Take took the wrong number of items.",t);
    test.BasicTest(nil,err,"Take took the wrong number of items.",t);
	cnt,err=SliceElems([]int{1, 2, 3, 4}).Take(5).Count()
    test.BasicTest(4,cnt,"Take took the wrong number of items.",t);
    test.BasicTest(nil,err,"Take took the wrong number of items.",t);
	cnt,err=SliceElems([]int{}).Take(1).Count()
    test.BasicTest(0,cnt,"Take took the wrong number of items.",t);
    test.BasicTest(nil,err,"Take took the wrong number of items.",t);
	cnt,err=SliceElems([]int{}).Take(0).Count()
    test.BasicTest(0,cnt,"Take took the wrong number of items.",t);
    test.BasicTest(nil,err,"Take took the wrong number of items.",t);
}

func TestTakeWhile(t *testing.T){
    cnt,err:=SliceElems([]int{1,2,3,4}).TakeWhile(func(val int) bool {
        return val<3;
    }).Count();
    test.BasicTest(2,cnt,"TakeWhile did not take correct number of elements.",t);
    test.BasicTest(nil,err,"TakeWhile did not take correct number of elements.",t);
    cnt,err=SliceElems([]int{1,2,3,4}).TakeWhile(func(val int) bool {
        return val<1;
    }).Count();
    test.BasicTest(0,cnt,"TakeWhile did not take correct number of elements.",t);
    test.BasicTest(nil,err,"TakeWhile did not take correct number of elements.",t);
    cnt,err=SliceElems([]int{1,2,3,4}).TakeWhile(func(val int) bool {
        return val<2;
    }).Count();
    test.BasicTest(1,cnt,"TakeWhile did not take correct number of elements.",t);
    test.BasicTest(nil,err,"TakeWhile did not take correct number of elements.",t);
    cnt,err=SliceElems([]int{1,2,3,4}).TakeWhile(func(val int) bool {
        return val<5;
    }).Count();
    test.BasicTest(4,cnt,"TakeWhile did not take correct number of elements.",t);
    test.BasicTest(nil,err,"TakeWhile did not take correct number of elements.",t);
}

func TestSkip(t *testing.T) {
    cnt,err:=SliceElems([]int{1,2,3,4,5,6,7,8,9}).Skip(0).Count();
    test.BasicTest(9,cnt,"Skip skipped a value it should not have.",t);
    test.BasicTest(nil,err,"Skip returned an error when it should not have.",t);
    cnt,err=SliceElems([]int{1,2,3,4,5,6,7,8,9}).Skip(1).Count();
    test.BasicTest(8,cnt,"Skip skipped a value it should not have.",t);
    test.BasicTest(nil,err,"Skip returned an error when it should not have.",t);
    cnt,err=SliceElems([]int{1,2,3,4,5,6,7,8,9}).Skip(2).Count();
    test.BasicTest(7,cnt,"Skip skipped a value it should not have.",t);
    test.BasicTest(nil,err,"Skip returned an error when it should not have.",t);
    cnt,err=SliceElems([]int{1,2,3,4,5,6,7,8,9}).Skip(8).Count();
    test.BasicTest(1,cnt,"Skip skipped a value it should not have.",t);
    test.BasicTest(nil,err,"Skip returned an error when it should not have.",t);
    cnt,err=SliceElems([]int{1,2,3,4,5,6,7,8,9}).Skip(9).Count();
    test.BasicTest(0,cnt,"Skip skipped a value it should not have.",t);
    test.BasicTest(nil,err,"Skip returned an error when it should not have.",t);
    cnt,err=SliceElems([]int{1,2,3,4,5,6,7,8,9}).Skip(10).Count();
    test.BasicTest(0,cnt,"Skip skipped a value it should not have.",t);
    test.BasicTest(nil,err,"Skip returned an error when it should not have.",t);
}

func mapIterHelper[T any](elems []T, t *testing.T){
    mapped,err:=Map(SliceElems(elems),func(index int, val T) (string,error) {
        return fmt.Sprintf("%v",val),nil;
    }).Collect();
    test.BasicTest(nil,err,"Map returned an error when it should not have",t);
    for i,v:=range(elems) {
        test.BasicTest(fmt.Sprintf("%v",v),mapped[i],
            "Mapping did not mutate elements as expected.",t,
        );
    }
}
func TestMap(t *testing.T){
    mapIterHelper([]int{1,2,3,4},t);
    mapIterHelper([]int{1},t);
    mapIterHelper([]int{},t);
}

func TestFilter(t *testing.T){
    cntr,err:=SliceElems([]int{1,2,3,4}).Filter(func(val int) bool {
        return val<3;
    }).Count();
    test.BasicTest(2,cntr,"Filter did not work appropriately.",t);
    test.BasicTest(nil,err,"Filter returned an error when it should not have",t);
    cntr,err=SliceElems([]int{1,2,3,4}).Filter(func(val int) bool {
        return val<5;
    }).Count();
    test.BasicTest(4,cntr,"Filter did not work appropriately.",t);
    test.BasicTest(nil,err,"Filter returned an error when it should not have",t);
    cntr,err=SliceElems([]int{1,2,3,4}).Filter(func(val int) bool {
        return val<2;
    }).Count();
    test.BasicTest(1,cntr,"Filter did not work appropriately.",t);
    test.BasicTest(nil,err,"Filter returned an error when it should not have",t);
    cntr,err=SliceElems([]int{1,2,3,4}).Filter(func(val int) bool {
        return val<1;
    }).Count();
    test.BasicTest(0,cntr,"Filter did not work appropriately.",t);
    test.BasicTest(nil,err,"Filter returned an error when it should not have",t);
}

//func testWindowIterHelper[T any](vals []T, size int, t *testing.T){
//    q,_:=dataStruct.NewCircularQueue[T](size);
//    Window(SliceElems(vals),q,false,
//    func(index int, q *types.Queue[T]) (IteratorFeedback, error) {
//
//    }).ForEach()
//}
//func TestWindow(t *testing.T){
//
//}
