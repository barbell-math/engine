package algo

import (
	"fmt"
	"github.com/barbell-math/block/util/dataStruct/base"
	customerr "github.com/barbell-math/block/util/err"
)

type Iter[T any] func()(T,error,bool);

func SliceElems[T any](s []T) Iter[T] {
    i:=-1;
    return func() (T,error,bool) {
        var rv T;
        i++;
        if i<len(s) {
            return s[i],nil,true;
        }
        return rv,nil,false;
    }
}

//func MapElems[K comparable, V any](s map[K]V) Iter[T] {
//
//}

func ChanElems[T any](c <-chan T) Iter[T] {
    return func() (T,error,bool) {
        next,ok:=<-c;
        return next,nil,ok;
    }
}

func Map[T any, U any](i Iter[T], op func(val T) (U,error)) Iter[U] {
    return func() (U,error,bool) {
        var tmp U;
        next,err,ok:=i();
        if !ok || err!=nil {
            return tmp,err,false;
        }
        mappedVal,err:=op(next);
        return mappedVal,err,true;
    }
}
func (i Iter[T])Map(op func(val T) (T,error)) Iter[T] {
    return Map(i,op);
}

func (i Iter[T])ForEach(op func(val T) error) error {
    var rv error;
    for next,err,cont:=i(); cont && err==nil && rv==nil; next,err,cont=i() {
        rv=op(next);
    }
    return rv;
}

func threadRunner[T any](jobs chan T, errs chan error, op func(val T) error){
    for val:=range(jobs) {
        errs <- op(val);
    }
}
func (i Iter[T])ForEachParallel(op func(val T) error, numThreads int) []error {
    if numThreads<1 {
        return []error{customerr.ValOutsideRange(fmt.Sprintf(
            "Expected >0 | Got: %d",numThreads,
        ))};
    }
    j,taken:=0,0;
    rv:=make([]error,0);
    jobs, errs:=make(chan T), make(chan error);
    for next,err,cont:=i(); cont && err==nil; next,err,cont=i() {
        fmt.Println("Next: ",next);
        fmt.Println("err: ",err);
        fmt.Println("Cont: ",cont);
        if j<numThreads {
            go threadRunner(jobs,errs,op);
        }
        fmt.Println("Writing to jobs...");
        //if all the threads are working, wait to add a value
        if j>=numThreads {
            if tmp,ok:= <- errs; ok {
                taken++;
                if tmp!=nil {
                    rv=append(rv,tmp);
                }
            }
        }
        jobs <- next;
        j++;
    }
    close(jobs);
    //need to wait for all results!!
    fmt.Println("len: ",taken);
    fmt.Println("Waiting on ",j-taken," threads");
    rv=append(rv,ChanElems(errs).Take(j-taken).
    Filter(func(val base.Pair[error,error]) bool {
        return val.First!=nil;
    }).Collect()...);
    close(errs);
    return rv;
}

//reduce

func (i Iter[T])Filter(op Filter[base.Pair[T,error]]) Iter[T] {
    return func() (T,error,bool) {
        var val T;
        var err error;
        var cont bool;
        for val,err,cont=i(); err==nil && cont && !op(base.Pair[T,error]{
            First: val, Second: err,
        }); val,err,cont=i() {}
        return val,err,cont;
    }
}
//filterParallel??

func (i Iter[T])Take(num int) Iter[T] {
    cntr:=0;
    return func() (T,error,bool) {
        var tmp T;
        if cntr<num {
            cntr++;
            return i();
        }
        return tmp,nil,false;
    }
}

func (i Iter[T])TakeWhile(op func(val T, err error) bool) Iter[T] {
    stop:=false;
    return func() (T,error,bool) {
        var tmp T;
        if !stop {
            val,err,cont:=i();
            if stop=!op(val,err); !stop {
                return val,err,cont;
            }
        }
        return tmp,nil,false;
    }
}

func (i Iter[T])Consume() {
    for _,_,cont:=i(); cont; _,_,cont=i() {}
}

func (i Iter[T])Collect() []T {
    rv:=make([]T,0);
    for val,err,cont:=i(); cont && err==nil; val,err,cont=i() {
        rv=append(rv,val);
    }
    return rv;
}

func (i Iter[T])CollectInto(buffer []T) int {
    j:=0;
    for val,err,cont:=i(); cont && err==nil && j<len(buffer); val,err,cont=i() {
        buffer[j]=val;
        j++;
    }
    return j;
}

func (i Iter[T])All(op func(val T) bool) bool {
    for val,err,cont:=i(); cont && err==nil; val,err,cont=i() {
        if !op(val) {
            return false;
        }
    }
    return true;
}

func (i Iter[T])Any(op func(val T) bool) bool {
    for val,err,cont:=i(); cont && err==nil; val,err,cont=i() {
        if op(val) {
            return true;
        }
    }
    return false;
}

func (i Iter[T])Count() int {
    rv:=0;
    for _,err,cont:=i(); cont && err==nil; _,err,cont=i() {
        rv++;
    }
    return rv;
}

func (i Iter[T])Find(op func(val T) bool) (T,bool) {
    var tmp T;
    for val,err,cont:=i(); cont && err==nil; val,err,cont=i() {
        if op(val) {
            return val,true;
        }
    }
    return tmp,false;
}

func (i Iter[T])Index(op func(val T) bool) int {
    j:=0;
    for val,err,cont:=i(); cont && err==nil; val,err,cont=i() {
        if op(val) {
            return j;
        }
        j++;
    }
    return -1;
}

func (i Iter[T])Nth(idx int) (T,bool) {
    j:=0;
    var tmp T;
    for val,err,cont:=i(); cont && err==nil; val,err,cont=i() {
        if j==idx {
            return val,true;
        }
        j++;
    }
    return tmp,false;
}
