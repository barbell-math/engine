package iter;

import (
//    "fmt"
    customerr "github.com/barbell-math/engine/util/err"
)

//ForEach is the only true consumer, all other consumers can be expressed using
//ForEach making them pseudo-consumers. By using this pattern all pseudo-consumers
//are abstracted away from the complex looping logic.
func (i Iter[T])ForEach(
        op func(index int, val T) (IteratorFeedback,error)) error {
    j:=0;
    f:=Continue;
    var next T;
    var err error;
    var cont bool=true;
    var opErr error=nil;
    for cont && err==nil && opErr==nil && f==Continue {
        next,err,cont=i(Iterate);
        if err==nil && cont {
            f,opErr=op(j,next);
            j++;
        }
    }
    _,cleanUpErr,_:=i(Break);
    return customerr.AppendError(opErr,
        customerr.AppendError(err,cleanUpErr),
    );
}

//Why is stop not a pseudo consumer? It breaks the parent calling convention
// that ForEach uses. For each will always get a parent iterators value before
// the op function is consulted. Stop should just stop, and not call the parent
// iterators one last time meaning it has to be separate.
func (i Iter[T])Stop() error {
    _,cleanUpErr,_:=i(Break);
    return cleanUpErr;
}
