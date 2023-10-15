package potentialSurface;

import (
    customerr "github.com/barbell-math/engine/util/err"
)

var InvalidPotentialSurfaceId,IsInvalidPotentialSurfaceId=customerr.ErrorFactory(
    "The supplied potential surface id is not mapped to any surface.",
);
