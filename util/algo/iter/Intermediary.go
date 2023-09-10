package iter

import (
	customerr "github.com/barbell-math/engine/util/err"
)

//Intermediaries sit between producers and consumers. They are responsible for
//ensuring all iter feedback messages get passed up the call chain to from the
//source of the message to the producer. This allows for resources to be managed
//properly.
//Rules:
//  1. Errors are propagated down to the consumer. The consumer will then call
//     it's parent iterator with the Break flag.
//  2. All Break flags will be passed up to the producer. This allows resources
//     to be destroyed in a top down fashion.
//  3. If a Break flag is generated in a intermediary, it should not clean up
//     its parents, but should return the command to not continue. The consumer
//     will start the destruction process once it sees the command to not continue.

//Next is the only true intermediary, all other intermediaries can be
//expressed using Next making them pseudo-intermediaries. By using this pattern
//all pseudo-intermediaries are abstracted away from the complex looping logic
//and do not need to worry about iter feedback message passing.
func Next[T any, U any](i Iter[T],
        op func(index int, val T, status IteratorFeedback) (IteratorFeedback,U,error)) Iter[U] {
    j:=0;
    return func(f IteratorFeedback) (U,error,bool) {
        var tmp U;
        next,err,cont:=i(f);
        if f==Break {
            _,_,err2:=op(j,next,Break); //Clean up current iterator
            return tmp,customerr.AppendError(err,err2),false;
        }
        var opErr error=nil;
        for ; cont && err==nil && opErr==nil && f==Iterate; next,err,cont=i(f) {
            f,tmp,opErr=op(j,next,f);
            j++;
            if f==Continue || f==Break || opErr!=nil {
                return tmp,opErr,(opErr==nil && f==Continue);
            } 
        }
        return tmp,err,cont;
    }
}
func (i Iter[T])Next(
        op func(index int, val T, status IteratorFeedback) (IteratorFeedback,T,error)) Iter[T] {
    return Next(i,op);
}

func (i Iter[T])Inject(op func(idx int, val T) (T,bool)) Iter[T] {
    j:=-1;
    injected:=false;
    var prevErr error;
    var prevCont bool;
    var prevVal T;
    return func(f IteratorFeedback) (T,error,bool) {
        if f==Break {
            var tmp T;
            return tmp,nil,false;
        }
        j++;
        if injected {
            injected=false;
            return prevVal,prevErr,prevCont;
        }
        next,err,cont:=i(f);
        if v,status:=op(j,next); status {
            injected=true;
            prevVal=next;
            prevCont=cont;
            prevErr=err;
            return v,nil,true;
        } else {
            return next,err,cont;
        }
    }
}
