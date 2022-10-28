package util;

import (
    "fmt"
    "errors"
    "strings"
)

var DataVersionNotAvailable=errors.New("Data version could not be selected.");

type errorType func(addendum string) error;
type isErrorType func(err error) bool;
func errorFactory(base string) (errorType,isErrorType) {
    errorTypeRv:=func(addendum string) error {
        return errors.New(fmt.Sprintf("%s | %s",base,addendum));
    }
    isErrorTypeRv:=func(err error) bool {
        return strings.Contains(fmt.Sprint(err),base);
    }
    return errorTypeRv,isErrorTypeRv;
}

var UnsupportedQueryType,IsUnsupportedQueryType=errorFactory(
    "The supplied query type is not supported for this operation.",
);

var SqlScriptNotFound,IsSqlScriptNotFound=errorFactory(
    "Could not open SQL script to run queries.",
);

var DataConversion,IsDataConversion=errorFactory(
    "An error occurred converting data between versions.",
);

var NoKnownDataConversion,IsNoKnownDataConversion=errorFactory(
    "No known data conversion.",
);

var FilterRemovedAllColumns,IsFilterRemovedAllColumns=errorFactory(
    "The filter passed resulted in no columns being selected.",
);

var DataVersionMalformed,IsDataVersionMalformed=errorFactory(
    "Data version is malformed.",
);
var SettingsFileNotFound,IsSettingsFileNotFound=errorFactory(
    "A file specified in the settings config file could not be found.",
);

var MatrixDimensionsDoNotAgree,IsMatrixDimensionsDoNotAgree=errorFactory(
    "Matrix dimensions do not agree.",
);

var InverseOfNonSquareMatrix,IsInverseOfNonSquareMatrix=errorFactory(
    "Only square matrices have inverses.",
);

var SingularMatrix,IsSingularMatrix=errorFactory(
    "The det of the matrix is zero. Some operations like inverse will produce erroneous results.",
);

var MatrixSingularToWorkingPrecision,IsMatrixSingularToWorkingPrecision=errorFactory(
    "Calculations resulted in a matrix that is <= working precision.",
);

var MissingVariable,IsMissingVariable=errorFactory(
    "The requested independent variable is not present.",
);

var MalformedCSVFile,IsMalformedCSVFile=errorFactory(
    "The CSV file cannot be converted to the requested struct.",
);

var NonStructValue,IsNonStructValue=errorFactory(
    "A struct value was expected but was not recieved.",
);

var SliceZippingError,IsSliceZippingError=errorFactory(
    "The given slices could not be zipped.",
);
