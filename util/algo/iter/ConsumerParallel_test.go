package iter;

import (
	"testing"
	"github.com/barbell-math/block/util/test"
    customerr "github.com/barbell-math/block/util/err"
)

func parallelIterHelper(vals []int, numThreads int, t *testing.T){
    i:=0;
    cpy:=make([]int,len(vals));
    rv:=SliceElems(vals).Parallel(func(val int) (int,error) {
        return val+1,nil;
    },func(val int, res int, err error){
        cpy[i]=res+1;
        i++;
    },numThreads);
    test.BasicTest(nil,rv,
        "ForEachParallel returned an error when it shouldn't have.",t,
    );
    for i,v:=range(cpy) {
        test.BasicTest(vals[i]+2,v,
            "ForEachParallel did not run correct operations.",t,
        );
    }
}
func TestParallel(t *testing.T) {
    rv:=SliceElems([]int{1,2,3,4}).Parallel(func(val int) (int,error) {
        return 0,nil;
    },NoOp[int,int],0);
    if !customerr.IsValOutsideRange(rv) {
        test.FormatError(customerr.ValOutsideRange(""),rv,
            "ForEachParallel returned incorrect error when one was expected.",t,
        );
    }
    vals:=make([]int,200);
    for i:=0; i<200; i++ {
        vals[i]=i;
    }
    for _,i:=range([]int{1,25,50,75,100}) {
        parallelIterHelper([]int{},i,t);
        parallelIterHelper([]int{1},i,t);
        parallelIterHelper(vals,i,t);
    }
}
