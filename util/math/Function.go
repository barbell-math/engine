package math;

import (
    "fmt"
)

type Vars[N Number] map[string]Symbol[N];
type Symbol[N Number] interface {
    Add(other Symbol[N]) (Symbol[N],error)
};

type Scalar[N Number] struct { v N };
type Vector[N Number] []Symbol[N];
type Matrix[N Number] [][]Symbol[N];
type Equation[N Number] func(iVars Vars[N]) (Symbol[N],error);

func (v Vars[N])Access(_var string) (Symbol[N],error){
    if v,ok:=v[_var]; ok {
        return v,nil;
    }
    return nil,MissingVariable(
        fmt.Sprintf("Requested: %s Have: %v",_var,v),
    );
}

func (v Vars[N])Apply(other Vars[N], op func(accum *N, iter N) error)  error {
    for oKey,oVal:=range(other) {
        if val,err:=v.Access(oKey); err==nil {
            if err:=op(&val,oVal); err==nil {
                v[oKey]=val;
            } else {
                return err;
            }
        } else {
            return err;
        }
    }
    return nil;
}

func (v Vars[N])Copy() Vars[N] {
    rv:=make(map[string]N,len(v));
    for k,v:=range(v) {
        rv[k]=v;
    }
    return rv;
}

func (f Function[N])Negate() Function[N] {
    return func(iVars Vars[N], err error) (Vars[N], error) {
        res,err:=f(iVars);
        for k,_:=range(res) {
            
        }
    }    
}
