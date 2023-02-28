package algo;

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

func ForEach[T any, U any](i Iter[T], op func(val T)) {
    for next,err,cont:=i(); cont && err==nil; next,err,cont=i() {
        op(next);
    }
}
//reduce
//filter
//filterParallel??

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
