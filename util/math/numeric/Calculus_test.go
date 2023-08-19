package numeric

import (
	"fmt"
	"testing"

	"github.com/barbell-math/block/util/test"
)

func TestDerivativeHorizontalLine(t *testing.T) {
    f:=Derivative(func(x float32) float32 { return 1; },0.01);
    test.BasicTest(float32(0),f(0),
        "Derivative of y=1 did not return 0.",t,
    );
    test.BasicTest(float32(0),f(-1),
        "Derivative of y=1 did not return 0.",t,
    );
    test.BasicTest(float32(0),f(1),
        "Derivative of y=1 did not return 0.",t,
    );
}

func TestDerivativeLine(t *testing.T) {
    f:=Derivative(func(x float32) float32 { return x; },0.01);
    test.BasicTest(true,Abs(float32(1)-f(0))<1e-6,
        "Derivative of y=x did not return value near 1.",t,
    );
    test.BasicTest(true,Abs(float32(1)-f(1))<1e-6,
        "Derivative of y=x did not return value near 1.",t,
    );
    test.BasicTest(true,Abs(float32(1)-f(-1))<1e-6,
        "Derivative of y=x did not return value near 1.",t,
    );
}

func TestDerivativeQuadratic(t *testing.T) {
    //for i:=0; i<30; i++ {
    //    v:=Derivative(func(x float32) float32 { return x*x; },0.1/float32(i))(1.0);
    //    fmt.Printf("h: %.10f, Val: %.10f, Err: %.10f\n",0.1/float32(i),v,Abs(2-v));
    //}
    f:=Derivative(func(x float32) float32 { return x*x; },0.1);
    test.BasicTest(float32(0),f(0),
        "Derivative of y=x^2 did not return 0 when evaluated at x=0.",t,
    );
    test.BasicTest(true,Abs(float32(1)-f(0.5))<1e-6,
        "Derivative of y=x^2 did not return 1 when evaluated at x=0.5.",t,
    );
    test.BasicTest(true,Abs(float32(-1)-f(-0.5))<1e-6,
        "Derivative of y=x^2 did not return -1 when evaluated at x=-0.5.",t,
    );
    test.BasicTest(true,Abs(float32(2)-f(1))<1e-6,
        "Derivative of y=x^2 did not return 1 when evaluated at x=0.5.",t,
    );
    test.BasicTest(true,Abs(float32(-2)-f(-1))<1e-6,
        "Derivative of y=x^2 did not return -1 when evaluated at x=-0.5.",t,
    );
}

func TestIntegralHorizontalLine(t *testing.T){
    f:=Integral(func(x float32) float32 {return 1;}) 
    v,_:=f(0,1,3);
    test.BasicTest(float32(1),v,
        "Integral of y=1 over [0,1] with step num 3 did not return 1.",t,
    );
    v,_=f(0,1,5);
    test.BasicTest(float32(1),v,
        "Integral of y=1 over [0,1] with step num 5 did not return 1.",t,
    );
    v,_=f(0,1,11);
    test.BasicTest(float32(1),v,
        "Integral of y=1 over [0,1] with step num 11 did not return 1.",t,
    );
}

func TestIntegralDiagonalLine(t *testing.T){
    f:=Integral(func(x float32) float32 {return x;}) 
    v,err:=f(0,2,3);
    fmt.Println(v,err);
    test.BasicTest(float32(2),v,
        "Integral of y=x over [0,2] with step num 3 did not return 2.",t,
    );
    v,_=f(0,2,5);
    test.BasicTest(float32(2),v,
        "Integral of y=x over [0,2] with step num 5 did not return 2.",t,
    );
    v,_=f(0,2,11);
    test.BasicTest(float32(2),v,
        "Integral of y=x over [0,2] with step num 11 did not return 2.",t,
    );
}

func TestIntegralQuadratic(t *testing.T){
    f:=Integral(func(x float32) float32 {return x*x;}) 
    v,err:=f(0,3,3);
    fmt.Println(v,err);
    test.BasicTest(float32(9),v,
        "Integral of y=x^2 over [0,3] with step num 3 did not return 9.",t,
    );
    v,_=f(0,3,5);
    test.BasicTest(float32(9),v,
        "Integral of y=x^2 over [0,3] with step num 5 did not return 9.",t,
    );
    v,_=f(0,3,11);
    test.BasicTest(float32(9),v,
        "Integral of y=x^2 over [0,3] with step num 11 did not return 9.",t,
    );
}

func TestDoubleIntegralHorizontalPlane(t *testing.T){
    f:=DoubleIntegral(func(x1, x2 float32) float32 { return 1;});
    val,_:=f(0,1,0,1,3);
    test.BasicTest(float32(1),val,
        "Double integral of z=1 over x[0,1] y[0,1] with step num 3 did not return 1.",t,
    );
    val,_=f(0,1,0,1,5);
    test.BasicTest(float32(1),val,
        "Double integral of z=1 over x[0,1] y[0,1] with step num 5 did not return 1.",t,
    );
    val,_=f(0,1,0,1,11);
    test.BasicTest(float32(1),val,
        "Double integral of z=1 over x[0,1] y[0,1] with step num 11 did not return 1.",t,
    );
}

func TestDoubleIntegralDiagonalPlane(t *testing.T){
    f:=DoubleIntegral(func(x1, x2 float32) float32 { return x1+x2;});
    val,_:=f(0,2,0,2,3);
    test.BasicTest(float32(8),val,
        "Double integral of z=x+y over x[0,2] y[0,2] with step num 3 did not return 8.",t,
    );
    val,_=f(0,2,0,2,5);
    test.BasicTest(float32(8),val,
        "Double integral of z=x+y over x[0,2] y[0,2] with step num 5 did not return 8.",t,
    );
    val,_=f(0,2,0,2,11);
    test.BasicTest(true,Abs(float32(8)-val)<1e-6,
        "Double integral of z=x+y over x[0,2] y[0,2] with step num 11 did not return 8.",t,
    );
}
