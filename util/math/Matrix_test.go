package math;

import (
    "testing"
    "github.com/barbell-math/block/util/test"
)

func TestCreateMatrix(t *testing.T){
    m:=NewMatrix[int](3,4,ZeroFill[int]);
    test.BasicTest(3,len(m.V),"Rows are not correct.",t);
    test.BasicTest(4,len(m.V[0]),"Columns are not correct.",t);
    m.Iter(func(r int, c int, v int){
        test.BasicTest(0, v,"Value is not correct.",t);
    });
}

func TestRowsCols(t *testing.T){
    m:=NewMatrix[int](3,4,ZeroFill[int]);
    test.BasicTest(3,m.Rows(),"Rows are not correct.",t);
    test.BasicTest(4,m.Cols(),"Columns are not correct.",t);
    m=NewMatrix[int](1,0, ZeroFill[int]);
    test.BasicTest(0 ,m.Rows(),"Rows are not correct.",t);
    test.BasicTest(0 ,m.Cols(),"Columns are not correct.",t);
    m=NewMatrix[int](0 ,1, ZeroFill[int]);
    test.BasicTest(0 ,m.Rows(),"Rows are not correct.",t);
    test.BasicTest(0 ,m.Cols(),"Columns are not correct.",t);
    m=NewMatrix[int](0 ,0, ZeroFill[int]);
    test.BasicTest(0 ,m.Rows(),"Rows are not correct.",t);
    test.BasicTest(0 ,m.Cols(),"Columns are not correct.",t);
}

func TestFill(t *testing.T){
    m:=NewMatrix[int](4,3,IdentityFill[int]);
    test.BasicTest(4,len(m.V),"Rows are not correct.",t);
    test.BasicTest(3,len(m.V[0]),"Columns are not correct.",t);
    m.Iter(func(r int, c int, v int){
        if r==c {
            test.BasicTest(1, v,"Value is not correct.",t);
        } else {
            test.BasicTest(0, v,"Value is not correct.",t);
        }
    });
    m=NewMatrix[int](4,3,func(r int, c int) int {
        if c==1 {
            return 1;
        }
        return 0;
    });
    test.BasicTest(4,len(m.V),"Rows are not correct.",t);
    test.BasicTest(3,len(m.V[0]),"Columns are not correct.",t);
    m.Iter(func(r int, c int, v int){
        if c==1 {
            test.BasicTest(1, v,"Value is not correct.",t);
        } else {
            test.BasicTest(0, v,"Value is not correct.",t);
        }
    });
    m1:=NewMatrix[int](4,3,DuplicateFill(&m));
    test.BasicTest(4,len(m1.V),"Rows are not correct.",t);
    test.BasicTest(3,len(m1.V[0]),"Columns are not correct.",t);
    m1.Iter(func(r int, c int, v int){
        if c==1 {
            test.BasicTest(1, v,"Value is not correct.",t);
        } else {
            test.BasicTest(0, v,"Value is not correct.",t);
        }
    });
    m=NewMatrix[int](4,3,ConstFill(5));
    test.BasicTest(4,len(m.V),"Rows are not correct.",t);
    test.BasicTest(3,len(m.V[0]),"Columns are not correct.",t);
    m.Iter(func(r int, c int, v int){
        test.BasicTest(5, v,"Value is not correct.",t);
    });
    m=NewMatrix[int](4,1,ArrayFill([][]int{{1},{2},{3},{4}}));
    test.BasicTest(4,len(m.V),"Rows are not correct.",t);
    test.BasicTest(1,len(m.V[0]),"Columns are not correct.",t);
    m.Iter(func(r int, c int, v int){
        test.BasicTest(r+1,v,"Value is not correct.",t);
    });
}

func TestAdd(t *testing.T){
    m1:=NewMatrix[int](3,4,func(r int, c int) int {
        return c-1;
    });
    m2:=NewMatrix[int](3,4,func(r int, c int) int {
        return 1;
    });
    m3:=NewMatrix[int](3,3,func(r int, c int) int {
        return 1;
    });
    err:=m1.Add(&m3);
    if !IsMatrixDimensionsDoNotAgree(err) {
        test.FormatError(
            MatrixDimensionsDoNotAgree(""),err,
            "Matrix dimensions were not successfully checked.",t,
        );
    }
    err=m1.Add(&m2);
    test.BasicTest(nil,err,"An error occurred during matrix addition.",t);
    m1.Iter(func(r int, c int, v int){
        test.BasicTest(c,v,"Addition was not executed properly.",t);
    });
}

func TestAddScalar(t *testing.T){
    m1:=NewMatrix[int](3,4,func(r int, c int) int {
        return c-1;
    });
    m1.AddScalar(1);
    m1.Iter(func(r int, c int, v int){
        test.BasicTest(c,v,"Addition was not executed properly.",t);
    });
}

func TestSub(t *testing.T){
    m1:=NewMatrix[int](3,4,func(r int, c int) int {
        return c+1;
    });
    m2:=NewMatrix[int](3,4,func(r int, c int) int {
        return 1;
    });
    m3:=NewMatrix[int](3,3,func(r int, c int) int {
        return 1;
    });
    err:=m1.Sub(&m3);
    if !IsMatrixDimensionsDoNotAgree(err) {
        test.FormatError(
            MatrixDimensionsDoNotAgree(""),err,
            "Matrix dimensions were not successfully checked.",t,
        );
    }
    err=m1.Sub(&m2);
    test.BasicTest(nil,err,"An error occurred during matrix addition.",t);
    m1.Iter(func(r int, c int, v int){
        test.BasicTest(c,v,"Addition was not executed properly.",t);
    });
}

func TestSubScalar(t *testing.T){
    m1:=NewMatrix[int](3,4,func(r int, c int) int {
        return c+1;
    });
    m1.SubScalar(1);
    m1.Iter(func(r int, c int, v int){
        test.BasicTest(c,v,"Addition was not executed properly.",t);
    });
}

func TestEquals(t *testing.T){
    m1:=NewMatrix[int](3,4,func(r int, c int) int {
        return c-1;
    });
    m2:=NewMatrix[int](4,4,func(r int, c int) int {
        return c-1;
    });
    res,err:=m1.Equals(&m2,0);
    if !IsMatrixDimensionsDoNotAgree(err) {
        test.FormatError(
            MatrixDimensionsDoNotAgree(""),err,
            "Matrix dimensions were not successfully checked.",t,
        );
    }
    test.BasicTest(false,res,"Equal returned false positive.",t);
    m2=NewMatrix[int](3,4,func(r int, c int) int {
        return c-1;
    });
    res,err=m1.Equals(&m2,0);
    test.BasicTest(true,res,"Equal returned false negative.",t);
    test.BasicTest(nil,err,"Equal returned error it wasn't supposed to.",t);
    m2.AddScalar(1);
    res,err=m1.Equals(&m2,0);
    test.BasicTest(false,res,"Equal returned false positive.",t);
    test.BasicTest(nil,err,"Equal returned error it wasn't supposed to.",t);
    res,err=m1.Equals(&m2,1);
    test.BasicTest(true,res,"Equal returned false negative.",t);
    test.BasicTest(nil,err,"Equal returned error it wasn't supposed to.",t);
}

func TestMul(t *testing.T){
    m1:=NewMatrix[int](3,4,func(r int, c int) int {
        return c+1;
    });
    m2:=NewMatrix[int](4,3,func(r int, c int) int {
        return r+1;
    });
    m3:=NewMatrix[int](3,4,func(r int, c int) int {
        return 1;
    });
    err:=m1.Mul(&m3);
    if !IsMatrixDimensionsDoNotAgree(err) {
        test.FormatError(
            MatrixDimensionsDoNotAgree(""),err,
            "Matrix dimensions were not successfully checked.",t,
        );
    }
    err=m1.Mul(&m2);
    test.BasicTest(nil,err,"An error occurred during matrix addition.",t);
    test.BasicTest(3,m1.Rows(),"Matrix is not proper size.",t);
    test.BasicTest(3,m1.Cols(),"Matrix is not proper size.",t);
    m1.Iter(func(r int, c int, v int){
        test.BasicTest(30,v,"Addition was not executed properly.",t);
    });
}

func TestMulScalar(t *testing.T){
    m1:=NewMatrix[int](3,4,func(r int, c int) int {
        return c+1;
    });
    m1.MulScalar(2);
    m1.Iter(func(r int, c int, v int){
        test.BasicTest((c+1)*2,v,"Addition was not executed properly.",t);
    });
}

func TestTranspose(t *testing.T){
    m1:=NewMatrix[int](2,4,func(r int, c int) int{
        return c+4*r;
    });
    m1.Transpose();
    test.BasicTest(4,m1.Rows(),"Matrix is not the proper size.",t);
    test.BasicTest(2,m1.Cols(),"Matrix is not the proper size.",t);
    m1.Iter(func(r int, c int, v int){
        test.BasicTest(r+4*c,v,"Transpose was not executed properly.",t);
    });
}

func TestInverse(t *testing.T){
    m1:=NewMatrix[float32](2,4,func(r int, c int) float32{
        return float32(c+4*r);
    });
    _,err:=m1.Inverse();
    if !IsInverseOfNonSquareMatrix(err) {
        test.FormatError(
            InverseOfNonSquareMatrix(""),err,
            "Inverse of non square matrix was taken.",t,
        );
    }
    m1=NewMatrix[float32](2,2,ZeroFill[float32]);
    _,err=m1.Inverse();
    if !IsSingularMatrix(err) {
        test.FormatError(
            SingularMatrix(""),err,
            "Inverse of singular matrix did not raise error.",t,
        );
    }
    m1=NewMatrix[float32](4,4,IdentityFill[float32]);
    _,err=m1.Inverse();
    test.BasicTest(nil,err,"Inverse returned error it wasn't supposed to.",t);
    m1.Iter(func(r int, c int, v float32){
        if r==c {
            test.BasicTest(
                float32(1.0),v,"Inverse was not calculated correctly.",t,
            );
        } else {
            test.BasicTest(
                float32(0.0),v,"Inverse was not calculated correctly.",t,
            );
        }
    });
    m2:=NewMatrix[float64](4,4,ZeroFill[float64]);
    m2.V=[][]float64{
        {1.0,4.0,1.0,6.0},
        {8.0,7.0,8.0,3.0},
        {0.0,4.0,6.0,3.0},
        {4.0,9.0,8.0,6.0},
    };
    correctMatrix:=NewMatrix[float64](4,4,ZeroFill[float64]);
    correctMatrix.V=[][]float64{
        {7.0/65.0, 27.0/130.0, -1.0/26.0, -5.0/26.0},
        {-4.0/13.0, -4.0/13.0, -6.0/13.0, 9.0/13.0},
        {2.0/65.0, 17.0/130.0, 9.0/26.0, -7.0/26.0},
        {68.0/195.0, 29.0/195.0, 10.0/39.0, -5.0/13.0},
    };
    _,err=m2.Inverse();
    //By default, MATLAB has 16 digits of precision, so I would say
    //this is acceptable
    res,err:=m2.Equals(&correctMatrix,1e-15);
    test.BasicTest(true,res,"Inverse was not calculated correctly.",t);
}

func TestRcond(t *testing.T){
    m1:=NewMatrix[float32](4,4,IdentityFill[float32]);
    rcond,err:=m1.Inverse();
    test.BasicTest(float64(1),rcond,"Rcond was not calculated correctly.",t);
    m3:=NewMatrix[float64](4,4,ZeroFill[float64]);
    m3.V=[][]float64{
        {1.0e+14, 0, 0, 0},
        {0, 7.0e-14, 0, 0},
        {0, 0, 6.0e+14, 0},
        {0, 0, 0, 6.0e-14},
    };
    rcond,err=m3.Inverse();
    if !IsMatrixSingularToWorkingPrecision(err) {
        test.FormatError(
            MatrixSingularToWorkingPrecision(""),err,
            "Rcond error was not created.",t,
        );
    }
}

func benchmarkDuplicateFill(r int, c int, b *testing.B){
    m1:=NewMatrix[int](r,c,IdentityFill[int]);
    for i:=0; i<b.N; i++ {
        NewMatrix[int](r,c,DuplicateFill[int](&m1));
    }
}
func benchmarkCopy(r int, c int, b *testing.B){
    m1:=NewMatrix[int](r,c,IdentityFill[int]);
    for i:=0; i<b.N; i++ {
        m1.Copy();
    }
}
func benchmarkAdd(r int, c int, b *testing.B){
    m1:=NewMatrix[int](r,c,ZeroFill[int]);
    m2:=NewMatrix[int](r,c,ConstFill(1));
    for i:=0; i<b.N; i++ {
        m1.addTraditional(&m2);
    }
}
func benchmarkAddRange(r int, c int, b *testing.B){
    m1:=NewMatrix[int](r,c,ZeroFill[int]);
    m2:=NewMatrix[int](r,c,ConstFill(1));
    for i:=0; i<b.N; i++ {
        m1.Add(&m2);
    }
}
func benchmarkAddIter(r int, c int, b *testing.B){
    m1:=NewMatrix[int](r,c,ZeroFill[int]);
    m2:=NewMatrix[int](r,c,ConstFill(1));
    for i:=0; i<b.N; i++ {
        m1.addIterCallback(&m2);
    }
}
func benchmarkAddFunctional(r int, c int, b *testing.B){
    m1:=NewMatrix[int](r,c,ZeroFill[int]);
    m2:=NewMatrix[int](r,c,ConstFill(1));
    for i:=0; i<b.N; i++ {
        m1.addFuncional(&m2);
    }
}
func benchmarkMul(r int, c int, b *testing.B){
    m1:=NewMatrix[int](r,c,ZeroFill[int]);
    m2:=NewMatrix[int](r,c,ConstFill(1));
    for i:=0; i<b.N; i++ {
        m1.Mul(&m2);
    }
}
func benchmarkInv(r int, c int, b *testing.B){
    m1:=NewMatrix[int](r,c,func(r int, c int) int {
        return c+r*4;
    });
    for i:=0; i<b.N; i++ {
        m1.Inverse();
    }
}

func BenchmarkDuplicateFill10(b *testing.B) { benchmarkDuplicateFill(10,10,b); }
func BenchmarkDuplicateFill100(b *testing.B) { benchmarkDuplicateFill(100,100,b); }
func BenchmarkCopy10(b *testing.B) { benchmarkCopy(10,10,b); }
func BenchmarkCopy100(b *testing.B) { benchmarkCopy(100,100,b); }
func BenchmarkAddTraditional1(b *testing.B) { benchmarkAdd(1,1,b); }
func BenchmarkAddTraditional10(b *testing.B) { benchmarkAdd(10,10,b); }
func BenchmarkAddTraditional100(b *testing.B) { benchmarkAdd(100,100,b); }
func BenchmarkAddRange1(b *testing.B) { benchmarkAddRange(1,1,b); }
func BenchmarkAddRange10(b *testing.B) { benchmarkAddRange(10,10,b); }
func BenchmarkAddRange100(b *testing.B) { benchmarkAddRange(100,100,b); }
func BenchmarkAddIter1(b *testing.B) { benchmarkAddIter(1,1,b); }
func BenchmarkAddIter10(b *testing.B) { benchmarkAddIter(10,10,b); }
func BenchmarkAddIter100(b *testing.B) { benchmarkAddIter(100,100,b); }
func BenchmarkAddFunctional1(b *testing.B) { benchmarkAddFunctional(1,1,b); }
func BenchmarkAddFunctional10(b *testing.B) { benchmarkAddFunctional(10,10,b); }
func BenchmarkAddFunctional100(b *testing.B) { benchmarkAddFunctional(100,100,b); }
func BenchmarkMul1(b *testing.B) { benchmarkMul(1,1,b); }
func BenchmarkMul10(b *testing.B) { benchmarkMul(10,10,b); }
func BenchmarkMul100(b *testing.B) { benchmarkMul(100,100,b); }
func BenchmarkInv10(b *testing.B) { benchmarkMul(10,10,b); }
func BenchmarkInv100(b *testing.B) { benchmarkMul(100,100,b); }
