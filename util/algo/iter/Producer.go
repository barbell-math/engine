package iter;

import (
	"bufio"
	"os"
)

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

func StrElems(s string) Iter[byte] {
    i:=-1;
    return func(f IteratorFeedback) (byte,error,bool) {
        i++;
        if i<len(s) && f!=Break {
            return s[i],nil,true;
        }
        return ' ',nil,false;
    }
}

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

func FileLines(path string) Iter[string] {
    var scanner *bufio.Scanner;
    file,err:=os.Open(path);
    if err==nil {
        scanner=bufio.NewScanner(file);
        scanner.Split(bufio.ScanLines);
    }
    return func(f IteratorFeedback) (string,error,bool) {
        if f==Break || err!=nil || !scanner.Scan() {
            file.Close();
            return "",err,false;
        }
        return scanner.Text(),nil,true;
    }
}
