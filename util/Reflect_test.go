package util;

import (
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
)

type testStruct struct {
    One int;
    Two int;
};

func TestNonStructGetName(t *testing.T){
    v:=0;
    name,err:=GetStructName(&v);
    testUtil.BasicTest("",name,"The name of a non struct type was returned.",t);
    if !IsNonStructValue(err) {
        testUtil.FormatError(NonStructValue(""),err,
            "Non struct val did not raise appropriate error.",t,
        );
    }
}

func TestGetStructName(t *testing.T){
    var s testStruct;
    name,err:=GetStructName(&s);
    testUtil.BasicTest(nil,err,
        "Getting name of struct returned error when it was not supposed to.",t,
    );
    testUtil.BasicTest("testStruct",name,"Name of the struct was not correct.",t);
}

func TestNonStructGetFieldNames(t *testing.T){
    v:=0;
    fName,err:=GetStructFieldNames(&v,NoFilter[string]);
    testUtil.BasicTest(0, len(fName),
        "The field names of a non struct type was returned.",t,
    );
    if !IsNonStructValue(err) {
        testUtil.FormatError(NonStructValue(""),err,
            "Non struct val did not raise appropriate error.",t,
        );
    }
}

func TestGetStructFieldNames(t *testing.T){
    var s testStruct;
    fNames,err:=GetStructFieldNames(&s,NoFilter[string]);
    testUtil.BasicTest(nil,err,
        "Getting struct field names returned error when it was not supposed to.",t,
    );
    testUtil.BasicTest(2,len(fNames),"Not all struct field names were returned.",t);
    testUtil.BasicTest("One",fNames[0],
        "First struct field name was not correct.",t,
    );
    testUtil.BasicTest("Two",fNames[1],
        "Second struct field name was not correct.",t,
    );
}

func TestNonStructGetStructVals(t *testing.T){
    v:=0;
    vals,err:=GetStructVals(&v,NoFilter[string]);
    testUtil.BasicTest(0, len(vals),
        "The vals of a non struct type was returned.",t,
    );
    if !IsNonStructValue(err) {
        testUtil.FormatError(NonStructValue(""),err,
            "Non struct val did not raise appropriate error.",t,
        );
    }
}

func TestGetStructVals(t *testing.T){
    var s testStruct=testStruct{One: 1, Two: 2};
    vals,err:=GetStructVals(&s,NoFilter[string]);
    testUtil.BasicTest(nil,err,
        "Getting struct field vals returned error when it was not supposed to.",t,
    );
    testUtil.BasicTest(2,len(vals),"Not all struct field vals were returned.",t);
    testUtil.BasicTest(s.One,vals[0].Interface().(int),
        "First struct field val was not correct.",t,
    );
    testUtil.BasicTest(s.Two,vals[1].Interface().(int),
        "Second struct field val was not correct.",t,
    );
}

func TestNonStructGetStructFieldPntrs(t *testing.T){
    v:=0;
    vals,err:=GetStructFieldPntrs(&v,NoFilter[string]);
    testUtil.BasicTest(0, len(vals),
        "The vals of a non struct type was returned.",t,
    );
    if !IsNonStructValue(err) {
        testUtil.FormatError(NonStructValue(""),err,
            "Non struct val did not raise appropriate error.",t,
        );
    }
}

func TestGetStructFieldPntrs(t *testing.T){
    var s testStruct=testStruct{One: 1, Two: 2};
    vals,err:=GetStructFieldPntrs(&s,NoFilter[string]);
    testUtil.BasicTest(nil,err,
        "Getting struct field vals returned error when it was not supposed to.",t,
    );
    testUtil.BasicTest(2,len(vals),"Not all struct field vals were returned.",t);
    testUtil.BasicTest(&s.One,vals[0].Interface().(*int),
        "First struct field val was not correct.",t,
    );
    testUtil.BasicTest(&s.Two,vals[1].Interface().(*int),
        "Second struct field val was not correct.",t,
    );
}
