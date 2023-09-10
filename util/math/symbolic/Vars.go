package symbolic;

import (
    "fmt"
	"github.com/barbell-math/engine/util/math"
)

type Vars[N math.Number] map[string]Symbol[N];

func (v Vars[N])Access(_var string) (Symbol[N],error){
    if v,ok:=v[_var]; ok {
        return v,nil;
    }
    return nil,math.MissingVariable(
        fmt.Sprintf("Requested: %s Have: %v",_var,v),
    );
}

func (v Vars[N])Add(other Symbol[N]) Symbol[N] {
    return v;
}
func (v Vars[N])Sub(other Symbol[N]) Symbol[N] {
    return v;
}
func (v Vars[N])Mul(other Symbol[N]) Symbol[N] {
    return v;
}
func (v Vars[N])Div(other Symbol[N]) Symbol[N] {
    return v;
}
