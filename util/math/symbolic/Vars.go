package symbolic;

import (
    "fmt"
	"github.com/barbell-math/block/util/math"
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
