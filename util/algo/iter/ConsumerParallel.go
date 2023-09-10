package iter;

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

func forEachWorker[T any, U any](
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

func Parallel[T any, U any](i Iter[T],
        workerOp func(val T) (U,error),
        resOp func(val T, res U, err error),
        numThreads int) error {
    if err:=numThreadsCheck(numThreads); err!=nil {
        return err;
    }
    taken,j:=0,0;
    jobs, results:=make(chan T), make(chan forEachParallelResult[T,U]);
    i.ForEach(func(index int, val T) (IteratorFeedback, error) {
        if j<numThreads {
            //If another worker can be created, make one
            go forEachWorker(jobs,results,workerOp);
        } else {
            //If all the worker threads are used wait for one to finish
            resOp((<-results).Unpack());
            taken++;
        }
        jobs <- val;
        j++;
        return Continue,nil;
    });
    close(jobs);
    //Need to wait for all the results!!
    for i:=0; i<j-taken; i++ { resOp((<-results).Unpack()); }
    close(results);
    return nil;
}
func (i Iter[T])Parallel(
        workerOp func(val T) (T,error),
        resOp func(val T, res T, err error),
        numThreads int) error {
    return Parallel(i,workerOp,resOp,numThreads);
}

func numThreadsCheck(numThreads int) error {
    if numThreads<1 {
        return customerr.ValOutsideRange(fmt.Sprintf(
            "Expected >0 | Got: %d",numThreads,
        ));
    }
    return nil;
}
