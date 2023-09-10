package symbolic;

import (
	"fmt"
	"reflect"

	"github.com/barbell-math/block/util/math"
)

type SymbolicError[N math.Number] struct { error; }

func invalidBinaryOpFormater[N math.Number](
    s1 Symbol[N], 
    s2 Symbol[N], 
    op string,
) SymbolicError[N] {
    return SymbolicError[N]{InvalidOperation(fmt.Sprintf(
        "%s %s %s",reflect.TypeOf(s1),op,reflect.TypeOf(s2),
    ))};
}

func (e SymbolicError[N])Add(other Symbol[N]) Symbol[N] { return e; }
func (e SymbolicError[N])Sub(other Symbol[N]) Symbol[N] { return e; }
func (e SymbolicError[N])Mul(other Symbol[N]) Symbol[N] { return e; }
func (e SymbolicError[N])Div(other Symbol[N]) Symbol[N] { return e; }
