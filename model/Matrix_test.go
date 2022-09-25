package model;

import (
    //"fmt"
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
)

func TestCreateMatrix(t *testing.T){
    m:=NewMatrix[int](3,4,ZeroFill[int]);
    testUtil.BasicTest(3,len(m.V),"Rows are not correct.",t);
    testUtil.BasicTest(4,len(m.V[0]),"Columns are not correct.",t);
    m.Iter(func(r int, c int, v int){
        testUtil.BasicTest(0, v,"Value is not correct.",t);
    });
}

func TestRowsCols(t *testing.T){
    m:=NewMatrix[int](3,4,ZeroFill[int]);
    testUtil.BasicTest(3,m.Rows(),"Rows are not correct.",t);
    testUtil.BasicTest(4,m.Cols(),"Columns are not correct.",t);
    m=NewMatrix[int](1,0, ZeroFill[int]);
    testUtil.BasicTest(0 ,m.Rows(),"Rows are not correct.",t);
    testUtil.BasicTest(0 ,m.Cols(),"Columns are not correct.",t);
    m=NewMatrix[int](0 ,1, ZeroFill[int]);
    testUtil.BasicTest(0 ,m.Rows(),"Rows are not correct.",t);
    testUtil.BasicTest(0 ,m.Cols(),"Columns are not correct.",t);
    m=NewMatrix[int](0 ,0, ZeroFill[int]);
    testUtil.BasicTest(0 ,m.Rows(),"Rows are not correct.",t);
    testUtil.BasicTest(0 ,m.Cols(),"Columns are not correct.",t);
}

func TestFill(t *testing.T){
    m:=NewMatrix[int](4,3,IdentityFill[int]);
    testUtil.BasicTest(4,len(m.V),"Rows are not correct.",t);
    testUtil.BasicTest(3,len(m.V[0]),"Columns are not correct.",t);
    m.Iter(func(r int, c int, v int){
        if r==c {
            testUtil.BasicTest(1, v,"Value is not correct.",t);
        } else {
            testUtil.BasicTest(0, v,"Value is not correct.",t);
        }
    });
    m=NewMatrix[int](4,3,func(r int, c int) int {
        if c==1 {
            return 1;
        }
        return 0;
    });
    testUtil.BasicTest(4,len(m.V),"Rows are not correct.",t);
    testUtil.BasicTest(3,len(m.V[0]),"Columns are not correct.",t);
    m.Iter(func(r int, c int, v int){
        if c==1 {
            testUtil.BasicTest(1, v,"Value is not correct.",t);
        } else {
            testUtil.BasicTest(0, v,"Value is not correct.",t);
        }
    });
    m=NewMatrix[int](4,3,ConstFill(5));
    testUtil.BasicTest(4,len(m.V),"Rows are not correct.",t);
    testUtil.BasicTest(3,len(m.V[0]),"Columns are not correct.",t);
    m.Iter(func(r int, c int, v int){
        testUtil.BasicTest(5, v,"Value is not correct.",t);
    });
}

func TestAdd( t *testing.T){
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
    if !util.IsMatrixDimensionsDoNotAgree(err) {
        testUtil.FormatError(
            util.MatrixDimensionsDoNotAgree(""),err,
            "Matrix dimensions were not successfully checked.",t,
        );
    }
    err=m1.Add(&m2);
    testUtil.BasicTest(nil,err,"An error occurred during matrix addition.",t);
    m1.Iter(func(r int, c int, v int){
        testUtil.BasicTest(c,v,"Addition was not executed properly.",t);
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
    if !util.IsMatrixDimensionsDoNotAgree(err) {
        testUtil.FormatError(
            util.MatrixDimensionsDoNotAgree(""),err,
            "Matrix dimensions were not successfully checked.",t,
        );
    }
    err=m1.Sub(&m2);
    testUtil.BasicTest(nil,err,"An error occurred during matrix addition.",t);
    m1.Iter(func(r int, c int, v int){
        testUtil.BasicTest(c,v,"Addition was not executed properly.",t);
    });
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
    if !util.IsMatrixDimensionsDoNotAgree(err) {
        testUtil.FormatError(
            util.MatrixDimensionsDoNotAgree(""),err,
            "Matrix dimensions were not successfully checked.",t,
        );
    }
    err=m1.Mul(&m2);
    testUtil.BasicTest(nil,err,"An error occurred during matrix addition.",t);
    testUtil.BasicTest(3,m1.Rows(),"Matrix is not proper size.",t);
    testUtil.BasicTest(3,m1.Cols(),"Matrix is not proper size.",t);
    m1.Iter(func(r int, c int, v int){
        testUtil.BasicTest(30,v,"Addition was not executed properly.",t);
    });
}

func TestTranspose(t *testing.T){
    m1:=NewMatrix[int](2,4,func(r int, c int) int{
        return c+4*r;
    });
    m1.Transpose();
    testUtil.BasicTest(4,m1.Rows(),"Matrix is not the proper size.",t);
    testUtil.BasicTest(2,m1.Cols(),"Matrix is not the proper size.",t);
    m1.Iter(func(r int, c int, v int){
        testUtil.BasicTest(r+4*c,v,"Transpose was not executed properly.",t);
    });
}

func benchmarkAdd(r int, c int, b *testing.B){
    m1:=NewMatrix[int](10,10,ZeroFill[int]);
    m2:=NewMatrix[int](10,10,ConstFill(1));
    for i:=0; i<b.N; i++ {
        m1.addTraditional(&m2);
    }
}
func benchmarkAddRange(r int, c int, b *testing.B){
    m1:=NewMatrix[int](10,10,ZeroFill[int]);
    m2:=NewMatrix[int](10,10,ConstFill(1));
    for i:=0; i<b.N; i++ {
        m1.Add(&m2);
    }
}
func benchmarkAddIter(r int, c int, b *testing.B){
    m1:=NewMatrix[int](10,10,ZeroFill[int]);
    m2:=NewMatrix[int](10,10,ConstFill(1));
    for i:=0; i<b.N; i++ {
        m1.addIterCallback(&m2);
    }
}
func benchmarkAddFunctional(r int, c int, b *testing.B){
    m1:=NewMatrix[int](10,10,ZeroFill[int]);
    m2:=NewMatrix[int](10,10,ConstFill(1));
    for i:=0; i<b.N; i++ {
        m1.addFuncional(&m2);
    }
}
func benchmarkMul(r int, c int, b *testing.B){
    m1:=NewMatrix[int](10,10,ZeroFill[int]);
    m2:=NewMatrix[int](10,10,ConstFill(1));
    for i:=0; i<b.N; i++ {
        m1.Mul(&m2);
    }
}

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
