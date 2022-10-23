package util;

import (
    "fmt"
    "time"
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
)

type csvTest struct {
    I int;
    I8 int8;
    Ui uint;
    Ui8 uint8;
    F32 float32;
    S string;
    S1 string;
    B bool;
    T time.Time;
};

func TestNonStruct(t *testing.T){
    cntr:=0;
    err:=CSVToStruct("../testData/csv/ValidCSV.csv",',',"",func(v *int){
        cntr++;
    });
    fmt.Println(err);
    testUtil.BasicTest(0, cntr,
        "CSV converter returned values it was not supposed to.",t,
    );
    if !IsNonStructValue(err) {
        testUtil.FormatError(NonStructValue(""),err,
            "Did not raise error with non-struct value.",t,
        );
    }
}

func TestValidStruct(t *testing.T){
    cntr:=0;
    baseTime,_:=time.Parse("01/02/2006","12/12/2012");
    err:=CSVToStruct("../testData/csv/ValidCSV.csv",',',"01/02/2006",
    func(v *csvTest){
        testUtil.BasicTest(cntr+1,v.I,"Did not parse 'I' column correctly.",t);
        testUtil.BasicTest(int8(-cntr-2),v.I8,
            "Did not parse 'I8' column correctly.",t,
        );
        testUtil.BasicTest(uint(cntr+100),v.Ui,
            "Did not parse 'Ui' column correctly.",t,
        );
        testUtil.BasicTest(uint8(cntr+101),v.Ui8,
            "Did not parse 'Ui8' column correctly.",t,
        );
        testUtil.BasicTest(fmt.Sprintf("str%d",cntr+1),v.S,
            "Did not parse 'S' column correctly.",t,
        );
        testUtil.BasicTest(fmt.Sprintf("str%d",cntr+1),v.S1,
            "Did not parse 'S1' column correctly.",t,
        );
        testUtil.BasicTest(cntr!=0, v.B,
            "Did not parse 'B' column correctly.",t,
        );
        testUtil.BasicTest(
            baseTime.Add(time.Hour*24*(time.Duration(cntr))),v.T,
            "Did not parse 'T' column correctly.",t,
        );
        cntr++;
    });
    testUtil.BasicTest(nil,err,
        "CSVToStruct encountered error when it shouldn't have.",t,
    );
    testUtil.BasicTest(2,cntr,"Did not extract all values from file.",t);
}

func TestMissingColumns(t *testing.T){
    cntr:=0;
    baseTime,_:=time.Parse("01/02/2006","12/12/2012");
    err:=CSVToStruct("../testData/csv/MissingColumns.csv",',',"01/02/2006",
    func(v *csvTest){
        testUtil.BasicTest(0, v.I,"Missing I column was not zero-initialized.",t);
        testUtil.BasicTest(int8(-cntr-2),v.I8,
            "Did not parse 'I8' column correctly.",t,
        );
        testUtil.BasicTest(uint(cntr+100),v.Ui,
            "Did not parse 'Ui' column correctly.",t,
        );
        testUtil.BasicTest(uint8(cntr+101),v.Ui8,
            "Did not parse 'Ui8' column correctly.",t,
        );
        testUtil.BasicTest(fmt.Sprintf("str%d",cntr+1),v.S,
            "Did not parse 'S' column correctly.",t,
        );
        testUtil.BasicTest("",v.S1,
            "Did not parse 'S1' column correctly.",t,
        );
        testUtil.BasicTest(cntr!=0, v.B,
            "Did not parse 'B' column correctly.",t,
        );
        testUtil.BasicTest(
            baseTime.Add(time.Hour*24*(time.Duration(cntr))),v.T,
            "Did not parse 'T' column correctly.",t,
        );
        cntr++;
    });
    testUtil.BasicTest(nil,err,
        "CSVToStruct encountered error when it shouldn't have.",t,
    );
    testUtil.BasicTest(2,cntr,"Did not extract all values from file.",t);
}

func TestMissingValues(t *testing.T){
    cntr:=0;
    baseTime,_:=time.Parse("01/02/2006","00/00/0000");
    err:=CSVToStruct("../testData/csv/MissingValues.csv",',',"01/02/2006",
    func(v *csvTest){
        testUtil.BasicTest(cntr+1, v.I,"Did not parse 'I' column correctly.",t);
        testUtil.BasicTest(int8(0),v.I8,
            "Blank 'I8' value was not zero initialized.",t,
        );
        testUtil.BasicTest(uint(cntr+100),v.Ui,
            "Did not parse 'Ui' column correctly.",t,
        );
        testUtil.BasicTest(uint8(cntr+101),v.Ui8,
            "Did not parse 'Ui8' column correctly.",t,
        );
        testUtil.BasicTest(fmt.Sprintf("str%d",cntr+1),v.S,
            "Did not parse 'S' column correctly.",t,
        );
        testUtil.BasicTest("",v.S1,
            "Blank 'S1' value was not zero initialized.",t,
        );
        testUtil.BasicTest(cntr!=0, v.B,
            "Did not parse 'B' column correctly.",t,
        );
        testUtil.BasicTest(baseTime,v.T,
            "Blank 'T' value was not zero initialized.",t,
        );
        cntr++;
    });
    testUtil.BasicTest(nil,err,
        "CSVToStruct encountered error when it shouldn't have.",t,
    );
    testUtil.BasicTest(2,cntr,"Did not extract all values from file.",t);
}

func TestMissingHeaders(t *testing.T){
    cntr:=0;
    err:=CSVToStruct("../testData/csv/MissingHeaders.csv",',',"01/02/2006",
    func(v *csvTest){
        cntr++;
    });
    fmt.Println(err);
    testUtil.BasicTest(0, cntr,
        "CSV converter returned values it was not supposed to.",t,
    );
    if !IsMalformedCSVFile(err) {
        testUtil.FormatError(MalformedCSVFile(""),err,
            "Did not raise error with no headers in CSV file.",t,
        );
    }
}

func TestMalformedTypes(t *testing.T){
    cntr:=0;
    err:=CSVToStruct("../testData/csv/MalformedInt.csv",',',"",
    func(v *csvTest){
        cntr++;
    });
    fmt.Println(err);
    testUtil.BasicTest(0, cntr,
        "CSV converter returned values it was not supposed to.",t,
    );
    if !IsMalformedCSVFile(err) {
        testUtil.FormatError(MalformedCSVFile(""),err,
            "Did not raise error with malformed integer.",t,
        );
    }
    err=CSVToStruct("../testData/csv/MalformedUint.csv",',',"",
    func(v *csvTest){
        cntr++;
    });
    fmt.Println(err);
    testUtil.BasicTest(0, cntr,
        "CSV converter returned values it was not supposed to.",t,
    );
    if !IsMalformedCSVFile(err) {
        testUtil.FormatError(MalformedCSVFile(""),err,
            "Did not raise error with malformed unsigned integer.",t,
        );
    }
    err=CSVToStruct("../testData/csv/MalformedFloat.csv",',',"",
    func(v *csvTest){
        cntr++;
    });
    fmt.Println(err);
    testUtil.BasicTest(0, cntr,
        "CSV converter returned values it was not supposed to.",t,
    );
    if !IsMalformedCSVFile(err) {
        testUtil.FormatError(MalformedCSVFile(""),err,
            "Did not raise error with malformed float.",t,
        );
    }
    err=CSVToStruct("../testData/csv/MalformedTime.csv",',',"01/02/2006",
    func(v *csvTest){
        cntr++;
    });
    fmt.Println(err);
    testUtil.BasicTest(0, cntr,
        "CSV converter returned values it was not supposed to.",t,
    );
    if !IsMalformedCSVFile(err) {
        testUtil.FormatError(MalformedCSVFile(""),err,
            "Did not raise error with malformed time.",t,
        );
    }
}
