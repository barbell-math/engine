package symbolic

import (
	"github.com/barbell-math/block/util/math"
)

type Scalar[N math.Number] struct { v N };

func (s Scalar[N])Add(other Symbol[N]) Symbol[N] {
    switch any(other).(type) {
        case Scalar[N]: return Scalar[N]{v: s.v+other.(Scalar[N]).v};
        case EuclideanVector[N]: return s.addVector(other.(EuclideanVector[N]));
        case SymbolicError[N]: return other;
        default: return invalidBinaryOpFormater[N](s,other,"+");
    }
}

func (s Scalar[N])addVector(v EuclideanVector[N]) Symbol[N] {
    rv:=NewVector[N](len(v));
    for i,val:=range(v) {
        rv[i]=val.Add(s);
    }
    return rv;
}

func (s Scalar[N])Sub(other Symbol[N]) Symbol[N] {
    switch any(other).(type) {
        case Scalar[N]: return Scalar[N]{v: s.v-other.(Scalar[N]).v};
        case EuclideanVector[N]: return s.subVector(other.(EuclideanVector[N]));
        case SymbolicError[N]: return other;
        default: return invalidBinaryOpFormater[N](s,other,"-");
    }
}

func (s Scalar[N])subVector(v EuclideanVector[N]) Symbol[N] {
    rv:=NewVector[N](len(v));
    for i,val:=range(v) {
        rv[i]=val.Sub(s);
    }
    return rv;
}

func (s Scalar[N])Mul(other Symbol[N]) Symbol[N] {
    switch any(other).(type) {
        case Scalar[N]: return Scalar[N]{v: s.v*other.(Scalar[N]).v};
        case EuclideanVector[N]: return s.subVector(other.(EuclideanVector[N]));
        case SymbolicError[N]: return other;
        default: return invalidBinaryOpFormater[N](s,other,"-");
    }
}

func (s Scalar[N])mulVector(v EuclideanVector[N]) Symbol[N] {
    rv:=NewVector[N](len(v));
    for i,val:=range(v) {
        rv[i]=val.Mul(s);
    }
    return rv;
}
