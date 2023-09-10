package symbolic

import (
	"reflect"
	"testing"

	"github.com/barbell-math/block/util/math"
	"github.com/barbell-math/block/util/test"
)

func TestScalarAddScalar(t *testing.T) {
    s:=Scalar[int]{v: 1}.Add(Scalar[int]{v: 1});
    switch s.(type) {
        case Scalar[int]: 
            test.BasicTest(int(2),s.(Scalar[int]).v,
                "Scalar+Scalar produced incorrect value.",t,
            );
            s=Scalar[int]{v: 3}.Add(Scalar[int]{v: 5});
            test.BasicTest(int(8),s.(Scalar[int]).v,
                "Scalar+Scalar produced incorrect value.",t,
            );
        default: test.FormatError("Scalar[int]",reflect.TypeOf(s),
            "Scalar+scalar returned the wrong type.",t,
        );
    }
}

func TestScalarAddEuclideanVector(t *testing.T) {
    v:=Scalar[int]{v: 1}.Add(EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }));
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !IsInvalidOperation(val) {
                test.FormatError(InvalidOperation(""),val,
                    "Scalar+EuclideanVector did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "Scalar+EuclideanVector returned the wrong type.",t,
        );
    }
}

func TestScalarAddPolarVector(t *testing.T) {
    v:=Scalar[int]{v: 1}.Add(PolarVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }));
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !IsInvalidOperation(val) {
                test.FormatError(InvalidOperation(""),val,
                    "Scalar+PolarVector did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "Scalar+PolarVector returned the wrong type.",t,
        );
    }
}

func TestScalarSubScalar(t *testing.T) {
    s:=Scalar[int]{v: 1}.Sub(Scalar[int]{v: 1});
    switch s.(type) {
        case Scalar[int]: 
            test.BasicTest(int(0),s.(Scalar[int]).v,
                "Scalar-Scalar produced incorrect value.",t,
            );
            s=Scalar[int]{v: 3}.Sub(Scalar[int]{v: 5});
            test.BasicTest(int(-2),s.(Scalar[int]).v,
                "Scalar-Scalar produced incorrect value.",t,
            );
        default: test.FormatError("Scalar[int]",reflect.TypeOf(s),
            "Scalar-Scalar returned the wrong type.",t,
        );
    }
}

func TestScalarSubEuclideanVector(t *testing.T) {
    v:=Scalar[int]{v: 1}.Sub(EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: 1}, Scalar[int]{v: 2}, Scalar[int]{v: 3},
    }));
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !IsInvalidOperation(val) {
                test.FormatError(InvalidOperation(""),val,
                    "Scalar-EuclideanVector did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "Scalar-EuclideanVector returned the wrong type.",t,
        );
    }
}

func TestScalarSubPolarVector(t *testing.T) {
    v:=Scalar[int]{v: 1}.Sub(PolarVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }));
    switch v.(type) {
        case SymbolicError[int]:
            val:=v.(SymbolicError[int]);
            if !IsInvalidOperation(val) {
                test.FormatError(InvalidOperation(""),val,
                    "Scalar-PolarVector did not result in an error.",t,
                );
            }
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "Scalar-PolarVector returned the wrong type.",t,
        );
    }
}

func TestScalarMulScalar(t *testing.T) {
    s:=Scalar[int]{v: 2}.Mul(Scalar[int]{v: 1});
    switch s.(type) {
        case Scalar[int]: 
            test.BasicTest(int(2),s.(Scalar[int]).v,
                "Scalar*Scalar produced incorrect value.",t,
            );
            s=Scalar[int]{v: 3}.Mul(Scalar[int]{v: 5});
            test.BasicTest(int(15),s.(Scalar[int]).v,
                "Scalar*Scalar produced incorrect value.",t,
            );
        default: test.FormatError("Scalar[int]",reflect.TypeOf(s),
            "Scalar*Scalar returned the wrong type.",t,
        );
    }
}

func TestScalarMulEuclideanVector(t *testing.T) {
    v:=Scalar[int]{v: 2}.Mul(EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: 0}, Scalar[int]{v: 1}, Scalar[int]{v: 2},
    }));
    switch v.(type) {
        case EuclideanVector[int]:
            val:=v.(EuclideanVector[int]);
            for i,s:=range(val) {
                test.BasicTest(i*2,s.(Scalar[int]).v,
                    "Scalar*EuclideanVector returned the wrong type.",t,
                );
            }
        default: test.FormatError("EuclideanVector[int]",reflect.TypeOf(v),
            "Scalar*EuclideanVector returned the wrong type.",t,
        );
    }
}

func TestScalarMulPolarVector(t *testing.T) {
    v:=Scalar[int]{v: 2}.Mul(PolarVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }));
    switch v.(type) {
        case PolarVector[int]:
            val:=v.(PolarVector[int]);
            test.BasicTest(int(-2),val[0].(Scalar[int]).v,
                "Scalar*PolarVector did not change radius.",t,
            );
            test.BasicTest(int(0),val[1].(Scalar[int]).v,
                "Scalar*PolarVector changed theta_1.",t,
            );
            test.BasicTest(int(1),val[2].(Scalar[int]).v,
                "Scalar*PolarVector changed theta_2.",t,
            );
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "Scalar*PolarVector returned the wrong type.",t,
        );

    }
}

func TestScalarDivScalar(t *testing.T) {
    s:=Scalar[float32]{v: 2}.Div(Scalar[float32]{v: 4});
    switch s.(type) {
        case Scalar[float32]: 
            test.BasicTest(float32(2.0/4),s.(Scalar[float32]).v,
                "Scalar/Scalar produced incorrect value.",t,
            );
            s=Scalar[float32]{v: 3}.Div(Scalar[float32]{v: 5});
            test.BasicTest(float32(3.0/5),s.(Scalar[float32]).v,
                "Scalar/Scalar produced incorrect value.",t,
            );
            s=Scalar[float32]{v: 0}.Div(Scalar[float32]{v: 3});
            test.BasicTest(float32(0),s.(Scalar[float32]).v,
                "Scalar/Scalar produced incorrect value.",t,
            );
        default: test.FormatError("Scalar[int]",reflect.TypeOf(s),
            "Scalar/Scalar returned the wrong type.",t,
        );
    }
}

func TestScalarDivZero(t *testing.T) {
    s:=Scalar[float32]{v: 2}.Div(Scalar[float32]{v: 0});
    switch s.(type) {
        case SymbolicError[float32]: 
            val:=s.(SymbolicError[float32]);
            if !math.IsDivByZero(val) {
                test.FormatError(math.DivByZero(""),val,
                    "Scalar/0 did not result in an error.",t,
                );
            }
        default: test.FormatError("Scalar[int]",reflect.TypeOf(s),
            "Scalar/0 returned the wrong type.",t,
        );
    }
}

func TestScalarDivEuclideanVector(t *testing.T) {
    v:=Scalar[float32]{v: 2}.Div(EuclideanVector[float32]([]Symbol[float32]{
        Scalar[float32]{v: 0}, Scalar[float32]{v: 1}, Scalar[float32]{v: 2},
    }));
    switch v.(type) {
        case EuclideanVector[float32]:
            val:=v.(EuclideanVector[float32]);
            for i,s:=range(val) {
                test.BasicTest(float32(float32(i)/2.0),s.(Scalar[float32]).v,
                    "Scalar/EuclideanVector returned the wrong type.",t,
                );
            }
        default: test.FormatError("EuclideanVector[int]",reflect.TypeOf(v),
            "Scalar/EuclideanVector returned the wrong type.",t,
        );
    }
}

func TestScalarDivPolarVector(t *testing.T) {
    v:=Scalar[float32]{v: 2}.Div(PolarVector[float32]([]Symbol[float32]{
        Scalar[float32]{v: -1}, Scalar[float32]{v: 0}, Scalar[float32]{v: 1},
    }));
    switch v.(type) {
        case PolarVector[float32]:
            val:=v.(PolarVector[float32]);
            test.BasicTest(float32(-1.0/2.0),val[0].(Scalar[float32]).v,
                "Scalar/PolarVector did not change radius.",t,
            );
            test.BasicTest(float32(0),val[1].(Scalar[float32]).v,
                "Scalar/PolarVector dhanged theta_1.",t,
            );
            test.BasicTest(float32(1.0),val[2].(Scalar[float32]).v,
                "Scalar/PolarVector dhanged theta_2.",t,
            );
        default: test.FormatError("SymbolicError[int]",reflect.TypeOf(v),
            "Scalar/PolarVector returned the wrong type.",t,
        );
    }
}


func BenchmarkAddScalarToScalar(b *testing.B) {
    for i:=0; i<b.N; i++ {
        s:=Scalar[int]{v: 1}.Add(Scalar[int]{v: 1}).(Scalar[int]);
        s.v++;
    }
}

func BenchmarkAddScalarToEuclideanVector(b *testing.B) {
    for i:=0; i<b.N; i++ {
        v:=Scalar[int]{v: 1}.Add(EuclideanVector[int]([]Symbol[int]{
            Scalar[int]{v: 1}, Scalar[int]{v: 1},
            Scalar[int]{v: 1}, Scalar[int]{v: 1},
        })).(EuclideanVector[int]);
        v[0]=v[0].Add(Scalar[int]{v: 1});
    }
}

func BenchmarkAddScalarToScalarNormal(b *testing.B) {
    for i:=0; i<b.N; i++ {
        s:=1+1;
        s++;
    }
}

func BenchmarkAddScalarToVectorNormal(b *testing.B) {
    for i:=0; i<b.N; i++ {
        v:=[]int{1,1,1,1};
        for i,_:=range(v) {
            v[i]+=1;
        }
        v[0]++;
    }
}
