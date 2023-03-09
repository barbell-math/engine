package algo;

import (
    "fmt"
    customerr "github.com/barbell-math/block/util/err"
)

type forEachParallelResult[T any, U any] struct {
    val T;
    res U;
    err error;
};
func (f forEachParallelResult[T,U])Unpack() (T,U,error){
    return f.val,f.res,f.err;
}

func forEachThreadRunner[T any, U any](
        jobs chan T,
        results chan forEachParallelResult[T,U],
        op func(val T) (U,error)){
    for val:=range(jobs) {
        tmp,err:=op(val);
        results <- forEachParallelResult[T, U]{
            val: val,
            res: tmp,
            err: err,
        };
    }
}

func NoOp[T any, U any](val T, res U, err error) { return; };

func ForEachParallel[T any, U any](
        i Iter[T],
        workerOp func(val T) (U,error),
        resOp func(val T, res U, err error),
        numThreads int) error {
    if err:=numThreadsCheck(numThreads); err!=nil {
        return err;
    }
    taken,j:=0,0;
    jobs, results:=make(chan T), make(chan forEachParallelResult[T,U]);
    for next,err,cont:=i(); cont && err==nil; next,err,cont=i() {
        if j<numThreads {
            //If another worker can be created, make one
            go forEachThreadRunner(jobs,results,workerOp);
        } else {
            //If all the worker threads are used wait for one to finish
            resOp((<-results).Unpack());
            taken++;
        }
        jobs <- next;
        j++;
    }
    close(jobs);
    //Need to wait for all the results!!
    for i:=0; i<j-taken; i++ { resOp((<-results).Unpack()); }
    close(results);
    return nil;
}
func (i Iter[T])ForEachParallel(
        workerOp func(val T) (T,error),
        resOp func(val T, res T, err error),
        numThreads int) error {
    return ForEachParallel(i,workerOp,resOp,numThreads);
}

func (i Iter[T])FilterParallel(op Filter[T], numThreads int) []T {
    rv:=make([]T,0);
    ForEachParallel(i,func(val T) (bool,error) {
        return op(val),nil;
    },func(val T, res bool, err error){
        if res {
            rv=append(rv,val);
        }
    },numThreads);
    return rv;
}

func numThreadsCheck(numThreads int) error {
    if numThreads<1 {
        return customerr.ValOutsideRange(fmt.Sprintf(
            "Expected >0 | Got: %d",numThreads,
        ));
    }
    return nil;
}
