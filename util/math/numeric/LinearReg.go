package math

import (
	"fmt"
)

//A summation op is the function that is associated with a variable
//Ex:
// Consider the following: y=b_1*x_1+b_2*x_2
//      b_1: first constant lin reg will find
//      b_2: second constant lin reg will find
//      x_1: the first 'summation op' associated with b_1 (In this case linear)
//      x_2: the second 'summation op' associated with b_2 (In this case linear)
//  It may help to think of linear reg as using the following generic form:
//      y=b_1*f_1(x_1,x_2,...,x_n)+b_2*f_2(x_1,...,x_n)+...+b_n*f_n(x_1,...,x_n)
//  Where each function f_n is a 'summation op'
type SummationOp[N Number] func(vals map[string]N) (N,error);
func ConstSummationOp[N Number](v N) SummationOp[N] {
    return func(vals map[string]N) (N,error){
        return v,nil;
    }
}
func LinearSummationOp[N Number](_var string) SummationOp[N] {
    return func(vals map[string]N) (N,error){
        return VarAcc[N](vals,_var);
    }
}
func NegatedLinearSummationOp[N Number](_var string) SummationOp[N] {
    return func(vals map[string]N) (N,error){
        val,err:=VarAcc[N](vals,_var);
        return -val,err;
    }
}

func VarAcc[N Number](vals map[string]N, _var string) (N,error){
    if v,ok:=vals[_var]; ok {
        return v,nil;
    }
    return N(0),MissingVariable(
        fmt.Sprintf("Requested: %s Have: %v",_var,vals),
    );
}

type SummationOpGen[N Number] func(
    iVars []string, dVar string,
) ([]SummationOp[N],SummationOp[N]);

func ConstSumOpGen[N Number](v N) SummationOpGen[N] {
    return func(iVars []string, dVar string) ([]SummationOp[N],SummationOp[N]) {
        rv:=make([]SummationOp[N],len(iVars));
        for i,_:=range(iVars) {
            rv[i]=ConstSummationOp[N](v);
        }
        return rv,ConstSummationOp[N](v);
    }
}
func LinearSumOpGen[N Number](
        iVars []string,
        dVar string) ([]SummationOp[N],SummationOp[N]) {
    rv:=make([]SummationOp[N],len(iVars));
    for i,v:=range(iVars) {
        rv[i]=LinearSummationOp[N](v);
    }
    return rv,LinearSummationOp[N](dVar);
}
func LinearSumOpGenWithError[N Number](
        iVars []string,
        dVar string) ([]SummationOp[N],SummationOp[N]) {
    rv:=make([]SummationOp[N],len(iVars)+1);
    for i,v:=range(iVars) {
        rv[i]=LinearSummationOp[N](v);
    }
    rv[len(rv)-1]=ConstSummationOp[N](1);
    return rv,LinearSummationOp[N](dVar);
}

type LinRegResult[N Number] struct {
    Matrix[N];
    Predict func(iVars map[string]N) (N,error);
};
func (l *LinRegResult[N])GetConstant(i int) N {
    if i<l.Matrix.Rows() {
        return l.Matrix.V[i][0];
    }
    return N(0);
}
func (l *LinearReg[N])genLinRegPredict(r *LinRegResult[N]){
    r.Predict=func(iVars map[string]N) (N,error) {
        var err error;
        var rv,v N=N(0),N(0);
        for i:=0; err==nil && i<len(l.iVarOps); i++ {
            v,err=l.iVarOps[i](iVars);
            rv+=r.Matrix.V[i][0]*v;
        }
        return rv,err;
    }
}

type LinearReg[N Number] struct {
    a Matrix[N];
    b Matrix[N];
    summationOps [][]SummationOp[N];
    iVarOps []SummationOp[N];
    dVarOp SummationOp[N];
};

func NewLinearReg[N Number](
        iVarOps []SummationOp[N],
        dVarOp SummationOp[N]) LinearReg[N] {
    var rv LinearReg[N];
    n:=len(iVarOps);
    rv.iVarOps=iVarOps;
    rv.dVarOp=dVarOp;
    rv.a=NewMatrix(n,n,ZeroFill[N]);
    rv.b=NewMatrix(n,1,ZeroFill[N]);
    rv.summationOps=make([][]SummationOp[N],n);
    for i,_:=range(rv.summationOps) {
        rv.summationOps[i]=make([]SummationOp[N],n+1);
    }
    rv.populateSumOps();
    return rv;
}

func (l *LinearReg[N])populateSumOps(){
    for i,_:=range(l.summationOps) {
        for j,_:=range(l.summationOps[i]) {
            if j<l.a.Cols() {
                l.summationOps[i][j]=l.aSumOpGen(i,j);
            } else {
                l.summationOps[i][j]=l.bSumOpGen(i,j-l.a.Cols());
            }
        }
    }
}

func (l *LinearReg[N])aSumOpGen(r int, c int) SummationOp[N] {
    return func(vals map[string]N) (N,error) {
        v1,err:=l.iVarOps[c](vals);
        if err!=nil {
            return 0, err;
        }
        v2,err:=l.iVarOps[r](vals);
        if err!=nil {
            return 0, err;
        }
        return v1*v2, err;
    }
}

func (l *LinearReg[N])bSumOpGen(r int, c int) SummationOp[N] {
    return func(vals map[string]N) (N,error) {
        v1,err:=l.dVarOp(vals);
        if err!=nil {
            return 0, err;
        }
        v2,err:=l.iVarOps[r](vals);
        if err!=nil {
            return 0, err;
        }
        return v1*v2, err;
    }
}

func (l *LinearReg[N])sumOpRows() int { return l.a.Rows(); }
func (l *LinearReg[N])sumOpCols() int { return l.a.Cols()+l.b.Cols(); }

func (l *LinearReg[N])IterSummationOps(f func(r int, c int, v SummationOp[N])){
    for i,v1:=range(l.summationOps) {
        for j,v2:=range(v1){
            f(i,j,v2);
        }
    }
}

func (l *LinearReg[N])IterLHS(f func(r int, c int, v N)){
    l.a.Iter(f);
}

func (l *LinearReg[N])IterRHS(f func(r int, c int, v N)){
    l.b.Iter(f);
}

func (l *LinearReg[N])UpdateSummations(vals map[string]N) error {
    for i,r:=range(l.summationOps) {
        for j,s:=range(r) {
            if v,err:=s(vals); err==nil {
                if j<l.a.Cols() {
                    l.a.V[i][j]+=v;
                } else {
                    l.b.V[i][j-l.a.Cols()]+=v;
                }
            } else {
                return err;
            }
        }
    }
    return nil;
}

func (l *LinearReg[N])Run() (LinRegResult[N],float64,error) {
    var rv LinRegResult[N];
    rv.Matrix=l.a.Copy();
    rcond,err:=rv.Matrix.Inverse();
    if !IsInverseOfNonSquareMatrix(err) {
        //err in RV can be ignored, matrices are guaranteed to have correct
        //dimensions because they are only managed by the linear reg struct
        rv.Matrix.Mul(&l.b);
    }
    l.genLinRegPredict(&rv);
    return rv,rcond,err;
}
