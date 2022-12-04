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
