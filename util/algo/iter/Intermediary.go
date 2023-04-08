package iter

import (
    customerr "github.com/barbell-math/block/util/err"
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

//Next is the only true intermediary, all other intermediaries can be expressed
//using Next making them pseudo-intermediaries. By using this pattern all
//pseudo-intermediaries are abstracted away from the complex looping logic and
//do not need to worry about iter feedback message passing.
func Next[T any, U any](i Iter[T],
        op func(index int, val T, status IteratorFeedback) (IteratorFeedback,U,error)) Iter[U] {
    j:=0;
    return func(f IteratorFeedback) (U,error,bool) {
        var tmp U;
        next,err,cont:=i(f);
        if f==Break {
            _,_,err2:=op(j,next,Break); //Clean up current iterator
            if err!=nil && err2!=nil {
                return tmp,customerr.AppendError(err,err2),false;
            } else if err!=nil && err2==nil {
                return tmp,err,false;
            }
            return tmp,err2,false;
        }
        var opErr error=nil;
        for ; cont && err==nil && opErr==nil && f==Iterate; next,err,cont=i(f) {
            f,tmp,opErr=op(j,next,f);
            if f==Continue || f==Break || opErr!=nil {
                return tmp,opErr,(opErr==nil && f==Continue);
            } 
            j++;
        }
        return tmp,err,cont;
    }
}
func (i Iter[T])Next(
        op func(index int, val T, status IteratorFeedback) (IteratorFeedback,T,error)) Iter[T] {
    return Next(i,op);
}
