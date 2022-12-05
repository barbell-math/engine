package reflect;

import (
    customerr "github.com/barbell-math/block/util/err"
)

var NonStructValue,IsNonStructValue=customerr.ErrorFactory(
    "A struct value was expected but was not recieved.",
);

