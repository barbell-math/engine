package iter

import (
	"bufio"
	"os"

	"github.com/barbell-math/engine/util/dataStruct/types"
	customerr "github.com/barbell-math/engine/util/err"
)

func NoElem[T any]() Iter[T] {
    return func(f IteratorFeedback) (T,error,bool) {
        var tmp T;
        return tmp,nil,false;
    }
}

func ValElem[T any](val T, err error, repeat int) Iter[T] {
    cntr:=0;
    return func(f IteratorFeedback) (T,error,bool) {
        var rv T;
        if cntr<repeat && f!=Break {
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

func Join[T any, U any](i1 Iter[T],
        i2 Iter[U],
        v types.Variant[T,U],
        decider func(left T, right U) bool) Iter[types.Variant[T,U]] {
    var i1Val T;
    var i2Val U;
    var err1, err2 error;
    cont1, cont2:=true, true;
    getI1Val, getI2Val:=true, true;
    return func(f IteratorFeedback) (types.Variant[T,U], error, bool) {
        if f==Break {
            return v, customerr.AppendError(i1.Stop(),i2.Stop()), false;
        }
        if getI1Val && cont1 && err1==nil {
            i1Val,err1,cont1=i1(f);
        }
        if getI2Val && cont2 && err2==nil {
            i2Val,err2,cont2=i2(f);
        }
        if err1==nil && err2==nil {
            if cont1 && cont2 {
                d:=decider(i1Val,i2Val);
                getI1Val=d;
                getI2Val=!d;
                if d {
                    return v.SetValA(i1Val),err1,cont1 && cont2;
                } else {
                    return v.SetValB(i2Val),err2,cont1 && cont2;
                }
            } else if cont1 && !cont2 {
                getI1Val=true;
                getI2Val=false;
                return v.SetValA(i1Val),err1,cont1;
            } else if !cont1 && cont2 {
                getI1Val=false;
                getI2Val=true;
                return v.SetValB(i2Val),err2,cont2;
            }
        }
        return v,customerr.AppendError(err1,err2),false;
    }
}

func JoinSame[T any](i1 Iter[T],
        i2 Iter[T],
        v types.Variant[T,T],
        decider func(left T, right T) bool) Iter[T] {
    var tmp T;
    realJoiner:=Join(i1,i2,v,decider);
    return func(f IteratorFeedback) (T, error, bool) {
        val,err,cont:=realJoiner(f);
        if cont && err==nil {
            if val.HasA() {
                return val.ValA(),err,cont;
            } else if val.HasB() {
                return val.ValB(),err,cont;
            }
        }
        return tmp,err,cont;
    }
}
