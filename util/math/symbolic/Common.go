package symbolic

import (
	"github.com/barbell-math/block/util/math"
)

type Symbol[N math.Number] interface {
    Add(other Symbol[N]) Symbol[N];
    Sub(other Symbol[N]) Symbol[N];
    Mul(other Symbol[N]) Symbol[N];
    //Cross
    //Dot
};
type Matrix[N math.Number] [][]Symbol[N];
type Equation[N math.Number] func(iVars Vars[N]) (Symbol[N],error);
