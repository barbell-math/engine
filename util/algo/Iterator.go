package algo

import (
	"fmt"

	"github.com/barbell-math/block/customerr"
)

func (i Iter[T])Filter(op Filter[T]) Iter[T] {
    return func(f IteratorFeedback) (T,error,bool) {
        var val T;
        if f==EarlyExit {
            i(EarlyExit);
            return val,nil,false;
        }
        var err error;
        var cont bool;
        for val,err,cont=i(f); err==nil && cont && !op(val); val,err,cont=i(f) {}
        return val,err,cont;
    }
}












