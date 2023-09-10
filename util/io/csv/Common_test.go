package csv;

import (
    "fmt"
    "time"
    "testing"
    "github.com/barbell-math/engine/util/test"
    "github.com/barbell-math/engine/util/algo/iter"
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

func TestFlatten(t *testing.T) {
    res,err:=Flatten(iter.SliceElems([][]string{}),",").Collect();
    test.BasicTest(nil,err,
        "Flatten returned an error when it was not supposed to",t,
    );
    test.BasicTest(0,len(res),
        "Flatten did not produce the correct value.",t,
    );
    test.BasicTest(nil,err,
        "Flatten returned an error when it was not supposed to",t,
    );
    res,err=Flatten(iter.SliceElems([][]string{
        {"1"},
        {"2"},
        {"3"},
        {"4"},
    }),",").Collect();
    test.BasicTest(nil,err,
        "Flatten returned an error when it was not supposed to",t,
    );
    for i,v:=range([]string{"1","2","3","4"}) {
        test.BasicTest(v,res[i],
            "Flatten did not produce the correct value.",t,
        );
    }
    res,err=Flatten(iter.SliceElems([][]string{
        {"1","2","3"},
        {"2","3","4"},
        {"3","4","5"},
        {"4","5","6"},
    }),",").Collect();
    test.BasicTest(nil,err,
        "Flatten returned an error when it was not supposed to",t,
    );
    for i,v:=range([]string{"1,2,3","2,3,4","3,4,5","4,5,6"}) {
        test.BasicTest(v,res[i],
            "Flatten did not produce the correct value.",t,
        );
    }
}

func TestCSVFileSplitter(t *testing.T){
    err:=CSVFileSplitter("./testData/ValidCSVTemplate.csv",',','#').ForEach(
    func(index int, val []string) (iter.IteratorFeedback, error) {
        if index==0 {
            for i,v:=range(val) {
                test.BasicTest(fmt.Sprintf("Column%d",i+1),v,
                    "Column of header row was not correct.",t,
                );
            }
        } else {
            for i,v:=range(val) {
                test.BasicTest(fmt.Sprintf("%d",index+i),v,
                    "Column of header row was not correct.",t,
                );
            }
        }
        return iter.Continue,nil;
    });
    test.BasicTest(nil,err,
        "CSV iterator returned an error when it should not have.",t,
    );
}

func TestCSVFileSplitterSkipHeaders(t *testing.T){
    err:=CSVFileSplitter("./testData/ValidCSVTemplate.csv",',','#').Skip(1).
    ForEach(func(index int, val []string) (iter.IteratorFeedback, error) {
        for i,v:=range(val) {
            test.BasicTest(fmt.Sprintf("%d",index+i+1),v,
                "Column of header row was not correct.",t,
            );
        }
        return iter.Continue,nil;
    });
    test.BasicTest(nil,err,
        "CSV iterator returned an error when it should not have.",t,
    );
}

func TestCSVFileSplitterWithStrings(t *testing.T){
    err:=CSVFileSplitter("./testData/ValidCSVTemplateStringVals.csv",',','#').
    Take(1).ForEach(func(index int, val []string) (iter.IteratorFeedback, error) {
        for i,v:=range(val) {
            test.BasicTest(fmt.Sprintf("Column%d,\"Column%d\"",i+1,i+2),v,
                "Column of header row was not correct.",t,
            );
        }
        return iter.Continue,nil;
    });
    test.BasicTest(nil,err,
        "CSV iterator returned an error when it should not have.",t,
    );
}


