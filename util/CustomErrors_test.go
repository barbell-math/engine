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
    };
    isErrs:=map[string]isErrorType{
        "SqlScriptNotFound": IsSqlScriptNotFound,
        "DataConversion": IsDataConversion,
        "NoKnownDataConversion": IsNoKnownDataConversion,
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
