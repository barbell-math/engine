package algo

type Iter[T any] func()(T,error,bool);

func ValElem[T any](val T, err error) Iter[T] {
    cntr:=0;
    return func() (T,error,bool) {
        var rv T;
        if cntr==0 {
            cntr++;
            return val,err,true;
        }
        return rv,nil,false;    
    }
}

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

func (i Iter[T])ForEach(op func(index int, val T) error) error {
    var rv error;
    j:=0;
    for next,err,cont:=i(); cont && err==nil && rv==nil; next,err,cont=i() {
        rv=op(j,next);
        j++;
    }
    return rv;
}

func (i Iter[T])Filter(op Filter[T]) Iter[T] {
    return func() (T,error,bool) {
        var val T;
        var err error;
        var cont bool;
        for val,err,cont=i(); err==nil && cont && !op(val); val,err,cont=i() {}
        return val,err,cont;
    }
}

func (i Iter[T])Reduce(op func(accum *T, iter T) error, start T) (T,error) {
    accum:=start;
    var rvErr error;
    for next,err,cont:=i(); cont && err==nil && rvErr==nil; next,err,cont=i() {
        rvErr=op(&accum,next);
    }
    return accum,rvErr;
}


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

func (i Iter[T])AppendTo(orig *[]T) int {
    j:=0;
    for val,err,cont:=i(); cont && err==nil; val,err,cont=i() {
        *orig=append(*orig,val);
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

//ToChan
