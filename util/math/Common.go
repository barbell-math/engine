package math;

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

func Constrain[N Number](given N, min N, max N) N {
    if given<min {
        return min;
    } else if given>max {
        return max;
    }
    return given;
}
