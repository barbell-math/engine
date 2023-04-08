package algo;

import (
    "fmt"
    "testing"
    "github.com/barbell-math/block/util/test"
)

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








