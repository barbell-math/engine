package iter;

import (
    "fmt"
	"errors"
	"testing"
	"github.com/barbell-math/block/util/test"
)

func forEachIterHelper[T any](
        vals []T,
        op func(index int, val T) T,
        t *testing.T){
    i:=0;
    cpy:=append([]T{},vals...);
    err:=SliceElems(vals).ForEach(
    func(index int, val T) (IteratorFeedback,error) {
        vals[i]=op(index,val);
        i++;
        return Continue,nil;
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

func TestForEachEarlyStopBool(t *testing.T){
    cntr:=0;
    err:=SliceElems([]int{0,1,2,3,4}).ForEach(
    func(i int, v int) (IteratorFeedback,error) {
        cntr++;
        if v==3 {
            return Break,nil;
        }
        return Continue,nil;
    });
    test.BasicTest(4,cntr,"ForEach did not stop after signaling to stop.",t);
    test.BasicTest(nil,err,
        "ForEach returned an error when it was not supposed to.",t,
    );
}

func TestForEachEarlyStopErr(t *testing.T){
    cntr:=0;
    err:=SliceElems([]int{0,1,2,3,4}).ForEach(
    func(i int, v int) (IteratorFeedback,error) {
        cntr++;
        if v==3 {
            return Continue,errors.New("NEW ERROR");
        }
        return Continue,nil;
    });
    test.BasicTest(4,cntr,"ForEach did not stop after signaling to stop.",t);
    if err==nil {
        fmt.Println(
            "ForEach did not return an error when it was supposed to.",t,
        );
    }
}
