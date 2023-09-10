package symbolic

import "github.com/barbell-math/block/util/math"

type PolarVector[N math.Number] []Symbol[N];

func (p PolarVector[N])Add(other Symbol[N]) Symbol[N] {
    return p;
}

func (p PolarVector[N])Sub(other Symbol[N]) Symbol[N] {
    return p;
}

func (p PolarVector[N])Mul(other Symbol[N]) Symbol[N] {
    return p;
}

func (p PolarVector[N])Div(other Symbol[N]) Symbol[N] {
    return p;
}
