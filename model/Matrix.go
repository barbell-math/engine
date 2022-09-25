package model;

import (
    "fmt"
    "github.com/carmichaeljr/powerlifting-engine/util"
)

type Number interface {
    int | int32 | int64 | float32 | float64
};

type Matrix[N Number] struct {
    V [][]N;
};

type RowColOps[N Number] func(r int, c int) N;
func ZeroFill[N Number](r int, c int) N { return N(0); }
func IdentityFill[N Number](r int, c int) N {
    if r==c {
        return N(1);
    }
    return N(0);
}
func ConstFill[N Number](v N) RowColOps[N] {
    return func(r int, c int) N {
        return v;
    }
}

func NewMatrix[N Number](r int, c int, fill RowColOps[N]) Matrix[N] {
    var rv Matrix[N];
    rv.V=make([][]N,r);
    for i,_:=range(rv.V) {
        rv.V[i]=make([]N,c);
        for j,_:=range(rv.V[i]){
            rv.V[i][j]=fill(i,j);
        }
    }
    return rv;
}

func (m *Matrix[N])Rows() int {
    if len(m.V)>0 && len(m.V[0])==0 {
        return 0;
    }
    return len(m.V);
}
func (m *Matrix[N])Cols() int {
    if len(m.V)==0 {
        return 0;
    }
    return len(m.V[0]);
}

func (m *Matrix[N])Fill(fill RowColOps[N]){
    for i,_:=range(m.V) {
        for j,_:=range(m.V[i]){
            m.V[i][j]=fill(i,j);
        }
    }
}

func (m *Matrix[N])Iter(f func(r int, c int, v N)){
    for i,_:=range(m.V) {
        for j,_:=range(m.V[i]){
            f(i,j,m.V[i][j]);
        }
    }
}

func (m *Matrix[N])Add(other *Matrix[N]) error {
    return m.addSubOps(other,true);
}

func (m *Matrix[N])AddScalar(v N){
    for i,_:=range(m.V) {
        for j,_:=range(m.V[0]) {
            m.V[i][j]+=v;
        }
    }
}

func (m *Matrix[N])Sub(other *Matrix[N]) error {
    return m.addSubOps(other,false);
}

func (m *Matrix[N])SubScalar(v N){
    for i,_:=range(m.V) {
        for j,_:=range(m.V[0]) {
            m.V[i][j]+=v;
        }
    }
}

func (m *Matrix[N])addSubOps(other *Matrix[N], add bool) error {
    if len(other.V)!=len(m.V) || len(other.V[0])!=len(m.V[0]) {
        return util.MatrixDimensionsDoNotAgree(
            fmt.Sprintf("[r1=%d c1=%d] [r2=%d c2=%d] | Need r1=r2 and c1=c2",
                m.Rows(),m.Cols(),other.Rows(),other.Cols(),
        ));
    }
    for i,_:=range(m.V) {
        for j,_:=range(m.V[0]) {
            if add {
                m.V[i][j]+=other.V[i][j];
            } else {
                m.V[i][j]-=other.V[i][j];
            }
        }
    }
    return nil;
}

func (m *Matrix[N])Mul(other *Matrix[N]) error {
    if m.Cols()!=other.Rows() {
        return util.MatrixDimensionsDoNotAgree(
            fmt.Sprintf("[r1=%d c1=%d] [r2=%d c2=%d] | Need c1=r2",
                m.Rows(),m.Cols(),other.Rows(),other.Cols(),
        ));
    }
    tmp:=NewMatrix(m.Rows(),other.Cols(),ZeroFill[N]);
    for mRow:=0; mRow<m.Rows(); mRow++ {
        for oCol:=0; oCol<other.Cols(); oCol++ {
            var sum N=0;
            for i:=0; i<m.Cols(); i++ {
                sum+=m.V[mRow][i]*other.V[i][oCol];
            }
            tmp.V[mRow][oCol]=sum;
        }
    }
    *m=tmp;
    return nil;
}

func (m *Matrix[N])Transpose(){
    tmp:=NewMatrix(m.Cols(),m.Rows(),ZeroFill[N]);
    for i:=0; i<m.Rows(); i++ {
        for j:=0; j<m.Cols(); j++ {
            tmp.V[j][i]=m.V[i][j];
        }
    }
    *m=tmp;
}

func (m *Matrix[N])Inverse(){

}

//=============================================================================
//These functions should *not* be used, and are only here for
// benchmarking in the associated matrix test file.
//=============================================================================
func (m *Matrix[N])addFuncional(other *Matrix[N]) error {
    return util.ChainedErrorOps(
        func(r ...any) (any,error) {
            return nil,util.ErrorOnBool(
                len(other.V)==len(m.V) && len(other.V[0])==len(m.V[0]),
                util.MatrixDimensionsDoNotAgree(
                    fmt.Sprintf("Given: [r=%d c=%d] Other: [r=%d c=%d]",
                        m.Rows(),m.Cols(),other.Rows(),other.Cols(),
            )));
        }, func(r ...any) (any,error){
            other.Iter(func(r int, c int, v N){
                m.V[r][c]+=v;
            });
            return nil,nil;
    });
}
func (m *Matrix[N])addIterCallback(other *Matrix[N]) error {
    if len(other.V)!=len(m.V) || len(other.V[0])!=len(m.V[0]) {
        return util.MatrixDimensionsDoNotAgree(
            fmt.Sprintf("Given: [r=%d c=%d] Other: [r=%d c=%d]",
                m.Rows(),m.Cols(),other.Rows(),other.Cols(),
            ),
        );
    }
    other.Iter(func(r int, c int, v N){
        m.V[r][c]+=v;
    });
    return nil;
}
func (m *Matrix[N])addTraditional(other *Matrix[N]) error {
    if len(other.V)!=len(m.V) || len(other.V[0])!=len(m.V[0]) {
        return util.MatrixDimensionsDoNotAgree(
            fmt.Sprintf("Given: [r=%d c=%d] Other: [r=%d c=%d]",
                m.Rows(),m.Cols(),other.Rows(),other.Cols(),
            ),
        );
    }
    for i:=0; i<len(m.V); i++ {
        for j:=0; j<len(m.V[0]); j++ {
            m.V[i][j]+=other.V[i][j];
        }
    }
    return nil;
}
