package iter;

import (
	"github.com/barbell-math/block/util/dataStruct/types"
)

func (i Iter[T])Take(num int) Iter[T] {
    cntr:=0;
    return i.Next(
    func(index int, val T, status IteratorFeedback) (IteratorFeedback, T, error) {
        if status!=Break && cntr<num {
            cntr++;
            return Continue,val,nil;
        }
        return Break,val,nil;
    });
}

func (i Iter[T])TakeWhile(op func(val T) bool) Iter[T] {
    return i.Next(
    func(index int, val T, status IteratorFeedback) (IteratorFeedback, T, error) {
        if status!=Break && op(val) {
            return Continue,val,nil;
        }
        return Break,val,nil;
    });
}

func (i Iter[T])Skip(num int) Iter[T] {
    cntr:=-1;
    return i.Next(
    func(index int, val T, status IteratorFeedback) (IteratorFeedback, T, error) {
        if status==Break {
            return Break,val,nil;
        }
        cntr++;
        if cntr<num {
            return Iterate,val,nil;
        } else {
            return Continue,val,nil;
        }
    });
}

func Map[T any, U any](i Iter[T],
        op func(index int, val T) (U,error)) Iter[U] {
    return Next(i,
    func(index int, val T, status IteratorFeedback) (IteratorFeedback, U, error) {
        if status==Break {
            var tmp U;
            return Break,tmp,nil;
        }
        tmp,err:=op(index,val);
        return Continue,tmp,err;
    });
}
func (i Iter[T])Map(op func(index int, val T) (T,error)) Iter[T] {
    return Map(i,op);
}

func (i Iter[T])Filter(op func(val T) bool) Iter[T] {
    return i.Next(
    func(index int, val T, status IteratorFeedback) (IteratorFeedback, T, error) {
        if status!=Break && op(val) {
            return Continue,val,nil;
        }
        return Iterate,val,nil;
    });
}

func Window[T any](i Iter[T],
        q types.Queue[T],
        allowPartials bool,
        op func(index int, q types.Queue[T]) (IteratorFeedback,error)) Iter[types.Queue[T]] {
    return Next(i,
    func(index int, val T, status IteratorFeedback) (IteratorFeedback, types.Queue[T], error) {
        if status==Break {
            return Break,q,nil;
        }
        if q.Length()==q.Capacity() {
            q.Pop(); //Ignoring potential queue empty error because of if stmt
        }
        //Ignoring potential queue full error because of if stmt
        q.Push(val);
        if !allowPartials && q.Length()!=q.Capacity() {
            return Iterate,q,nil;
        }
        return Continue,q,nil;
    });
}
