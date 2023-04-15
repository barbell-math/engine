package iter;

func (i Iter[T])FilterParallel(op func(val T) bool, numThreads int) ([]T,error) {
    rv:=make([]T,0);
    err:=Parallel(i,func(val T) (bool,error) {
        return op(val),nil;
    },func(val T, res bool, err error){
        if res {
            rv=append(rv,val);
        }
    },numThreads);
    return rv,err;
}
