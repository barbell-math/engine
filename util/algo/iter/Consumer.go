package iter;

import (
//    "fmt"
    customerr "github.com/barbell-math/block/util/err"
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
