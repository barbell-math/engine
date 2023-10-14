package numeric

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/barbell-math/engine/util/algo/iter"
	"github.com/barbell-math/engine/util/io/csv"
	"github.com/barbell-math/engine/util/test"
	"github.com/barbell-math/engine/util/math"
)

func TestCreateLinReg(t *testing.T){
    l:=NewLinearReg[float32](ConstSumOpGen[float32](5.0)(
        []string{"x1","x2","x3","x4","x5"},"y"),
    );
    test.BasicTest(5,l.a.Rows(),"A matrix has wrong number of rows.",t);
    test.BasicTest(5,l.a.Cols(),"A matrix has wrong number of cols.",t);
    test.BasicTest(5,l.b.Rows(),"B matrix has wrong number of rows.",t);
    test.BasicTest(1,l.b.Cols(),"B matrix has wrong number of cols.",t);
    test.BasicTest(5,len(l.summationOps),
        "SummationOps matrix has wrong number of rows.",t,
    );
    test.BasicTest(6,len(l.summationOps[0]),
        "SummationOps matrix has wrong number of cols.",t,
    );
    test.BasicTest(5,len(l.iVarOps),"iVarOps has wrong number of elements.",t);
    l=NewLinearReg[float32](ConstSumOpGen[float32](5.0)([]string{"x1"},"y"));
    test.BasicTest(1,l.a.Rows(),"A matrix has wrong number of rows.",t);
    test.BasicTest(1,l.a.Cols(),"A matrix has wrong number of cols.",t);
    test.BasicTest(1,l.b.Rows(),"B matrix has wrong number of rows.",t);
    test.BasicTest(1,l.b.Cols(),"B matrix has wrong number of cols.",t);
    test.BasicTest(1,len(l.summationOps),
        "SummationOps matrix has wrong number of rows.",t,
    );
    test.BasicTest(2,len(l.summationOps[0]),
        "SummationOps matrix has wrong number of cols.",t,
    );
    test.BasicTest(1,len(l.iVarOps),"iVarOps has wrong number of elements.",t);
}

func TestConstantLinearReg(t *testing.T){
    l:=NewLinearReg[float32](ConstSumOpGen[float32](5.0)(
        []string{"x1","x2","x3"},"y"),
    );
    for i:=0; i<5; i++ {
        l.IterRHS(func(r int, c int, v float32){
            test.BasicTest(float32(25*i),v,
                "RHS matrix not updated by const summation op correctly.",t);
        });
        l.IterLHS(func(r int, c int, v float32){
            test.BasicTest(float32(25*i),v,
                "LHS matrix not updated by const summation op correctly.",t);
        });
        test.BasicTest(nil,l.UpdateSummations(map[string]float32{}),
            "Summation returned error when it was not supposed to.",t,
        );
    }
    _,_,err:=l.Run();
    if !math.IsSingularMatrix(err) {
        test.FormatError(math.SingularMatrix(""),err,
            "Result should have been singular.",t,
        );
    }
}

func Test1DStdLinearReg(t *testing.T){
    l:=NewLinearReg[float32](LinearSumOpGen[float32]([]string{"x1"},"y"));
    //Create and use data points (0,0) (1,1) (2,2) ... (10,10)
    for i:=0; i<11; i++ {
        l.UpdateSummations(map[string]float32{"x1": float32(i),"y": float32(i)});
    }
    test.BasicTest(float32(385),l.a.V[0][0],
        "A Summation was not run properly.",t,
    );
    test.BasicTest(float32(385),l.b.V[0][0],
        "B Summation was not run properly.",t,
    );
    res,rcond,err:=l.Run();
    test.BasicTest(float64(1),rcond,
        "Appropriate linear relationship was not found.",t,
    );
    test.BasicTest(nil,err,
        "Linear reg returned error when it shouldn't have.",t,
    );
    for i:=-12; i<14; i+=2 {
        v,err:=res.Predict(map[string]float32{"x1": float32(i)});
        test.BasicTest(nil,err,
            "Appropriate linear relationship was not found.",t,
        );
        test.BasicTest(float32(i),v,
            "Appropriate linear relationship was not found.",t,
        );
    }
}

func Test2DLinearReg(t *testing.T){
    l:=NewLinearReg[float64](LinearSumOpGen[float64]([]string{"x1","x2"},"y"));
    //Create and use data points (0,0) (1,1) (2,2) ... (10,10)
    for i:=0; i<11; i++ {
        l.UpdateSummations(map[string]float64{
            "x1": float64(i),
            "x2": float64(i/2.0),   //This is int div, makes non-singular matrix
            "y": float64(i),
        });
    }
    test.BasicTest(float64(385),l.a.V[0][0],
        "A Summation was not run properly.",t,
    );
    test.BasicTest(float64(85),l.a.V[1][1],
        "A Summation was not run properly.",t,
    );
    test.BasicTest(float64(180),l.a.V[1][0],
        "A Summation was not run properly.",t,
    );
    test.BasicTest(float64(180),l.a.V[0][1],
        "A Summation was not run properly.",t,
    );
    test.BasicTest(float64(385),l.b.V[0][0],
        "B Summation was not run properly.",t,
    );
    test.BasicTest(float64(180),l.b.V[1][0],
        "B Summation was not run properly.",t,
    );
    res,_,err:=l.Run();
    test.BasicTest(nil,err,
        "Linear reg returned error when it shouldn't have.",t,
    );
    for i:=-12; i<14; i+=2 {
        v,err:=res.Predict(map[string]float64{"x1": float64(i)});
        if !math.IsMissingVariable(err) {
            test.FormatError(math.MissingVariable(""),err,
                "Missing variable not caught.",t,
            );
        }
        v,err=res.Predict(map[string]float64{
            "x1": float64(i), "x2": float64(i/2.0),
        });
        test.BasicTest(nil,err,
            "Appropriate linear relationship was not found.",t,
        );
        if Abs(float64(i)-v)>WORKING_PRECISION {
            test.FormatError(float64(i),v,
                "Value is not within working precision of expected value.",t,
            );
        }
    }
}

func Test2DLinearRegWithError(t *testing.T){
    l:=NewLinearReg[float64](LinearSumOpGenWithError[float64](
        []string{"x1","x2"},"y"),
    );
    //Create and use data points (0,0) (1,1) (2,2) ... (10,10)
    for i:=0; i<11; i++ {
        l.UpdateSummations(map[string]float64{
            "x1": float64(i),
            "x2": float64(i/2.0),   //This is int div, makes non-singular matrix
            "y": float64(i),
        });
    }
    correctA:=[][]float64{
        []float64{385, 180, 55},
        []float64{180, 85, 25},
        []float64{55, 25, 11},
    };
    correctB:=[][]float64{
        []float64{385},
        []float64{180},
        []float64{55},
    };
    l.IterLHS(func(r int, c int, v float64){
         test.BasicTest(correctA[r][c],l.a.V[r][c],
             "A summation was not run properly.",t,
         );
    });
    l.IterRHS(func(r int, c int, v float64){
         test.BasicTest(correctB[r][c],l.b.V[r][c],
             "A summation was not run properly.",t,
         );
    });
    res,_,err:=l.Run();
    test.BasicTest(nil,err,
        "Linear reg returned error when it shouldn't have.",t,
    );
    for i:=-12; i<14; i+=2 {
        v,err:=res.Predict(map[string]float64{
            "x1": float64(i), "x2": float64(i/2.0),
        });
        test.BasicTest(nil,err,
            "Appropriate linear relationship was not found.",t,
        );
        if Abs(float64(i)-v)>WORKING_PRECISION*100 {
            test.FormatError(float64(i),v,
                "Value is not within working precision of expected value.",t,
            );
        }
    }
}

func TestNonStdLinearReg(t *testing.T){
    l:=NewLinearReg[float64]([]SummationOp[float64]{
        func(vals Vars[float64]) (float64,error) {
            v1,_:=vals.Access("x1");
            v2,_:=vals.Access("x2");
            return v1*v2,nil;
        },
        func (vals Vars[float64]) (float64,error) {
            v1,_:=vals.Access("x2");
            return v1*v1,nil;
        },
    },LinearSummationOp[float64]("y"));
    for i:=0; i<11; i++ {
        l.UpdateSummations(map[string]float64{
            "x1": float64(i),
            "x2": float64(i/2.0),
            "y": float64(i),
        });
    }
    test.BasicTest(float64(5762),l.a.V[0][0],
        "A Summation was not run properly.",t,
    );
    test.BasicTest(float64(2766),l.a.V[0][1],
        "A Summation was not run properly.",t,
    );
    test.BasicTest(float64(2766),l.a.V[1][0],
        "A Summation was not run properly.",t,
    );
    test.BasicTest(float64(1333),l.a.V[1][1],
        "A Summation was not run properly.",t,
    );
    test.BasicTest(float64(1430),l.b.V[0][0],
        "A Summation was not run properly.",t,
    );
    test.BasicTest(float64(680),l.b.V[1][0],
        "A Summation was not run properly.",t,
    );
    res,_,err:=l.Run();
    test.BasicTest(nil,err,
        "Linear reg returned error when it shouldn't have.",t,
    );
    csv.CSVFileSplitter("testData/NonStdLinRegActual.csv",',','#').ForEach(
    func(index int, val []string) (iter.IteratorFeedback, error) {
        cntr,_:=strconv.Atoi(val[0]);
        v,err:=res.Predict(map[string]float64{
            "x1": float64(cntr),"x2": float64(cntr/2.0),
        });
        actual,_:=strconv.ParseFloat(val[1],64);
        test.BasicTest(nil,err,
            "Appropriate linear relationship was not found.",t,
        );
        if Abs(v-actual)>WORKING_PRECISION*10 {
            test.FormatError(actual,v,
                "Value is not within working precision of expected value.",t,
            );
        }
        return iter.Continue,nil;

    });
}

//func TestLinearityOfBenchmarkResults(t *testing.T){
//    l:=NewLinearReg[float64](LinearSumOpGen[float64]([]string{"x1"},"y"));
//    util.CSVFileSplitter("testData/LinRegBenchmarkResults.csv",',',false,
//        func(cols []string) bool {
//            x1,_:=strconv.ParseFloat(cols[0],64);
//            y,_:=strconv.ParseFloat(cols[0],64);
//            l.UpdateSummations(map[string]float64{"x1": x1, "y": y});
//            return true;
//    });
//    _,rcond,_:=l.Run();
//    fmt.Println("rcond: ",rcond);
//}

func benchmarkStdLinReg(n int, nPnts int, b *testing.B){
    iVars:=make([]string,n);
    for i:=0; i<n; i++ {
        iVars[i]=fmt.Sprintf("x%d",i);
    }
    l:=NewLinearReg[float64](LinearSumOpGen[float64](iVars,"y"));
    pnts:=make(map[string]float64,n);
    for i:=0; i<b.N; i++ {
        for j:=0; j<nPnts; j++ {
            for k:=0; k<n; k++ {
                pnts[fmt.Sprintf("x%d",k)]=float64(i/(j+1));
            }
            l.UpdateSummations(pnts);
        }
        l.Run();
    }
}

func BenchmarkStdLinReg10_100(b *testing.B){ benchmarkStdLinReg(10,100,b); }
func BenchmarkStdLinReg10_1000(b *testing.B){ benchmarkStdLinReg(10,1000,b); }
func BenchmarkStdLinReg10_10000(b *testing.B){ benchmarkStdLinReg(10,10000,b); }
func BenchmarkStdLinReg10_100000(b *testing.B){ benchmarkStdLinReg(10,100000,b); }
func BenchmarkStdLinReg10_1000000(b *testing.B){ benchmarkStdLinReg(10,1000000,b); }
func BenchmarkStdLinReg100_100(b *testing.B){ benchmarkStdLinReg(100,100,b); }
func BenchmarkStdLinReg1000_1000(b *testing.B){ benchmarkStdLinReg(1000,1000,b); }
