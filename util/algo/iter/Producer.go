package iter;


func ValElem[T any](val T, err error) Iter[T] {
    cntr:=0;
    return func(f IteratorFeedback) (T,error,bool) {
        var rv T;
        if cntr==0 && f!=Break {
            cntr++;
            return val,err,true;
        }
        return rv,nil,false;    
    }
}

func SliceElems[T any](s []T) Iter[T] {
    i:=-1;
    return func(f IteratorFeedback) (T,error,bool) {
        var rv T;
        i++;
        if i<len(s) && f!=Break {
            return s[i],nil,true;
        }
        return rv,nil,false;
    }
}

//func StrElems(s string) Iter[rune] {
//
//}

//func MapElems[K comparable, V any](s map[K]V) Iter[T] {
//
//}

func ChanElems[T any](c <-chan T) Iter[T] {
    return func(f IteratorFeedback) (T,error,bool) {
        if f!=Break {
            next,ok:=<-c;
            return next,nil,ok;
        }
        var rv T;
        return rv,nil,false;
    }
}
