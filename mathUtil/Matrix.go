package mathUtil;

import (
    "fmt"
    "github.com/carmichaeljr/powerlifting-engine/util"
)

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
func DuplicateFill[N Number](m *Matrix[N]) RowColOps[N] {
    return func(r int, c int) N {
        return m.V[r][c];
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

//Copy is more performant than duplicate fill, use copy when possible
func (m *Matrix[N])Copy() Matrix[N] {
    var rv Matrix[N];
    rv.V=make([][]N,m.Rows());
    for i,_:=range(rv.V) {
        rv.V[i]=make([]N,m.Cols());
        copy(rv.V[i],m.V[i]);
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

func (m *Matrix[N])Equals(other *Matrix[N], tol N) (bool,error) {
    var rv bool=true;
    if err:=matchingDimensionsErrorCheck(m,other); err!=nil {
        return false,err;
    }
    for i:=0; i<len(m.V) && rv; i++ {
        for j:=0; j<len(m.V[i]) && rv; j++ {
            rv=(Abs(m.V[i][j]-other.V[i][j])<=tol);
        }
    }
    return rv,nil;
}

func (m *Matrix[N])AddScalar(v N){
    for i,_:=range(m.V) {
        for j,_:=range(m.V[0]) {
            m.V[i][j]+=v;
        }
    }
}

func (m *Matrix[N])SubScalar(v N){
    for i,_:=range(m.V) {
        for j,_:=range(m.V[0]) {
            m.V[i][j]-=v;
        }
    }
}

func (m *Matrix[N])MulScalar(v N){
    for i,_:=range(m.V) {
        for j,_:=range(m.V[0]) {
            m.V[i][j]*=v;
        }
    }
}

func (m *Matrix[N])Add(other *Matrix[N]) error {
    return m.addSubOps(other,true);
}

func (m *Matrix[N])Sub(other *Matrix[N]) error {
    return m.addSubOps(other,false);
}

func (m *Matrix[N])addSubOps(other *Matrix[N], add bool) error {
    if err:=matchingDimensionsErrorCheck(m,other); err!=nil {
        return err;
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
    if err:=matchingInnerDimensionsErrorCheck(m,other); err!=nil {
        return err;
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

//www.mathworks.com/help/dsp/ref/reciprocalcondition.html
//Reciprocal condition number:
//  if return val is near 1.0 the matrix is well conditioned (low error)
//  if return val is near 0.0 the matrix is badly conditioned (high error)
func (m *Matrix[N])Inverse() (float64,error) {
    if err:=squareMatrixErrorCheck(m); err!=nil {
        return -1,err;
    }
    var rv error=nil;
    colMax:=Abs(m.getMaxColSum());
    tmp:=NewMatrix(m.Rows(),m.Cols(),IdentityFill[N]);
    for i:=0; i<m.Rows(); i++ {
        if m.V[i][i]==N(0) {
            rv=util.SingularMatrix(fmt.Sprintf("M[r=%d c=%d]=0",i,i));
        }
        tmp.divideRow(i,m.V[i][i]);
        m.divideRow(i,m.V[i][i]);
        for j:=0; j<m.Rows(); j++ {
            if j!=i {
                zeroVal(m,&tmp,i,j);
            }
        }
    }
    *m=tmp;
    invColMax:=Abs(m.getMaxColSum());
    rcond:=1.0/float64(colMax*invColMax);
    if rv==nil && rcond<WORKING_PRECISION {
        rv=util.MatrixSingularToWorkingPrecision(fmt.Sprintf("RCOND=%e",rcond));
    }
    return rcond,rv;
}

func (m *Matrix[N])getMaxColSum() N {
    var rv=make([]N,m.Cols());
    for _,col:=range(m.V) { //Iterate over rows first to avoid cache misses
        for c,v:=range(col) {
            rv[c]+=v;
        }
    }
    return Max(rv...);
}

func (m *Matrix[N])divideRow(r int, divVal N){
    for i,_:=range(m.V[r]) {
        m.V[r][i]/=divVal;
    }
}

func zeroVal[N Number](
        m1 *Matrix[N],
        m2 *Matrix[N],
        pivot int,
        changedRow int){
    mulVal:=-m1.V[changedRow][pivot];
    for i,_:=range(m1.V[changedRow]) {
        m1.V[changedRow][i]+=(m1.V[pivot][i]*mulVal);
        m2.V[changedRow][i]+=(m2.V[pivot][i]*mulVal);
    }
}

func matchingDimensionsErrorCheck[N Number](
        m *Matrix[N],
        other *Matrix[N]) error {
    if m.Rows()!=other.Rows() || m.Cols()!=other.Cols() {
        return util.MatrixDimensionsDoNotAgree(
            fmt.Sprintf("[r1=%d c1=%d] [r2=%d c2=%d] | Need r1=r2 and c1=c2",
                m.Rows(),m.Cols(),other.Rows(),other.Cols(),
        ));
    }
    return nil;
}

func matchingInnerDimensionsErrorCheck[N Number](
        m *Matrix[N],
        other *Matrix[N]) error {
    if m.Cols()!=other.Rows() {
        return util.MatrixDimensionsDoNotAgree(
            fmt.Sprintf("[r1=%d c1=%d] [r2=%d c2=%d] | Need c1=r2",
                m.Rows(),m.Cols(),other.Rows(),other.Cols(),
        ));
    }
    return nil;
}

func squareMatrixErrorCheck[N Number](m *Matrix[N]) error {
    if m.Cols()!=m.Rows() {
        return util.InverseOfNonSquareMatrix(
            fmt.Sprintf("[r1=%d c1=%d] | Need r1=c1",m.Rows(),m.Cols()),
        );
    }
    return nil;
}

//=============================================================================
//These functions should *not* be used, and are only here for
//benchmarking in the associated matrix test file.
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
