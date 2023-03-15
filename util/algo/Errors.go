package algo;

import (
    customerr "github.com/barbell-math/block/util/err"
)

var SliceZippingError,IsSliceZippingError=customerr.ErrorFactory(
    "The given slices could not be zipped.",
);
