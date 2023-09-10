package iter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/barbell-math/block/util/test"
	customerr "github.com/barbell-math/block/util/err"
)

func TestNext(t *testing.T){
    //Should produce sequence:
    //  (1,nil,true)
    //  (3,nil,true)
    //  (5,nil,true)
    //  (5,err,false)
    n:=SliceElems([]int{1,2,3,4,5,6,7}).Next(
    func(index, val int, status IteratorFeedback) (IteratorFeedback, int, error) {
        if status==Break {
            return Break,0,nil;
        }
        if val%2==0 {
            return Continue,val+1,nil;
        } else if val==3 {
            return Iterate,val,nil;
        } else if val==5 {
            return Continue,val,errors.New("NEW ERROR");
        }
        return Continue,val,nil;
    });
    next,err,cont:=n(Iterate);
    test.BasicTest(1,next,"Next did not return the correct value.",t);
    test.BasicTest(nil,err,
        "Next returned an error when it was not supposed to.",t,
    );
    test.BasicTest(true,cont,"Next did not reutrn correct cont status.",t);
    next,err,cont=n(Iterate);
    test.BasicTest(3,next,"Next did not return the correct value.",t);
    test.BasicTest(nil,err,
        "Next returned an error when it was not supposed to.",t,
    );
    test.BasicTest(true,cont,"Next did not reutrn correct cont status.",t);
    next,err,cont=n(Iterate);
    test.BasicTest(5,next,"Next did not return the correct value.",t);
    test.BasicTest(nil,err,
        "Next returned an error when it was not supposed to.",t,
    );
    test.BasicTest(true,cont,"Next did not reutrn correct cont status.",t);
    next,err,cont=n(Iterate);
    test.BasicTest(5,next,"Next did not return the correct value.",t);
    if err==nil {
        test.FormatError("!nil",err,
            "Next did not return an error when it was supposed to.",t,
        );
    }
    test.BasicTest(false,cont,"Next did not reutrn correct cont status.",t);
    next,err,cont=n(Break);
    test.BasicTest(0,next,"Next did not return the correct value.",t);
    test.BasicTest(nil,err,
        "Next returned an error when it was not supposed to.",t,
    );
    test.BasicTest(false,cont,"Next did not reutrn correct cont status.",t);
}

func TestNextReachesBreak(t *testing.T){
    breakReached:=false;
    breakReached2:=false;
    n:=SliceElems([]int{1,2,3,4,5,6,7}).Next(
    func(index, val int, status IteratorFeedback) (IteratorFeedback, int, error) {
        if status==Break {
            breakReached=true;
        }
        return Continue,val,nil;
    }).Next(func(index, val int, status IteratorFeedback) (IteratorFeedback, int, error) {
        if status==Break {
            breakReached2=true;
        }
        return Continue,val,nil;
    });
    err:=n.Consume();
    test.BasicTest(true,breakReached,
        "Next did not properly call parrent iterators with break flag.",t,
    );
    test.BasicTest(true,breakReached2,
        "Next did not properly call parrent iterators with break flag.",t,
    );
    test.BasicTest(nil,err,
        "Next returned an error when it shouldn't have.",t,
    );
}

func TestNextReachesBreakParentErr(t *testing.T){
    expectedErr:=errors.New("NEW ERROR");
    breakReached:=false;
    breakReached2:=false;
    n:=SliceElems([]int{1,2,3,4,5,6,7}).Next(
    func(index, val int, status IteratorFeedback) (IteratorFeedback, int, error) {
        if status==Break {
            breakReached=true;
        }
        if val==5 {
            return Break,val,expectedErr;
        }
        return Continue,val,nil;
    }).Next(func(index, val int, status IteratorFeedback) (IteratorFeedback, int, error) {
        if status==Break {
            breakReached2=true;
        }
        return Continue,val,nil;
    });
    err:=n.Consume();
    test.BasicTest(true,breakReached,
        "Next did not properly call parrent iterators with break flag.",t,
    );
    test.BasicTest(true,breakReached2,
        "Next did not properly call parrent iterators with break flag.",t,
    );
    test.BasicTest(expectedErr,err,
        "Next returned an error when it shouldn't have.",t,
    );
}

func TestNextReachesBreakParentCleanUpErr(t *testing.T){
    expectedErr:=errors.New("NEW ERROR");
    breakReached:=false;
    breakReached2:=false;
    n:=SliceElems([]int{1,2,3,4,5,6,7}).Next(
    func(index, val int, status IteratorFeedback) (IteratorFeedback, int, error) {
        if status==Break {
            breakReached=true;
            return Continue,val,expectedErr;
        }
        return Continue,val,nil;
    }).Next(func(index, val int, status IteratorFeedback) (IteratorFeedback, int, error) {
        if status==Break {
            breakReached2=true;
        }
        return Continue,val,nil;
    });
    err:=n.Consume();
    test.BasicTest(true,breakReached,
        "Next did not properly call parrent iterators with break flag.",t,
    );
    test.BasicTest(true,breakReached2,
        "Next did not properly call parrent iterators with break flag.",t,
    );
    test.BasicTest(expectedErr,err,
        "Next returned an error when it shouldn't have.",t,
    );
}

func TestNextReachesBreakParentErrAndCleanUpErr(t *testing.T){
    expectedErr:=errors.New("NEW ERROR");
    breakReached:=false;
    breakReached2:=false;
    n:=SliceElems([]int{1,2,3,4,5,6,7}).Next(
    func(index, val int, status IteratorFeedback) (IteratorFeedback, int, error) {
        if status==Break {
            breakReached=true;
            return Continue,val,expectedErr;
        }
        return Continue,val,nil;
    }).Next(func(index, val int, status IteratorFeedback) (IteratorFeedback, int, error) {
        if status==Break {
            breakReached2=true;
        }
        if val==5 {
            return Break,val,expectedErr;
        }
        return Continue,val,nil;
    });
    err:=n.Consume();
    test.BasicTest(true,breakReached,
        "Next did not properly call parrent iterators with break flag.",t,
    );
    test.BasicTest(true,breakReached2,
        "Next did not properly call parrent iterators with break flag.",t,
    );
    tmp:=fmt.Sprintf("%s",customerr.AppendError(expectedErr,expectedErr));
    if fmt.Sprintf("%s",expectedErr)==tmp {
        test.FormatError(tmp,err,
            "Next returned an error when it shouldn't have.",t,
        );
    }
}

func injectIterHelper[T any](initialVals []T,
        desiredSeq []T, 
        op func(idx int, val T) (T,bool), t *testing.T){
    result,err:=SliceElems(initialVals).Inject(op).Collect();
    test.BasicTest(nil,err,"Inject created an error when it should not have.",t);
    test.BasicTest(len(desiredSeq),len(result),
        "Using inject did not return the correct number of elements.",t,
    );
    for i,v:=range(desiredSeq) {
        if i<len(result) {
            test.BasicTest(v,result[i],
                "Inject incorrectly modified values.",t,
            );
        }
    }
}
func TestInject(t *testing.T) {
    injectIterHelper([]int{},[]int{},
        func(idx, val int) (int, bool) { return 0,idx==1; },t,
    );
    injectIterHelper([]int{1,2,3,4},[]int{1,2,3,4},
        func(idx, val int) (int, bool) { return 0,idx==5; },t,
    );
    injectIterHelper([]int{},[]int{0},
        func(idx, val int) (int, bool) { return 0,idx==0; },t,
    );
    injectIterHelper([]int{1},[]int{0,1},
        func(idx, val int) (int, bool) { return 0,idx==0; },t,
    );
    injectIterHelper([]int{1,2,3,4},[]int{0,1,2,3,4},
        func(idx, val int) (int, bool) { return 0,idx==0; },t,
    );
    injectIterHelper([]int{1},[]int{1,0},
        func(idx, val int) (int, bool) { return 0,idx==1; },t,
    );
    injectIterHelper([]int{1,2,3,4},[]int{1,2,3,4,0},
        func(idx, val int) (int, bool) { return 0,idx==4; },t,
    );
    injectIterHelper([]int{1,2,3,4},[]int{1,2,0,3,4},
        func(idx, val int) (int, bool) { return 0,idx==2; },t,
    );
}
