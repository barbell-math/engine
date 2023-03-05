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

func SqErr[N Number](act []N, given []N) ([]N,error) {
    if err:=arrayDimsArgree(
        act,given,"MSE requires lists of equal length.",
    ); err!=nil {
        return []N{},err;
    }
    if len(act)==0 {
        return []N{},nil;
    }
    rv:=make([]N,len(act));
    for i,actIter:=range(act) {
        rv[i]=(actIter-given[i])*(actIter-given[i]);
    }
    return rv,nil;
}
func MeanSqErr[N Number](act []N, given []N) (N,error) {
    var sum N=N(0);
    sqErrs,err:=SqErr(act,given);
    if err!=nil || len(sqErrs)==0 {
        return sum,err;
    }
    for _,v:=range(sqErrs) {
        sum+=v;
    }
    return sum/N(len(act)),err;
}

func arrayDimsArgree[N any](act []N, given []N, message string) error {
    if la,lg:=len(act),len(given); la!=lg {
        return DimensionsDoNotAgree(fmt.Sprintf(
            "%s | len(Given)=%d len(act)=%d",
            message,la,lg,
        ));
    }
    return nil;
}

func Constrain[N Number](given N, min N, max N) N {
    if given<min {
        return min;
    } else if given>max {
        return max;
    }
    return given;
}
