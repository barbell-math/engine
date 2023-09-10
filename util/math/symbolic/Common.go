package symbolic

import (
	"github.com/barbell-math/engine/util/math"
)

type Symbol[N math.Number] interface {
    Add(other Symbol[N]) Symbol[N];
    Sub(other Symbol[N]) Symbol[N];
    Mul(other Symbol[N]) Symbol[N];
    Div(other Symbol[N]) Symbol[N];
    //Pow
    //Inverse
    //Derivative
    //Integral
    //ToEucidean
    //ToPolar
    //Cross
    //Dot
};
type Matrix[N math.Number] [][]Symbol[N];
type Equation[N math.Number] func(iVars Vars[N]) (Symbol[N],error);

type Vectors[N math.Number] interface {
    EuclideanVector[N] | PolarVector[N]
};

func NewVector[N math.Number, T Vectors[N]](l int) T {
    return T(make([]Symbol[N],l));
}

func CopyVector[N math.Number, T Vectors[N]](v T) T {
    rv:=NewVector[N,T](len(v));
    for i,v:=range(v) {
        rv[i]=v;
    }
    return rv;
}
