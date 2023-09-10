package dataStruct

import (
	"testing"

	"github.com/barbell-math/block/util/algo/iter"
	"github.com/barbell-math/block/util/dataStruct/types"
	customerr "github.com/barbell-math/block/util/err"
	"github.com/barbell-math/block/util/test"
)

func TestNewCircularQueue(t *testing.T) {
    tmp,err:=NewCircularQueue[int](5);
    test.BasicTest(nil,err,
        "NewCircularQueue returned an error when it should not have.",t,
    );
    test.BasicTest(5,len(tmp.vals),
        "NewCircularQueue added values to empty queue during initialization.",t,
    );
    test.BasicTest(5,cap(tmp.vals),
        "NewCircularQueue did not set capacity correctly.",t,
    );
    test.BasicTest(0,tmp.startEnd.A,
        "NewCircularQueue added values to empty queue during initialization.",t,
    );
    test.BasicTest(-1,tmp.startEnd.B,
        "NewCircularQueue added values to empty queue during initialization.",t,
    );
}

func TestNewCircularQueueBadSize(t *testing.T) {
    tmp,err:=NewCircularQueue[int](0);
    if !customerr.IsValOutsideRange(err) {
        test.FormatError(customerr.ValOutsideRange(""),err,
            "Invalid queue size did not raise the correct error.",t,
        );
    }
    test.BasicTest(0,len(tmp.vals),
        "NewCircularQueue added values to empty queue during initialization.",t,
    );
    test.BasicTest(0,cap(tmp.vals),
        "NewCircularQueue did not set capacity correctly.",t,
    );
    test.BasicTest(0,tmp.startEnd.A,
        "NewCircularQueue added values to empty queue during initialization.",t,
    );
    test.BasicTest(0,tmp.startEnd.B,
        "NewCircularQueue added values to empty queue during initialization.",t,
    );
}

func TestCircularQueuePush(t *testing.T){
    tmp,err:=NewCircularQueue[int](5);
    test.BasicTest(nil,err,
        "NewCircularQueue returned an error when it should not have.",t,
    );
    for i:=0; i<5; i++ {
        res:=tmp.Push(i);
        test.BasicTest(nil,res,
            "Push returned an error when it should not have.",t,
        );
        test.BasicTest(i+1,tmp.NumElems,
            "Push did not increment NumElems after adding value.",t,
        );
        test.BasicTest(i,tmp.vals[i],"Push did not save value.",t);
        test.BasicTest(0,tmp.startEnd.A,"Push modified the start index.",t);
        test.BasicTest(i,tmp.startEnd.B,
            "Push did not modify the end index.",t,
        );
    }
    res:=tmp.Push(5);
    if !IsQueueFull(res) {
        test.FormatError(QueueFull(""),res,
            "Push did not detect the queue was full.",t,
        );
    }
    test.BasicTest(5,tmp.NumElems,
        "Push incremented NumElems when queue was full.",t,
    );
    test.BasicTest(0,tmp.startEnd.A,"Push modified the start index.",t);
    test.BasicTest(len(tmp.vals)-1,tmp.startEnd.B,
        "Push modified the end index when the queue was full.",t,
    );
}

func TestCircularQueuePushStartFromMiddle(t *testing.T) {
    tmp,err:=NewCircularQueue[int](5);
    test.BasicTest(nil,err,
        "NewCircularQueue returned an error when it should not have.",t,
    );
    tmp.startEnd.A=2;
    tmp.startEnd.B=1;
    for i:=0; i<5; i++ {
        res:=tmp.Push(i);
        test.BasicTest(nil,res,
            "Push returned an error when it should not have.",t,
        );
        test.BasicTest(i+1,tmp.NumElems,
            "Push did not increment NumElems after adding value.",t,
        );
        test.BasicTest(2,tmp.startEnd.A,"Push modified the start index.",t);
        if i<3 {
            test.BasicTest(i,tmp.vals[i+2],"Push did not save value.",t);
            test.BasicTest(i+2,tmp.startEnd.B,
                "Push did not modify the end index.",t,
            );
        } else {
            test.BasicTest(i,tmp.vals[i-3],"Push did not save value.",t);
            test.BasicTest(i-3,tmp.startEnd.B,
                "Push did not modify the end index.",t,
            );
        }
    }
    res:=tmp.Push(5);
    if !IsQueueFull(res) {
        test.FormatError(QueueFull(""),res,
            "Push did not detect the queue was full.",t,
        );
    }
    test.BasicTest(5,tmp.NumElems,
        "Push incremented NumElems when queue was full.",t,
    );
    test.BasicTest(2,tmp.startEnd.A,"Push modified the start index.",t);
    test.BasicTest(1,tmp.startEnd.B,
        "Push modified the end index when the queue was full.",t,
    );
}

func testCircularQueuePeekHelper(c CircularQueue[int], t *testing.T){
    for i:=0; i<len(c.vals); i++ {
        c.Push(i);
    }
    for i:=0; i<len(c.vals); i++ {
        v,err:=c.Peek(i);
        test.BasicTest(nil,err,
            "Peek returned an error when it should not have.",t,
        );
        test.BasicTest(i,v,"Peek did not return correct value.",t);
    }
    _,err:=c.Peek(5);
    if !customerr.IsValOutsideRange(err) {
        test.FormatError(customerr.ValOutsideRange(""),err,
            "Peek did not return the correct error.",t,
        );
    }
    _,err=c.Peek(-1);
    if !customerr.IsValOutsideRange(err) {
        test.FormatError(customerr.ValOutsideRange(""),err,
            "Peek did not return the correct error.",t,
        );
    }
}
func TestCircularQueuePeek(t *testing.T){
    tmp,err:=NewCircularQueue[int](5);
    test.BasicTest(nil,err,
        "NewCircularQueue returned an error when it should not have.",t,
    );
    testCircularQueuePeekHelper(tmp,t);
    tmp,err=NewCircularQueue[int](5);
    test.BasicTest(nil,err,
        "NewCircularQueue returned an error when it should not have.",t,
    );
    tmp.startEnd.A=2;
    tmp.startEnd.B=1;
    testCircularQueuePeekHelper(tmp,t);
}

func testCircularQueuePeekPntrHelper(c CircularQueue[int], t *testing.T){
    for i:=0; i<len(c.vals); i++ {
        c.Push(i);
    }
    for i:=0; i<len(c.vals); i++ {
        v,err:=c.PeekPntr(i);
        test.BasicTest(nil,err,
            "Peek returned an error when it should not have.",t,
        );
        test.BasicTest(i,*v,"Peek did not return correct value.",t);
        *v=100;
    }
    for i:=0; i<len(c.vals); i++ {
        v,err:=c.PeekPntr(i);
        test.BasicTest(nil,err,
            "Peek returned an error when it should not have.",t,
        );
        test.BasicTest(100,*v,"Peek did not return correct value.",t);
    }
    _,err:=c.Peek(5);
    if !customerr.IsValOutsideRange(err) {
        test.FormatError(customerr.ValOutsideRange(""),err,
            "Peek did not return the correct error.",t,
        );
    }
    _,err=c.Peek(-1);
    if !customerr.IsValOutsideRange(err) {
        test.FormatError(customerr.ValOutsideRange(""),err,
            "Peek did not return the correct error.",t,
        );
    }
}
func TestCircularQueuePeekPntr(t *testing.T){
    tmp,err:=NewCircularQueue[int](5);
    test.BasicTest(nil,err,
        "NewCircularQueue returned an error when it should not have.",t,
    );
    testCircularQueuePeekPntrHelper(tmp,t);
    tmp,err=NewCircularQueue[int](5);
    test.BasicTest(nil,err,
        "NewCircularQueue returned an error when it should not have.",t,
    );
    tmp.startEnd.A=2;
    tmp.startEnd.B=1;
    testCircularQueuePeekPntrHelper(tmp,t);
}

func TestCircularQueuePop(t *testing.T){
    tmp,err:=NewCircularQueue[int](5);
    test.BasicTest(nil,err,
        "NewCircularQueue returned an error when it should not have.",t,
    );
    tmp.Push(0);
    for i:=1; i<21; i++ {
        tmp.Push(i);
        v,err:=tmp.Pop();
        test.BasicTest(nil,err,
            "Pop returned an error when it should not have.",t,
        );
        test.BasicTest(i-1,v,"Pop did not return correct value.",t);
    }
    v,err:=tmp.Pop();
    test.BasicTest(nil,err,
        "Pop returned an error when it should not have.",t,
    );
    test.BasicTest(20,v,"Pop did not return correct value.",t);
    _,err=tmp.Pop();
    if !IsQueueEmpty(err) {
        test.FormatError(QueueEmpty(""),err,
            "Pop did not return correct error.",t,
        );
    }
}

func testCircularQueueElemsHelper(c CircularQueue[int], t *testing.T){
    for i:=0; i<len(c.vals); i++ {
        c.Push(i);
    }
    c.Elems().ForEach(func(index, val int) (iter.IteratorFeedback, error) {
        test.BasicTest(index,val,"Element was skipped while iterating.",t);
        return iter.Continue,nil;
    });
}
func TestCircularQueueElems(t *testing.T){
    tmp,err:=NewCircularQueue[int](5);
    test.BasicTest(nil,err,
        "NewCircularQueue returned an error when it should not have.",t,
    );
    testCircularQueueElemsHelper(tmp,t);
    tmp,err=NewCircularQueue[int](5);
    test.BasicTest(nil,err,
        "NewCircularQueue returned an error when it should not have.",t,
    );
    tmp.startEnd.A=2;
    tmp.startEnd.B=1;
    testCircularQueueElemsHelper(tmp,t);
}

func testCircularQueuePntrElemsHelper(c CircularQueue[int], t *testing.T){
    for i:=0; i<len(c.vals); i++ {
        c.Push(i);
    }
    c.PntrElems().ForEach(func(index int, val *int) (iter.IteratorFeedback, error) {
        test.BasicTest(index,*val,"Element was skipped while iterating.",t);
        *val=100;
        return iter.Continue,nil;
    });
    c.Elems().ForEach(func(index int, val int) (iter.IteratorFeedback, error) {
        test.BasicTest(100,val,"Element was not updated while iterating.",t);
        return iter.Continue,nil;
    });
}
func TestCircularQueuePntrElems(t *testing.T){
    tmp,err:=NewCircularQueue[int](5);
    test.BasicTest(nil,err,
        "NewCircularQueue returned an error when it should not have.",t,
    );
    testCircularQueueElemsHelper(tmp,t);
    tmp,err=NewCircularQueue[int](5);
    test.BasicTest(nil,err,
        "NewCircularQueue returned an error when it should not have.",t,
    );
    tmp.startEnd.A=2;
    tmp.startEnd.B=1;
    testCircularQueueElemsHelper(tmp,t);
}

func queueInterfaceTypeCheck[T any](q types.Queue[T]){}
func TestCircularQueueQueueTypeInterface(t *testing.T) {
    tmp,err:=NewCircularQueue[int](5);
    test.BasicTest(nil,err,
        "NewCircularQueue returned an error when it should not have.",t,
    );
    queueInterfaceTypeCheck[int](&tmp);
}
