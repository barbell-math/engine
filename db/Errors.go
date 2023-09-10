package db;

import (
    "errors"
    customerr "github.com/barbell-math/engine/util/err"
)

var DataVersionNotAvailable=errors.New("Data version could not be selected.");

var UnsupportedQueryType,IsUnsupportedQueryType=customerr.ErrorFactory(
    "The supplied query type is not supported for this operation.",
);

var SqlScriptNotFound,IsSqlScriptNotFound=customerr.ErrorFactory(
    "Could not open SQL script to run queries.",
);

var DataConversion,IsDataConversion=customerr.ErrorFactory(
    "An error occurred converting data between versions.",
);

var NoKnownDataConversion,IsNoKnownDataConversion=customerr.ErrorFactory(
    "No known data conversion.",
);

var FilterRemovedAllColumns,IsFilterRemovedAllColumns=customerr.ErrorFactory(
    "The filter passed resulted in no columns being selected.",
);
