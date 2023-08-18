package symbolic;

import (
    "testing"
    "reflect"
	"github.com/barbell-math/block/util/test"
)

func TestScalarAddScalar(t *testing.T) {
    s:=Scalar[int]{v: 1}.Add(Scalar[int]{v: 1});
    switch s.(type) {
        case Scalar[int]: 
            test.BasicTest(int(2),s.(Scalar[int]).v,
                "Adding two scalars produced incorrect value.",t,
            );
            s=Scalar[int]{v: 3}.Add(Scalar[int]{v: 5});
            test.BasicTest(int(8),s.(Scalar[int]).v,
                "Adding two scalars produced incorrect value.",t,
            );
        default: test.FormatError("Scalar[int]",reflect.TypeOf(s),
            "Adding two scalars returned the wrong type.",t,
        );
    }
}

func TestScalarAddVector(t *testing.T) {
    v:=Scalar[int]{v: 1}.Add(EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: -1}, Scalar[int]{v: 0}, Scalar[int]{v: 1},
    }));
    switch v.(type) {
        case EuclideanVector[int]:
            val:=v.(EuclideanVector[int]);
            for i,s:=range(val) {
                test.BasicTest(i,s.(Scalar[int]).v,
                    "Adding a vector to a scalar returned the wrong type.",t,
                );
            }
        default: test.FormatError("Vector[int]",reflect.TypeOf(v),
            "Adding a vector to a scalar returned the wrong type.",t,
        );

    }
}

func TestScalarSubScalar(t *testing.T) {
    s:=Scalar[int]{v: 1}.Sub(Scalar[int]{v: 1});
    switch s.(type) {
        case Scalar[int]: 
            test.BasicTest(int(0),s.(Scalar[int]).v,
                "Adding two scalars produced incorrect value.",t,
            );
            s=Scalar[int]{v: 3}.Sub(Scalar[int]{v: 5});
            test.BasicTest(int(-2),s.(Scalar[int]).v,
                "Adding two scalars produced incorrect value.",t,
            );
        default: test.FormatError("Scalar[int]",reflect.TypeOf(s),
            "Adding two scalars returned the wrong type.",t,
        );
    }
}

func TestScalarSubVector(t *testing.T) {
    v:=Scalar[int]{v: 1}.Sub(EuclideanVector[int]([]Symbol[int]{
        Scalar[int]{v: 1}, Scalar[int]{v: 2}, Scalar[int]{v: 3},
    }));
    switch v.(type) {
        case EuclideanVector[int]:
            val:=v.(EuclideanVector[int]);
            for i,s:=range(val) {
                test.BasicTest(i,s.(Scalar[int]).v,
                    "Adding a vector to a scalar returned the wrong type.",t,
                );
            }
        default: test.FormatError("Vector[int]",reflect.TypeOf(v),
            "Adding a vector to a scalar returned the wrong type.",t,
        );

    }
}

func BenchmarkAddScalarToScalar(b *testing.B) {
    for i:=0; i<b.N; i++ {
        s:=Scalar[int]{v: 1}.Add(Scalar[int]{v: 1}).(Scalar[int]);
        s.v++;
    }
}

func BenchmarkAddScalarToVector(b *testing.B) {
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
