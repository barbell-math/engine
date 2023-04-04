package algo

import (
	"fmt"

	"github.com/barbell-math/block/customerr"
)

func Map[T any, U any](i Iter[T], op func(val T) (U,error)) Iter[U] {
    return func(f IteratorFeedback) (U,error,bool) {
        var tmp U;
        if f==EarlyExit {
            i(EarlyExit);
            return tmp,nil,false;
        }
        next,err,ok:=i(Continue);
        if !ok || err!=nil {
            return tmp,err,false;
        }
        mappedVal,err:=op(next);
        return mappedVal,err,true;
    }
}
func (i Iter[T])Map(op func(val T) (T,error)) Iter[T] {
    return Map(i,op);
}

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

func (i Iter[T])Take(num int) Iter[T] {
    cntr:=0;
    return func() (T,error,bool) {
        var tmp T;
        if cntr<num {
            cntr++;
            return i();
        }
        return tmp,nil,false;
    }
}

func (i Iter[T])TakeWhile(op func(val T, err error) bool) Iter[T] {
    stop:=false;
    return func() (T,error,bool) {
        var tmp T;
        if !stop {
            val,err,cont:=i();
            if stop=!op(val,err); !stop {
                return val,err,cont;
            }
        }
        return tmp,nil,false;
    }
}











