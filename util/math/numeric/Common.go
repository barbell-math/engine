package numeric

import (
	stdMath "math"

	"github.com/barbell-math/engine/util/algo/iter"
	"github.com/barbell-math/engine/util/dataStruct"
	customerr "github.com/barbell-math/engine/util/err"
	"github.com/barbell-math/engine/util/math"
)

const WORKING_PRECISION float64=1e-16;

func Max[N math.Number](vals ...N) N {
    rv:=vals[0];
    for _,v:=range(vals) {
        if v>rv {
            rv=v;
        }
    }
    return rv;
}

func Min[N math.Number](vals ...N) N {
    rv:=vals[0];
    for _,v:=range(vals) {
        if v<rv {
            rv=v;
        }
    }
    return rv;
}

func Abs[N math.Number](v N) N {
    if v<0 {
        return -v;
    }
    return v;
}

func SqErr[N math.Number](act []N, given []N) ([]N,error) {
    if err:=customerr.ArrayDimsArgree(
        act,given,"MSE requires lists of equal length.",
    ); err!=nil || len(act)==0 {
        return []N{},err;
    }
    rv:=make([]N,len(act));
    for i,actIter:=range(act) {
        rv[i]=(actIter-given[i])*(actIter-given[i]);
    }
    return rv,nil;
}
func MeanSqErr[N math.Number](act []N, given []N) (N,error) {
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

func NoOpConstraint[N math.Number]() dataStruct.Pair[N,N] {
    return dataStruct.Pair[N, N]{A: N(stdMath.Inf(-1)), B: N(stdMath.Inf(1))};
}
func PositiveConstraint[N math.Number]() dataStruct.Pair[N,N] {
    return dataStruct.Pair[N, N]{A: N(0), B: N(stdMath.Inf(1))};
}
func NegativeConstraint[N math.Number]() dataStruct.Pair[N,N] {
    return dataStruct.Pair[N, N]{A: N(stdMath.Inf(-1)), B: N(0)};
}
func Constrain[N math.Number](given N, minMax dataStruct.Pair[N,N]) N {
    if given<minMax.A {
        return minMax.A;
    } else if given>minMax.B {
        return minMax.B;
    }
    return given;
}

func Range[N math.Number](start N, stop N, step N) iter.Iter[N] {
    cntr:=N(0);
    return func(f iter.IteratorFeedback) (N,error,bool) {
        iterVal:=step*cntr+start;
        if f!=iter.Break && ((step>0 && iterVal<stop) || (step<0 && iterVal>stop)) {
            cntr++;
            return iterVal,nil,true;
        }
        return stop,nil,false;
    }
}
