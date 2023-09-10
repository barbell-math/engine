package dataStruct;

import (
	"github.com/barbell-math/block/util/dataStruct/types"
)

type VariantFlag int;
const (
    A VariantFlag=iota
    B
);

type Variant[T any, U any] struct {
    val any
    aOrB VariantFlag;
};

//The return type for these two functions has to be types.Variant because that 
// is what the interface expects. The interface cannot use a specific return 
//value because that would require the interface to import this module, creating 
//circular imports.
func (v Variant[T,U])SetValA(newVal T) types.Variant[T,U] {
    v.val=newVal;
    v.aOrB=A;
    return v;
}

func (v Variant[T,U])SetValB(newVal U) types.Variant[T,U] {
    v.val=newVal;
    v.aOrB=B;
    return v;
}

func (v Variant[T,U])HasA() bool { return v.aOrB==A };
func (v Variant[T,U])HasB() bool { return v.aOrB==B };

func (v Variant[T,U])ValA() T {
    return v.val.(T);
}

func (v Variant[T,U])ValB() U {
    return v.val.(U);
}

func (v Variant[T,U])ValAOr(_default T) T {
    if v.aOrB==A {
        return v.val.(T);
    }
    return _default;
}

func (v Variant[T,U])ValBOr(_default U) U {
    if v.aOrB==B {
        return v.val.(U);
    }
    return _default;
}
