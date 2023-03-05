package math;

import (
    "fmt"
)

func Add[N Number](accum *N, iter N) error { 
    *accum=*accum+iter;
    return nil;
}

func Sub[N Number](accum *N, iter N) error { 
    *accum=*accum-iter;
    return nil;
}

func Mul[N Number](accum *N, iter N) error { 
    *accum=*accum*iter;
    return nil;
}

func Div[N Number](accum *N, iter N) error { 
    if iter==N(0) {
        return DivByZero(fmt.Sprintf("%v/%v",*accum,iter));
    }
    *accum=*accum/iter;
    return nil;
}
