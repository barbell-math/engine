package iter

import (
	"fmt"
	"math"
	"testing"
)

func sequenceGenerator(start float64, end float64, step float64) Iter[float64] {
    cntr:=start;
    return func(f IteratorFeedback) (float64,error,bool) {
        if f!=Break && cntr<=end {
            cntr+=step;
            return cntr-1,nil,true;
        }
        return end,nil,false;
    }
}

//Find the area under the function y=5cos(2pi/5(x-5))+5 between [-100,100]
// using a Riemann sum with a given step size
func testGeneratorIterator(step float64, print bool) func() {
    return func(){
        amp:=5.0;
        period:=5.0;
        hShift:=5.0;
        vShift:=5.0;
        val,err:=sequenceGenerator(-100.0,100.0,step).Map(func(index int, val float64) (float64,error) {
            height:=amp*math.Cos(2*math.Pi/period*(val-hShift))+vShift;
            return height*step,nil;
        }).Reduce(0.0, func(accum *float64, iter float64) error {
            *accum+=iter;
            return nil;
        });
        if print {
            fmt.Printf("Area is: %f Using step size: %f\n",val,step);
            fmt.Printf("Err is: %v\n",err);
        }
    }
}

func testGeneratorIterator2(step float64, print bool) func() {
    return func(){
        amp:=5.0;
        period:=5.0;
        hShift:=5.0;
        vShift:=5.0;
        total:=0.0;
        err:=sequenceGenerator(-100.0,100.0,step).ForEach(func(index int, val float64) (IteratorFeedback,error) {
            height:=amp*math.Cos(2*math.Pi/period*(val-hShift))+vShift;
            total+=height*step;
            return Continue,nil;
        });
        if print {
            fmt.Printf("Area is: %f Using step size: %f\n",total,step);
            fmt.Printf("Err is: %v\n",err);
        }
    }
}

func testGeneratorForLoop(step float64, print bool) func() {
    return func(){
        amp:=5.0;
        period:=5.0;
        hShift:=5.0;
        vShift:=5.0;
        total:=0.0;
        for x:=-100.0; x<=100; x+=step {
            height:=amp*math.Cos(2*math.Pi/period*(x-hShift))+vShift;
            total+=height*step;
        }
        if print {
            fmt.Printf("Area is: %f Using step size: %f\n",total,step);
        }
    }
}

func TestExample1(t *testing.T){
    testGeneratorIterator(1,true)();
}
func TestExample01(t *testing.T){
    testGeneratorIterator(0.1,true)();
}
func TestExample001(t *testing.T){
    testGeneratorIterator(0.01,true)();
}
func TestExample0001(t *testing.T){
    testGeneratorIterator(0.001,true)();
}
func TestExample00001(t *testing.T){
    testGeneratorIterator(0.0001,true)();
}

//Benchmarks
func BenchmarkExampleIterator1_1(b *testing.B){
    for i:=0; i<b.N; i++ {
        testGeneratorIterator(1,false)();
    }
}
func BenchmarkExampleIterator1_2(b *testing.B){
    for i:=0; i<b.N; i++ {
        testGeneratorIterator2(1,false)();
    }
}
func BenchmarkExampleForLoop1(b *testing.B) {
    for i:=0; i<b.N; i++ {
        testGeneratorForLoop(1,false)();
    }
}

func BenchmarkExampleIterator01(b *testing.B){
    for i:=0; i<b.N; i++ {
        testGeneratorIterator(0.1,false)();
    }
}
func BenchmarkExampleIterator01_2(b *testing.B){
    for i:=0; i<b.N; i++ {
        testGeneratorIterator2(0.1,false)();
    }
}
func BenchmarkExampleForLoop01(b *testing.B) {
    for i:=0; i<b.N; i++ {
        testGeneratorForLoop(0.1,false)();
    }
}

func BenchmarkExampleIterator001(b *testing.B){
    for i:=0; i<b.N; i++ {
        testGeneratorIterator(0.01,false)();
    }
}
func BenchmarkExampleIterator001_2(b *testing.B){
    for i:=0; i<b.N; i++ {
        testGeneratorIterator2(0.01,false)();
    }
}
func BenchmarkExampleForLoop001(b *testing.B) {
    for i:=0; i<b.N; i++ {
        testGeneratorForLoop(0.01,false)();
    }
}

func BenchmarkExampleIterator0001(b *testing.B){
    for i:=0; i<b.N; i++ {
        testGeneratorIterator(0.001,false)();
    }
}
func BenchmarkExampleIterator0001_2(b *testing.B){
    for i:=0; i<b.N; i++ {
        testGeneratorIterator2(0.001,false)();
    }
}
func BenchmarkExampleForLoop0001(b *testing.B) {
    for i:=0; i<b.N; i++ {
        testGeneratorForLoop(0.001,false)();
    }
}

func BenchmarkExampleIterator00001(b *testing.B){
    for i:=0; i<b.N; i++ {
        testGeneratorIterator(0.0001,false)();
    }
}
func BenchmarkExampleIterator00001_2(b *testing.B){
    for i:=0; i<b.N; i++ {
        testGeneratorIterator2(0.0001,false)();
    }
}
func BenchmarkExampleForLoop00001(b *testing.B) {
    for i:=0; i<b.N; i++ {
        testGeneratorForLoop(0.0001,false)();
    }
}
