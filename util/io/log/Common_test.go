package log

import (
	"fmt"
	"testing"
	"github.com/barbell-math/block/util/test"
	"github.com/barbell-math/block/util/algo/iter"
	"github.com/barbell-math/block/util/dataStruct"
	"github.com/barbell-math/block/util/dataStruct/types"
)

func generateLog(l *Logger[int], numLines int) {
    for i:=0; i<numLines; i++ {
        l.Log(fmt.Sprintf("Line %d",i),i);
    }
}

func TestLogIteration(t *testing.T){
    l:=NewLog[int](Error,"./testData/generateLog.log");
    generateLog(&l,1000);
    LogElems[int](l).ForEach(func(index int, val LogEntry[int]) (iter.IteratorFeedback, error) {
        test.BasicTest(Error,val.Status,"Log line status was not set properly.",t);
        test.BasicTest(fmt.Sprintf("Line %d",index),val.Message,
            "Log line message was not set properly.",t,
        );
        test.BasicTest(index,val.Val,"Log line value was not set properly.",t);
        return iter.Continue,nil;
    });
}

func TestLogIterationTime(t *testing.T){
    l:=NewLog[int](Error,"./testData/generateLog.log");
    generateLog(&l,1000);
    c,_:=dataStruct.NewCircularQueue[LogEntry[int]](2);
    cnt,err:=iter.Window[LogEntry[int]](LogElems(l),&c,false,
    func(index int, q types.Queue[LogEntry[int]]) (iter.IteratorFeedback, error) {
        fmt.Printf("%v",q);
        f,_:=q.Peek(0);
        s,_:=q.Peek(1);
        if !f.Time.Before(s.Time) {
            test.FormatError(
                "t1<t2",
                fmt.Sprintf("t1='%v' >= t2='%v'",f.Time,s.Time),
                "Time did not increase as it ought to.",t,
            );
        }
        return iter.Continue,nil;
    }).Count();
    test.BasicTest(999,cnt,
        "Window did not produce correct number of elements.",t,
    );
    test.BasicTest(nil,err,
        "Window returned an error when it should not have.",t,
    );
}

func BenchmarkBigLog(b *testing.B) {
    cntr:=0;
    l:=NewLog[int](Error,"./testData/big.log");
    fmt.Println("Generating log file with 1M lines...");
    //generateLog(&l,1000000);
    fmt.Println("Running tests...");
    for i:=0; i<b.N; i++ {
        LogElems(l).ForEach(
        func(index int, val LogEntry[int]) (iter.IteratorFeedback, error) {
            cntr++;
            return iter.Continue,nil;
        });
        fmt.Println(cntr);
    }
}
