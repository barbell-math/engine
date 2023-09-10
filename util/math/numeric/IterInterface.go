package numeric;

import (
    "fmt"
    "github.com/barbell-math/engine/util/math"
)

func Add[N math.Number](accum *N, iter N) error { 
    *accum=*accum+iter;
    return nil;
}

func Sub[N math.Number](accum *N, iter N) error { 
    *accum=*accum-iter;
    return nil;
}

func Mul[N math.Number](accum *N, iter N) error { 
    *accum=*accum*iter;
    return nil;
}

func Div[N math.Number](accum *N, iter N) error { 
    if iter==N(0) {
        return math.DivByZero(fmt.Sprintf("%v/%v",*accum,iter));
    }
    *accum=*accum/iter;
    return nil;
}
