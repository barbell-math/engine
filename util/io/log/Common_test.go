package log

import (
	"fmt"
	"testing"

	"github.com/barbell-math/block/util/algo/iter"
	"github.com/barbell-math/block/util/dataStruct"
	"github.com/barbell-math/block/util/dataStruct/types"
	"github.com/barbell-math/block/util/test"
)

func generateLog(l *Logger[int], numLines int) {
    for i:=0; i<numLines; i++ {
        l.Log(fmt.Sprintf("Line %d",i),i);
    }
}

func generateIntertwinedLogs(l1 *Logger[int], l2 *Logger[int], numLines int) {
    cntr:=0;
    for i:=0; i<numLines; i++ {
        cntr++;
        l1.Log(fmt.Sprintf("L1 Line %d",cntr),cntr);
        cntr++;
        l2.Log(fmt.Sprintf("L2 Line %d",cntr),cntr);
    }
}

func TestNewLogBadPath(t *testing.T) {
    _,err:=NewLog[int](Error,"./non/existant/path/to/file.txt",false);
    if err==nil {
        test.FormatError("!<nil>",err,"An error was expected to be returned.",t);
    }
}

func TestLogIteration(t *testing.T){
    l,_:=NewLog[int](Error,"./testData/generateLog.log",false);
    generateLog(&l,1000);
    LogElems(l).ForEach(func(index int, val LogEntry[int]) (iter.IteratorFeedback, error) {
        test.BasicTest(Error,val.Status,"Log line status was not set properly.",t);
        test.BasicTest(fmt.Sprintf("Line %d",index),val.Message,
            "Log line message was not set properly.",t,
        );
        test.BasicTest(index,val.Val,"Log line value was not set properly.",t);
        return iter.Continue,nil;
    });
    l.Close();
}

func TestLogIterationTime(t *testing.T){
    l,_:=NewLog[int](Error,"./testData/generateLog.log",false);
    generateLog(&l,1000);
    cntr:=0;
    c,_:=dataStruct.NewCircularQueue[LogEntry[int]](2);
    err:=iter.Window[LogEntry[int]](LogElems(l),&c,false).ForEach(
    func(index int, q types.Queue[LogEntry[int]]) (iter.IteratorFeedback, error) {
        cntr++;
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
    });
    test.BasicTest(999,cntr,
        "Window did not produce correct number of elements.",t,
    );
    test.BasicTest(nil,err,
        "Window returned an error when it should not have.",t,
    );
    l.Close();
}

func TestLogAppend(t *testing.T){
    l,_:=NewLog[int](Error,"./testData/generateLog.log",false);
    generateLog(&l,1000);
    l.Close();
    l,_=NewLog[int](Error,"./testData/generateLog.log",true);
    generateLog(&l,1000);
    cntr:=0;
    err:=LogElems(l).ForEach(
    func(index int, val LogEntry[int]) (iter.IteratorFeedback, error) {
        test.BasicTest(Error,val.Status,"Log line status was not set properly.",t);
        test.BasicTest(fmt.Sprintf("Line %d",index%1000),val.Message,
            "Log line message was not set properly.",t,
        );
        test.BasicTest(index%1000,val.Val,
            "Log line value was not set properly.",t,
        );
        cntr++;
        return iter.Continue,nil;
    });
    test.BasicTest(2000,cntr,"Log was not appended to correctly.",t);
    test.BasicTest(nil,err,
        "Log elems returned an error when it was not supposed to.",t,
    );
    l.Close();
}

func TestLogJoin(t *testing.T){
    cntr:=0;
    l1,_:=NewLog[int](Error,"./testData/generateLog.part1.log",false);
    l2,_:=NewLog[int](Error,"./testData/generateLog.part2.log",false);
    generateIntertwinedLogs(&l1,&l2,1000);
    l1S,_:=NewLog[int](Error,"./testData/generateLog.part1.log",true);
    l2S,_:=NewLog[int](Error,"./testData/generateLog.part2.log",true);
    err:=iter.JoinSame[LogEntry[int]](LogElems(l1S),LogElems(l2S),
        dataStruct.Variant[LogEntry[int],LogEntry[int]]{},
        JoinLogByTime[int,int],
    ).ForEach(
    func(index int, val LogEntry[int]) (iter.IteratorFeedback, error) {
        test.BasicTest(Error,val.Status,"Log line status was not set properly.",t);
        test.BasicTest(fmt.Sprintf("L%d Line %d",index%2+1,index+1),val.Message,
            "Log line message was not set properly.",t,
        );
        test.BasicTest(index+1,val.Val,
            "Log line value was not set properly.",t,
        );
        cntr++;
        return iter.Continue,nil;
    });
    test.BasicTest(2000,cntr,"Log was not appended to correctly.",t);
    test.BasicTest(nil,err,
        "Log elems returned an error when it was not supposed to.",t,
    );
    l1.Close();
    l2.Close();
}


func BenchmarkJoinLog(b *testing.B) {
    l1,_:=NewLog[int](Error,"./testData/generateLog.part1.log",false);
    l2,_:=NewLog[int](Error,"./testData/generateLog.part2.log",false);
    generateIntertwinedLogs(&l1,&l2,1000);
    l1S,_:=NewLog[int](Error,"./testData/generateLog.part1.log",true);
    l2S,_:=NewLog[int](Error,"./testData/generateLog.part2.log",true);
    for i:=0; i<b.N; i++ {
        iter.JoinSame[LogEntry[int]](LogElems(l1S),LogElems(l2S),
            dataStruct.Variant[LogEntry[int],LogEntry[int]]{},
            JoinLogByTime[int,int],
        ).ForEach(
        func(index int, val LogEntry[int]) (iter.IteratorFeedback, error) {
            return iter.Continue,nil;
        });
    }
    l1.Close();
    l2.Close();
}

func BenchmarkBigLog(b *testing.B) {
    cntr:=0;
    l,_:=NewLog[int](Error,"./testData/big.log",false);
    fmt.Println("Generating log file with 500K lines (this may take a while)...");
    generateLog(&l,500000);
    fmt.Println("Running tests...");
    for i:=0; i<b.N; i++ {
        LogElems(l).ForEach(
        func(index int, val LogEntry[int]) (iter.IteratorFeedback, error) {
            cntr++;
            return iter.Continue,nil;
        });
        fmt.Printf("Processed %d lines.\n",cntr);
    }
    l.Close();
}
