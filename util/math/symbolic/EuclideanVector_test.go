package symbolic

import (
	"reflect"
	"testing"

	customerr "github.com/barbell-math/block/util/err"
	"github.com/barbell-math/block/util/test"
)

func TestEuclideanVectorAddScalar(t *testing.T) {
    v:=EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }).Add(Scalar[int]{v: 1});
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !IsInvalidOperation(val) {
                test.FormatError(InvalidOperation(""),val,
                    "EuclideanVector+Scalar did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector+Scalar returned the wrong type.",t,
        );
    }
}

func TestEuclideanVectorAddEuclideanVector(t *testing.T) {
    v:=EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }).Add(EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }));
    switch v.(type) {
        case EuclideanVector[int]:
            val:=v.(EuclideanVector[int]);
            for i,v:=range(val) {
                test.BasicTest(int(2*i-2),v.(Scalar[int]).v,
                    "EuclideanVector+EuclideanVector did not procuce correct value.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector+EuclideanVector returned the wrong type.",t,
        );
    }
}

func TestEuclideanVectorAddEuclideanVectorWrongDimension(t *testing.T) {
    v:=EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }).Add(EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0},
    }));
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !customerr.IsDimensionsDoNotAgree(val) {
                test.FormatError(InvalidOperation(""),val,
                    "EuclideanVector+EuclideanVector with wrong dimensions did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector+EuclideanVector returned the wrong type.",t,
        );
    }
}

func TestEuclideanVectorAddPolarVector(t *testing.T) {
    v:=EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }).Add(PolarVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0},
    }));
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !IsInvalidOperation(val) {
                test.FormatError(InvalidOperation(""),val,
                    "EuclideanVector+PolarVector did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector+PolarVector returned the wrong type.",t,
        );
    }
}

func TestEuclideanVectorSubScalar(t *testing.T) {
    v:=EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }).Sub(Scalar[int]{v: 1});
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !IsInvalidOperation(val) {
                test.FormatError(InvalidOperation(""),val,
                    "EuclideanVector-Scalar did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector-Scalar returned the wrong type.",t,
        );
    }
}

func TestEuclideanVectorSubEuclideanVector(t *testing.T) {
    v:=EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }).Sub(EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: 1}, Scalar[int]{v: 0}, Scalar[int]{v: -1},
    }));
    switch v.(type) {
        case EuclideanVector[int]:
            val:=v.(EuclideanVector[int]);
            for i,v:=range(val) {
                test.BasicTest(int(2*i-2),v.(Scalar[int]).v,
                    "EuclideanVector-EuclideanVector did not procuce correct value.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector-EuclideanVector returned the wrong type.",t,
        );
    }
}

func TestEuclideanVectorSubEuclideanVectorWrongDimension(t *testing.T) {
    v:=EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }).Sub(EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0},
    }));
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !customerr.IsDimensionsDoNotAgree(val) {
                test.FormatError(InvalidOperation(""),val,
                    "EuclideanVector-EuclideanVector with wrong dimensions did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector-EuclideanVector returned the wrong type.",t,
        );
    }
}

func TestEuclideanVectorSubPolarVector(t *testing.T) {
    v:=EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }).Sub(PolarVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0},
    }));
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !IsInvalidOperation(val) {
                test.FormatError(InvalidOperation(""),val,
                    "EuclideanVector-PolarVector did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector-PolarVector returned the wrong type.",t,
        );
    }
}

func TestEuclideanVectorMulScalar(t *testing.T) {
    v:=EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }).Mul(Scalar[int]{v: 2});
    switch v.(type) {
        case EuclideanVector[int]:
            val:=v.(EuclideanVector[int]);
            for i,v:=range(val) {
                test.BasicTest(int(2*i-2),v.(Scalar[int]).v,
                    "EuclideanVector*Scalar did not procuce correct value.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector*Scalar returned the wrong type.",t,
        );
    }
}

func TestEuclideanVectorMulEuclideanVector(t *testing.T) {
    v:=EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }).Mul(EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: 1}, Scalar[int]{v: 0}, Scalar[int]{v: -1},
    }));
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !IsInvalidOperation(val) {
                test.FormatError(InvalidOperation(""),val,
                    "EuclideanVector*EuclideanVector did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector*EuclideanVector returned the wrong type.",t,
        );
    }
}

func TestEuclideanVectorMulPolarVector(t *testing.T) {
    v:=EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }).Mul(PolarVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0},
    }));
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !IsInvalidOperation(val) {
                test.FormatError(InvalidOperation(""),val,
                    "EuclideanVector*PolarVector did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector*PolarVector returned the wrong type.",t,
        );
    }
}

func TestEuclideanVectorDivScalar(t *testing.T) {
    v:=EuclideanVector[float32]([]Symbol[float32]{
        Scalar[float32]{v: -1}, Scalar[float32]{v: 0}, Scalar[float32]{v: 1},
    }).Mul(Scalar[float32]{v: 2});
    switch v.(type) {
        case EuclideanVector[float32]:
            val:=v.(EuclideanVector[float32]);
            for i,v:=range(val) {
                test.BasicTest(float32(2.0*float32(i)-2.0),v.(Scalar[float32]).v,
                    "EuclideanVector/Scalar did not procuce correct value.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector/Scalar returned the wrong type.",t,
        );
    }
}

func TestEuclideanVectorDivEuclideanVector(t *testing.T) {
    v:=EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }).Mul(EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: 1}, Scalar[int]{v: 0}, Scalar[int]{v: -1},
    }));
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !IsInvalidOperation(val) {
                test.FormatError(InvalidOperation(""),val,
                    "EuclideanVector/EuclideanVector did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector/EuclideanVector returned the wrong type.",t,
        );
    }
}

func TestEuclideanVectorDivPolarVector(t *testing.T) {
    v:=EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }).Mul(PolarVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0},
    }));
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !IsInvalidOperation(val) {
                test.FormatError(InvalidOperation(""),val,
                    "EuclideanVector/PolarVector did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "EuclideanVector/PolarVector returned the wrong type.",t,
        );
    }
}
