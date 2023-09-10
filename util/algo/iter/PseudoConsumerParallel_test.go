package iter;

import (
    "testing"
    "github.com/barbell-math/block/util/test"
    customerr "github.com/barbell-math/block/util/err"
)

func filterParallelHelper(vals []int, numThreads int, t *testing.T) {
    rv,err:=SliceElems(vals).FilterParallel(func(val int) bool {
        return val==1 || val==2;
    },numThreads);
    test.BasicTest(nil,err,
        "FilterParallel returned and error when it wasn't supposed to.",t,
    );
    for _,v:=range(rv) {
        if v!=1 && v!=2 {
            test.FormatError("1 | 2",v,
                "FilterParallel returned errors it was not supposed to.",t,
            );
        }
    }
}
func TestFilterParallel(t *testing.T) {
    _,err:=SliceElems([]int{1,2,3,4}).FilterParallel(func(val int) bool {
        return false;
    },0);
    if !customerr.IsValOutsideRange(err) {
        test.FormatError(customerr.ValOutsideRange(""),err,
            "FilterParallel returned incorrect error when one was expected.",t,
        );
    }
    vals:=make([]int,200);
    for i:=0; i<200; i++ {
        vals[i]=i;
    }
    for _,i:=range([]int{1,25,50,75,100}) {
        filterParallelHelper([]int{},i,t);
        filterParallelHelper([]int{1},i,t);
        filterParallelHelper(vals,i,t);
    }
}
