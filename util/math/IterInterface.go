package math;

import (
    "fmt"
)

func ReduceAdd[N Number](accum *N, iter N) error { 
    *accum=*accum+iter;
    return nil;
}

func ReduceSub[N Number](accum *N, iter N) error { 
    *accum=*accum-iter;
    return nil;
}

func ReduceMul[N Number](accum *N, iter N) error { 
    *accum=*accum*iter;
    return nil;
}

func ReduceDiv[N Number](accum *N, iter N) error { 
    if iter==N(0) {
        return DivByZero(fmt.Sprintf("%v/%v",*accum,iter));
    }
    *accum=*accum/iter;
    return nil;
}
