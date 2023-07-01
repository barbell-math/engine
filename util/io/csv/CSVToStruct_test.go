package csv;

import (
    "fmt"
    "time"
    "testing"
    "github.com/barbell-math/block/util/test"
    "github.com/barbell-math/block/util/algo/iter"
)

func TestNonStruct(t *testing.T){
    cntr,err:=CSVToStruct[int](
        CSVFileSplitter("testData/ValidCSV.csv",',','#'),"",
    ).Count();
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
    err:=CSVToStruct[csvTest](
        CSVFileSplitter("testData/ValidCSV.csv",',','#'),"01/02/2006",
    ).ForEach(func(index int, val csvTest) (iter.IteratorFeedback, error) {
        test.BasicTest(index+1,val.I,"Did not parse 'I' column correctly.",t);
        test.BasicTest(int8(-index-2),val.I8,
            "Did not parse 'I8' column correctly.",t,
        );
        test.BasicTest(uint(index+100),val.Ui,
            "Did not parse 'Ui' column correctly.",t,
        );
        test.BasicTest(uint8(index+101),val.Ui8,
            "Did not parse 'Ui8' column correctly.",t,
        );
        test.BasicTest(fmt.Sprintf("str%d",index+1),val.S,
            "Did not parse 'S' column correctly.",t,
        );
        test.BasicTest(fmt.Sprintf("str%d",index+1),val.S1,
            "Did not parse 'S1' column correctly.",t,
        );
        test.BasicTest(index!=0, val.B,
            "Did not parse 'B' column correctly.",t,
        );
        test.BasicTest(
            baseTime.Add(time.Hour*24*(time.Duration(index))),val.T,
            "Did not parse 'T' column correctly.",t,
        );
        cntr++;
        return iter.Continue, nil;
    });
    test.BasicTest(nil,err,
        "CSVToStruct encountered error when it shouldn't have.",t,
    );
    test.BasicTest(2,cntr,"Did not extract all values from file.",t);
}

func TestMissingColumns(t *testing.T){
    cntr:=0;
    baseTime,_:=time.Parse("01/02/2006","12/12/2012");
    err:=CSVToStruct[csvTest](
        CSVFileSplitter("testData/MissingColumns.csv",',','#'),"01/02/2006",
    ).ForEach(func(index int, val csvTest) (iter.IteratorFeedback, error) {
        test.BasicTest(0, val.I,"Missing I column was not zero-initialized.",t);
        test.BasicTest(int8(-index-2),val.I8,
            "Did not parse 'I8' column correctly.",t,
        );
        test.BasicTest(uint(index+100),val.Ui,
            "Did not parse 'Ui' column correctly.",t,
        );
        test.BasicTest(uint8(index+101),val.Ui8,
            "Did not parse 'Ui8' column correctly.",t,
        );
        test.BasicTest(fmt.Sprintf("str%d",index+1),val.S,
            "Did not parse 'S' column correctly.",t,
        );
        test.BasicTest("",val.S1,
            "Did not parse 'S1' column correctly.",t,
        );
        test.BasicTest(index!=0, val.B,
            "Did not parse 'B' column correctly.",t,
        );
        test.BasicTest(
            baseTime.Add(time.Hour*24*(time.Duration(index))),val.T,
            "Did not parse 'T' column correctly.",t,
        );
        cntr++;
        return iter.Continue,nil;
    });
    test.BasicTest(nil,err,
        "CSVToStruct encountered error when it shouldn't have.",t,
    );
    test.BasicTest(2,cntr,"Did not extract all values from file.",t);
}

func TestMissingValues(t *testing.T){
    cntr:=0;
    baseTime,_:=time.Parse("01/02/2006","00/00/0000");
    err:=CSVToStruct[csvTest](
        CSVFileSplitter("testData/MissingValues.csv",',','#'),"01/02/2006",
    ).ForEach(func(index int, val csvTest) (iter.IteratorFeedback, error) {
        test.BasicTest(index+1, val.I,"Did not parse 'I' column correctly.",t);
        test.BasicTest(int8(0),val.I8,
            "Blank 'I8' value was not zero initialized.",t,
        );
        test.BasicTest(uint(index+100),val.Ui,
            "Did not parse 'Ui' column correctly.",t,
        );
        test.BasicTest(uint8(index+101),val.Ui8,
            "Did not parse 'Ui8' column correctly.",t,
        );
        test.BasicTest(fmt.Sprintf("str%d",index+1),val.S,
            "Did not parse 'S' column correctly.",t,
        );
        test.BasicTest("",val.S1,
            "Blank 'S1' value was not zero initialized.",t,
        );
        test.BasicTest(index!=0, val.B,
            "Did not parse 'B' column correctly.",t,
        );
        test.BasicTest(baseTime,val.T,
            "Blank 'T' value was not zero initialized.",t,
        );
        cntr++;
        return iter.Continue,nil;
    });
    test.BasicTest(nil,err,
        "CSVToStruct encountered error when it shouldn't have.",t,
    );
    test.BasicTest(2,cntr,"Did not extract all values from file.",t);
}

func TestMissingHeaders(t *testing.T){
    cntr:=0;
    err:=CSVToStruct[csvTest](
        CSVFileSplitter("testData/MissingHeaders.csv",',','#'),"01/02/2006",
    ).ForEach(func(index int, val csvTest) (iter.IteratorFeedback, error) {
        cntr++;
        return iter.Continue,nil;
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
    err:=CSVToStruct[csvTest](
        CSVFileSplitter("testData/MalformedInt.csv",',','#'),"01/02/2006",
    ).ForEach(func(index int, val csvTest) (iter.IteratorFeedback, error) {
        cntr++;
        return iter.Continue,nil;
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
    err=CSVToStruct[csvTest](
        CSVFileSplitter("testData/MalformedUint.csv",',','#'),"01/02/2006",
    ).ForEach(func(index int, val csvTest) (iter.IteratorFeedback, error) {
        cntr++;
        return iter.Continue,nil;
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
    err=CSVToStruct[csvTest](
        CSVFileSplitter("testData/MalformedFloat.csv",',','#'),"01/02/2006",
    ).ForEach(func(index int, val csvTest) (iter.IteratorFeedback, error) {
        cntr++;
        return iter.Continue,nil;
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
    err=CSVToStruct[csvTest](
        CSVFileSplitter("testData/MalformedTime.csv",',','#'),"01/02/2006",
    ).ForEach(func(index int, val csvTest) (iter.IteratorFeedback, error) {
        cntr++;
        return iter.Continue,nil;
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
