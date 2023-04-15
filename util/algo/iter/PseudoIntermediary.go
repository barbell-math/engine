package iter;


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

func (i Iter[T])Window(size int,
        op func(index int, vals []T) (IteratorFeedback,error)) error {
    if size<1 {
        return customerr.ValOutsideRange(fmt.Sprintf(
            "Window size must be >=1 | Have: %d",size,
        ));
    }

}
