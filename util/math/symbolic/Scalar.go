package symbolic

import (
	"github.com/barbell-math/engine/util/math"
)

type Scalar[N math.Number] struct { v N };

func (s Scalar[N])Add(other Symbol[N]) Symbol[N] {
    switch any(other).(type) {
        case Scalar[N]: return Scalar[N]{v: s.v+other.(Scalar[N]).v};
        case SymbolicError[N]: return other;
        default: return invalidBinaryOpFormater[N](s,other,"+");
    }
}

func (s Scalar[N])Sub(other Symbol[N]) Symbol[N] {
    switch any(other).(type) {
        case Scalar[N]: return Scalar[N]{v: s.v-other.(Scalar[N]).v};
        case SymbolicError[N]: return other;
        default: return invalidBinaryOpFormater[N](s,other,"-");
    }
}

func (s Scalar[N])Mul(other Symbol[N]) Symbol[N] {
    switch any(other).(type) {
        case Scalar[N]: return Scalar[N]{v: s.v*other.(Scalar[N]).v};
        case EuclideanVector[N]: return s.mulEuclideanVector(other.(EuclideanVector[N]));
        case PolarVector[N]: return s.mulPolarVector(other.(PolarVector[N]));
        case SymbolicError[N]: return other;
        default: return invalidBinaryOpFormater[N](s,other,"*");
    }
}

func (s Scalar[N])mulEuclideanVector(v EuclideanVector[N]) Symbol[N] {
    rv:=NewVector[N,EuclideanVector[N]](len(v));
    for i,val:=range(v) {
        rv[i]=val.Mul(s);
    }
    return rv;
}

func (s Scalar[N])mulPolarVector(v PolarVector[N]) Symbol[N] {
    rv:=CopyVector(v);
    rv[0]=rv[0].Mul(s);
    return rv;
}

func (s Scalar[N])Div(other Symbol[N]) Symbol[N] {
    switch any(other).(type) {
        case Scalar[N]: return s.divByScalar(other.(Scalar[N]).v);
        case EuclideanVector[N]: return s.divEuclideanVector(other.(EuclideanVector[N]));
        case PolarVector[N]: return s.divPolarVector(other.(PolarVector[N]));
        case SymbolicError[N]: return other;
        default: return invalidBinaryOpFormater[N](s,other,"/");
    }
}

func (s Scalar[N])divByScalar(val N) Symbol[N] {
    if val==N(0) {
        return SymbolicError[N]{math.DivByZero("")};
    } else {
        return Scalar[N]{v: s.v/val};
    }
}

func (s Scalar[N])divEuclideanVector(v EuclideanVector[N]) Symbol[N] {
    rv:=NewVector[N,EuclideanVector[N]](len(v));
    for i,val:=range(v) {
        rv[i]=val.Div(s);
    }
    return rv;
}

func (s Scalar[N])divPolarVector(v PolarVector[N]) Symbol[N] {
    rv:=CopyVector(v);
    rv[0]=rv[0].Div(s);
    return rv;
}
