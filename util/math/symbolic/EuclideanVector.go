package symbolic

import (
	"fmt"

	customerr "github.com/barbell-math/block/util/err"
	"github.com/barbell-math/block/util/math"
)

type EuclideanVector[N math.Number] []Symbol[N];

func NewVector[N math.Number](l int) EuclideanVector[N] {
    return EuclideanVector[N](make([]Symbol[N],l));
}

func (v EuclideanVector[N])Add(other Symbol[N]) Symbol[N] {
    switch any(other).(type) {
        case Scalar[N]: return v.addScalar(other.(Scalar[N]));
        case EuclideanVector[N]: return v.addVector(other.(EuclideanVector[N]));
        case SymbolicError[N]: return other;
        default: return invalidBinaryOpFormater[N](v,other,"+");
    }
}

func (v EuclideanVector[N])addScalar(s Scalar[N]) Symbol[N] {
     rv:=NewVector[N](len(v));
     for i,val:=range(v) {
        rv[i]=val.Add(s);
     }
     return rv;
}

func (v EuclideanVector[N])addVector(otherV EuclideanVector[N]) Symbol[N] {
     if len(v)!=len(otherV) {
         return SymbolicError[N]{customerr.DimensionsDoNotAgree(fmt.Sprintf(
            "len(v1)=%d len(v2)=%d",len(v),len(otherV),
         ))};
     }
     rv:=NewVector[N](len(v));
     for i,val:=range(v) {
        rv[i]=val.Add(otherV[i]);
     }
     return rv;
}

func (v EuclideanVector[N])Sub(other Symbol[N]) Symbol[N] {
    switch any(other).(type) {
        case Scalar[N]: return v.subScalar(other.(Scalar[N]));
        case EuclideanVector[N]: return v.subVector(other.(EuclideanVector[N]));
        case SymbolicError[N]: return other;
        default: return invalidBinaryOpFormater[N](v,other,"+");
    }
}

func (v EuclideanVector[N])subScalar(s Scalar[N]) Symbol[N] {
     rv:=NewVector[N](len(v));
     for i,val:=range(v) {
        rv[i]=val.Sub(s);
     }
     return rv;
}

func (v EuclideanVector[N])subVector(otherV EuclideanVector[N]) Symbol[N] {
     if len(v)!=len(otherV) {
         return SymbolicError[N]{customerr.DimensionsDoNotAgree(fmt.Sprintf(
            "len(v1)=%d len(v2)=%d",len(v),len(otherV),
         ))};
     }
     rv:=NewVector[N](len(v));
     for i,val:=range(v) {
        rv[i]=val.Sub(otherV[i]);
     }
     return rv;
}

func (v EuclideanVector[N])Mul(other Symbol[N]) Symbol[N] {
    switch any(other).(type) {
        case Scalar[N]: return v.mulScalar(other.(Scalar[N]));
        case EuclideanVector[N]: return invalidBinaryOpFormater[N](v,other,"*");
        case SymbolicError[N]: return other;
        default: return invalidBinaryOpFormater[N](v,other,"*");
    }
}

func (v EuclideanVector[N])mulScalar(s Scalar[N]) Symbol[N] {
     rv:=NewVector[N](len(v));
     for i,val:=range(v) {
        rv[i]=val.Mul(s);
     }
     return rv;
}
