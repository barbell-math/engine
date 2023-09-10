package symbolic;

import (
    customerr "github.com/barbell-math/block/util/err"
)

var InvalidOperation,IsInvalidOperation=customerr.ErrorFactory(
    "The supplied types cannot perform the requested operation.",
);
