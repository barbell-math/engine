package csv;

import (
    "fmt"
    "time"
    "testing"
    "github.com/barbell-math/block/util/test"
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
    err:=CSVToStruct("testData/ValidCSV.csv",',',"",func(v *int){
        cntr++;
    });
    fmt.Println(err);
    test.BasicTest(0, cntr,
        "CSV converter returned values it was not supposed to.",t,
    );
    if !IsNonStructValue(err) {
        test.FormatError(NonStructValue(""),err,
            "Did not raise error with non-struct value.",t,
        );
    }
}

func TestValidStruct(t *testing.T){
    cntr:=0;
    baseTime,_:=time.Parse("01/02/2006","12/12/2012");
    err:=CSVToStruct("testData/ValidCSV.csv",',',"01/02/2006",
    func(v *csvTest){
        test.BasicTest(cntr+1,v.I,"Did not parse 'I' column correctly.",t);
        test.BasicTest(int8(-cntr-2),v.I8,
            "Did not parse 'I8' column correctly.",t,
        );
        test.BasicTest(uint(cntr+100),v.Ui,
            "Did not parse 'Ui' column correctly.",t,
        );
        test.BasicTest(uint8(cntr+101),v.Ui8,
            "Did not parse 'Ui8' column correctly.",t,
        );
        test.BasicTest(fmt.Sprintf("str%d",cntr+1),v.S,
            "Did not parse 'S' column correctly.",t,
        );
        test.BasicTest(fmt.Sprintf("str%d",cntr+1),v.S1,
            "Did not parse 'S1' column correctly.",t,
        );
        test.BasicTest(cntr!=0, v.B,
            "Did not parse 'B' column correctly.",t,
        );
        test.BasicTest(
            baseTime.Add(time.Hour*24*(time.Duration(cntr))),v.T,
            "Did not parse 'T' column correctly.",t,
        );
        cntr++;
    });
    test.BasicTest(nil,err,
        "CSVToStruct encountered error when it shouldn't have.",t,
    );
    test.BasicTest(2,cntr,"Did not extract all values from file.",t);
}

func TestMissingColumns(t *testing.T){
    cntr:=0;
    baseTime,_:=time.Parse("01/02/2006","12/12/2012");
    err:=CSVToStruct("testData/MissingColumns.csv",',',"01/02/2006",
    func(v *csvTest){
        test.BasicTest(0, v.I,"Missing I column was not zero-initialized.",t);
        test.BasicTest(int8(-cntr-2),v.I8,
            "Did not parse 'I8' column correctly.",t,
        );
        test.BasicTest(uint(cntr+100),v.Ui,
            "Did not parse 'Ui' column correctly.",t,
        );
        test.BasicTest(uint8(cntr+101),v.Ui8,
            "Did not parse 'Ui8' column correctly.",t,
        );
        test.BasicTest(fmt.Sprintf("str%d",cntr+1),v.S,
            "Did not parse 'S' column correctly.",t,
        );
        test.BasicTest("",v.S1,
            "Did not parse 'S1' column correctly.",t,
        );
        test.BasicTest(cntr!=0, v.B,
            "Did not parse 'B' column correctly.",t,
        );
        test.BasicTest(
            baseTime.Add(time.Hour*24*(time.Duration(cntr))),v.T,
            "Did not parse 'T' column correctly.",t,
        );
        cntr++;
    });
    test.BasicTest(nil,err,
        "CSVToStruct encountered error when it shouldn't have.",t,
    );
    test.BasicTest(2,cntr,"Did not extract all values from file.",t);
}

func TestMissingValues(t *testing.T){
    cntr:=0;
    baseTime,_:=time.Parse("01/02/2006","00/00/0000");
    err:=CSVToStruct("testData/MissingValues.csv",',',"01/02/2006",
    func(v *csvTest){
        test.BasicTest(cntr+1, v.I,"Did not parse 'I' column correctly.",t);
        test.BasicTest(int8(0),v.I8,
            "Blank 'I8' value was not zero initialized.",t,
        );
        test.BasicTest(uint(cntr+100),v.Ui,
            "Did not parse 'Ui' column correctly.",t,
        );
        test.BasicTest(uint8(cntr+101),v.Ui8,
            "Did not parse 'Ui8' column correctly.",t,
        );
        test.BasicTest(fmt.Sprintf("str%d",cntr+1),v.S,
            "Did not parse 'S' column correctly.",t,
        );
        test.BasicTest("",v.S1,
            "Blank 'S1' value was not zero initialized.",t,
        );
        test.BasicTest(cntr!=0, v.B,
            "Did not parse 'B' column correctly.",t,
        );
        test.BasicTest(baseTime,v.T,
            "Blank 'T' value was not zero initialized.",t,
        );
        cntr++;
    });
    test.BasicTest(nil,err,
        "CSVToStruct encountered error when it shouldn't have.",t,
    );
    test.BasicTest(2,cntr,"Did not extract all values from file.",t);
}

func TestMissingHeaders(t *testing.T){
    cntr:=0;
    err:=CSVToStruct("testData/MissingHeaders.csv",',',"01/02/2006",
    func(v *csvTest){
        cntr++;
    });
    fmt.Println(err);
    test.BasicTest(0, cntr,
        "CSV converter returned values it was not supposed to.",t,
    );
    if !IsMalformedCSVFile(err) {
        test.FormatError(MalformedCSVFile(""),err,
            "Did not raise error with no headers in CSV file.",t,
        );
    }
}

func TestMalformedTypes(t *testing.T){
    cntr:=0;
    err:=CSVToStruct("testData/MalformedInt.csv",',',"",
    func(v *csvTest){
        cntr++;
    });
    fmt.Println(err);
    test.BasicTest(0, cntr,
        "CSV converter returned values it was not supposed to.",t,
    );
    if !IsMalformedCSVFile(err) {
        test.FormatError(MalformedCSVFile(""),err,
            "Did not raise error with malformed integer.",t,
        );
    }
    err=CSVToStruct("testData/MalformedUint.csv",',',"",
    func(v *csvTest){
        cntr++;
    });
    fmt.Println(err);
    test.BasicTest(0, cntr,
        "CSV converter returned values it was not supposed to.",t,
    );
    if !IsMalformedCSVFile(err) {
        test.FormatError(MalformedCSVFile(""),err,
            "Did not raise error with malformed unsigned integer.",t,
        );
    }
    err=CSVToStruct("testData/MalformedFloat.csv",',',"",
    func(v *csvTest){
        cntr++;
    });
    fmt.Println(err);
    test.BasicTest(0, cntr,
        "CSV converter returned values it was not supposed to.",t,
    );
    if !IsMalformedCSVFile(err) {
        test.FormatError(MalformedCSVFile(""),err,
            "Did not raise error with malformed float.",t,
        );
    }
    err=CSVToStruct("testData/MalformedTime.csv",',',"01/02/2006",
    func(v *csvTest){
        cntr++;
    });
    fmt.Println(err);
    test.BasicTest(0, cntr,
        "CSV converter returned values it was not supposed to.",t,
    );
    if !IsMalformedCSVFile(err) {
        test.FormatError(MalformedCSVFile(""),err,
            "Did not raise error with malformed time.",t,
        );
    }
}
