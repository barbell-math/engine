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

//Why is join not a pseudo intermediary? It breaks the order that next calls
// the parrent consumer. Join not only needs to call two parrent consumers, but
// it also needs to manage caching of values between calls.
//func (i Iter[T])Join(i2 Iter[T], decider func(left T, right T) bool) Iter[T] {
//    var tmp T;
//    var i1Val, i2Val T;
//    var err1, err2 error;
//    cont1, cont2:=true, true;
//    getI1Val, getI2Val:=true, true;
//    return func(f IteratorFeedback) (T, error, bool) {
//        if f==Break {
//            return tmp, customerr.AppendError(i.Stop(),i2.Stop()), false;
//        }
//        if getI1Val && cont1 {
//            i1Val,err1,cont1=i(f);
//        }
//        if getI2Val && cont2 {
//            i2Val,err2,cont2=i2(f);
//        }
//        if err1==nil && err2==nil {
//            if cont1 && cont2 {
//                d:=decider(i1Val,i2Val);
//                getI1Val=d;
//                getI2Val=!d;
//                if d {
//                    return i1Val,err1,cont1 && cont2;
//                } else {
//                    return i2Val,err2,cont1 && cont2;
//                }
//            } else if cont1 && !cont2 {
//                getI1Val=true;
//                getI2Val=false;
//                return i1Val,err1,cont1;
//            } else if !cont1 && cont2 {
//                getI1Val=false;
//                getI2Val=true;
//                return i2Val,err2,cont2;
//            }
//        }
//        return tmp,customerr.AppendError(err1,err2),false;
//    }
//}
