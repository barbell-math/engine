package math;

import (
    customerr "github.com/barbell-math/block/util/err"
)

var DivByZero,IsDivByZero=customerr.ErrorFactory(
    "Attempted division by zero.",
);

var DimensionsDoNotAgree,IsDimensionsDoNotAgree=customerr.ErrorFactory(
    "Dimensions do not agree.",
);

var MatrixDimensionsDoNotAgree,IsMatrixDimensionsDoNotAgree=customerr.ErrorFactory(
    "Matrix dimensions do not agree.",
);

var InverseOfNonSquareMatrix,IsInverseOfNonSquareMatrix=customerr.ErrorFactory(
    "Only square matrices have inverses.",
);

var SingularMatrix,IsSingularMatrix=customerr.ErrorFactory(
    "The det of the matrix is zero. Some operations like inverse will produce erroneous results.",
);

var MatrixSingularToWorkingPrecision,IsMatrixSingularToWorkingPrecision=customerr.ErrorFactory(
    "Calculations resulted in a matrix that is <= working precision.",
);

var MissingVariable,IsMissingVariable=customerr.ErrorFactory(
    "The requested independent variable is not present.",
);

