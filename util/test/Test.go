package test;

import (
    "fmt"
    "testing"
)

func BasicTest(expected any, got any, base string, t *testing.T){
    if expected!=got {
        FormatError(expected,got,base,t);
    }
}
func FormatError(expected any, got any, base string, t *testing.T){
    t.Error(fmt.Sprintf("Err: %s Expected: '%v' Got: '%v'",base,expected,got));
}

func SlicesMatch[T any](actual []T, generated []T, t *testing.T){
    BasicTest(len(actual),len(generated),"Slices do not match in length.",t);
    min:=len(generated);
    if len(actual)<len(generated) {
        min=len(actual);
    }
    for i:=0; i<min; i++ {
        BasicTest(actual[i],generated[i],fmt.Sprintf(
            "Values do not match | Index: %d",i,
        ),t);
    }
}
