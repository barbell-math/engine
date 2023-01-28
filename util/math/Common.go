package math;

import (
    "fmt"
)

const WORKING_PRECISION float64=1e-16;

type Int interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
};

type Uint interface {
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
};

type Float interface {
    ~float32 | ~float64
};

type Number interface {
    Int |
    Uint |
    Float
};

func Max[N Number](vals ...N) N {
    rv:=vals[0];
    for _,v:=range(vals) {
        if v>rv {
            rv=v;
        }
    }
    return rv;
}

func Min[N Number](vals ...N) N {
    rv:=vals[0];
    for _,v:=range(vals) {
        if v<rv {
            rv=v;
        }
    }
    return rv;
}

func Abs[N Number](v N) N {
    if v<0 {
        return -v;
    }
    return v;
}

func SqErr[N Number](act N, given N) N {
    return (act-given)*(act-given);
}
func MeanSqErr[N Number](act []N, given []N) (N,error) {
    var rv N=N(0);
    if la,lg:=len(act),len(given); la!=lg {
        return rv,DimensionsDoNotAgree(fmt.Sprintf(
            "MSE requires lists of equal length. | len(Given)=%d len(act)=%d",
            la,lg,
        ));
    }
    if len(act)==0 {
        return N(0),nil;
    }
    for i,actIter:=range(act) {
        rv+=(actIter-given[i])*(actIter-given[i]);
    }
    return rv/N(len(act)),nil;
}

func Constrain[N Number](given N, min N, max N) N {
    if given<min {
        return min;
    } else if given>max {
        return max;
    }
    return given;
}
