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
