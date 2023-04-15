package dataStruct

import (
	"fmt"
	"github.com/barbell-math/block/util/algo/iter"
	customerr "github.com/barbell-math/block/util/err"
)

type CircularQueue[T any] struct {
    vals []T;
    NumElems int;
    startEnd Pair[int,int];
};

func NewCircularQueue[T any](size int) (CircularQueue[T],error) {
    if size<=0 {
        return CircularQueue[T]{},customerr.ValOutsideRange(
            fmt.Sprintf("Size of queue must be >=0 | Have: %d",size),
        );
    }
    return CircularQueue[T]{
        vals: make([]T,size),
        startEnd: Pair[int, int]{First: 0, Second: -1},
    },nil;
}

func (c *CircularQueue[T])Capacity() int {
    return len(c.vals);
}

func (c *CircularQueue[T])Push(v T) error {
    if c.NumElems<len(c.vals) {
        c.NumElems++;
        c.startEnd.Second=c.wrapAroundIndex(c.startEnd.Second);
        c.vals[c.startEnd.Second]=v;
        return nil;
    }
    return QueueFull(fmt.Sprintf("Queue size: %d",len(c.vals)));
}

func (c *CircularQueue[T])Peek(idx int) (T,error) {
    v,err:=c.PeekPntr(idx);
    if v!=nil {
        return *v,err;
    }
    var tmp T;
    return tmp,err;
}

func (c *CircularQueue[T])PeekPntr(idx int) (*T,error) {
    if idx>=0 && idx<c.NumElems {
        properIndex:=idx+c.startEnd.First;
        if properIndex>=len(c.vals) {
            properIndex-=len(c.vals);
        }
        return &c.vals[properIndex],nil;
    }
    return nil,customerr.ValOutsideRange(
        fmt.Sprintf(
            "Index must be >0 and <num elems in queue. | Num Elems: %d Index: %d",
            c.NumElems,idx,
    ));
}

func (c *CircularQueue[T])Pop() (T,error) {
    if c.NumElems>0 {
        rv:=c.vals[c.startEnd.First];
        c.startEnd.First=c.wrapAroundIndex(c.startEnd.First);
        c.NumElems--;
        return rv,nil;
    }
    var tmp T;
    return tmp,QueueEmpty("Nothing to pop!");
}

func (c *CircularQueue[T])Elems() iter.Iter[T] {
    i:=-1;
    return func(f iter.IteratorFeedback) (T,error,bool) {
        var rv T;
        i++;
        if i<c.NumElems && f!=iter.Break {
            v,err:=c.Peek(i);
            return v,err,true;
        }
        return rv,nil,false;
    }
}

func (c *CircularQueue[T])PntrElems() iter.Iter[*T] {
    i:=-1;
    return func(f iter.IteratorFeedback) (*T,error,bool) {
        i++;
        if i<c.NumElems && f!=iter.Break {
            v,err:=c.PeekPntr(i);
            return v,err,true;
        }
        return nil,nil,false;
    }
}

func (c *CircularQueue[T])wrapAroundIndex(val int) int {
    if val+1>=len(c.vals) {
        return 0;
    }
    return val+1;
}
