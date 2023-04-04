package algo;

import (
    "fmt"
    "testing"
    "github.com/barbell-math/block/util/test"
)

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







