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

func (i Iter[T])

//func Map[T any, U any](i Iter[T], op func(val T) (U,error)) Iter[U] {
//    return func() (U,error,bool) {
//        var tmp U;
//        for {
//            next,err,ok:=i();
//            if !ok || err!=nil {
//                return tmp,err,false;
//            }
//            mappedVal,err:=op(next);
//            return mappedVal,err,true;
//        }
//    }
//}

func (i Iter[T])Collect() []T {
    rv:=make([]T,0);
    _break:=false;
    for !_break {
        val,_,cont:=i();
        _break=!cont;
        rv=append(rv,val);
    }
    return rv;
}
