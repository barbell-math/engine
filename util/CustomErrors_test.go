package util;

import (
    "fmt"
    "testing"
)

func TestErrorEquality(t *testing.T){
    errs:=map[string]errorType{
        "SqlScriptNotFound": SqlScriptNotFound,
        "DataConversion": DataConversion,
        "NoKnownDataConversion": NoKnownDataConversion,
        "FilterRemovedAllColumns": FilterRemovedAllColumns,
        "DataVersionMalformed": DataVersionMalformed,
        "MatrixDimensionsDoNotAgree": MatrixDimensionsDoNotAgree,
        "InverseOfNonSquareMatrix": InverseOfNonSquareMatrix,
        "SingularMatrix": SingularMatrix,
        "MatrixSingularToWorkingPrecision": MatrixSingularToWorkingPrecision,
        "MissingVariable": MissingVariable,
    };
    isErrs:=map[string]isErrorType{
        "SqlScriptNotFound": IsSqlScriptNotFound,
        "DataConversion": IsDataConversion,
        "NoKnownDataConversion": IsNoKnownDataConversion,
        "FilterRemovedAllColumns": IsFilterRemovedAllColumns,
        "DataVersionMalformed": IsDataVersionMalformed,
        "MatrixDimensionsDoNotAgree": IsMatrixDimensionsDoNotAgree,
        "InverseOfNonSquareMatrix": IsInverseOfNonSquareMatrix,
        "SingularMatrix": IsSingularMatrix,
        "MatrixSingularToWorkingPrecision": IsMatrixSingularToWorkingPrecision,
        "MissingVariable": IsMissingVariable,
    };
    for k,_:=range(errs){
        iterErr:=errs[k]("testAddendum");
        if !isErrs[k](iterErr) {
            t.Errorf(fmt.Sprintf("%s is returning false negative.",k));
        }
        if isErrs[k](DataVersionNotAvailable) {
            t.Errorf(fmt.Sprintf("%s is returning false positive.",k));
        }
    }
}
