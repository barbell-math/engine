package err;

import (
    "fmt"
    "errors"
    "strings"
)

type errorType func(addendum string) error;
type isErrorType func(err error) bool;
func ErrorFactory(base string) (errorType,isErrorType) {
    errorTypeRv:=func(addendum string) error {
        return errors.New(fmt.Sprintf("%s | %s",base,addendum));
    }
    isErrorTypeRv:=func(err error) bool {
        return strings.Contains(fmt.Sprint(err),base);
    }
    return errorTypeRv,isErrorTypeRv;
}

type ErrorOp func(results ...any) (any,error);
func ChainedErrorOps(ops ...ErrorOp) error {
    var err error=nil;
    res:=make([]any,len(ops));
    for i:=0; err==nil && i<len(ops); i++ {
        res[i],err=ops[i](res[:i]...);
    }
    return err;
}

func ErrorOnBool(in bool, e error) error {
    if in {
        return nil;
    }
    return e;
}

func PanicOnError(op func() error){
    err:=op();
    if err!=nil {
        panic(err);
    }
}

func ArrayDimsArgree[N any, P any](one []N, given []P, message string) error {
    if lOne,lTwo:=len(one),len(given); lOne!=lTwo {
        return DimensionsDoNotAgree(fmt.Sprintf(
            "%s | len(one)=%d len(two)=%d",
            message,lOne,lTwo,
        ));
    }
    return nil;
}

func AppendError(first error, second error) error {
    if second!=nil && first==nil {
        first=second;
    } else if second!=nil && first!=nil {
        first=fmt.Errorf(
            "%s \nThe following error was also generated: %s",
            first,second,
        );
    }
    return first;
}

//func safeMapAcc[K comparable, V any](m map[K]V, _var K, err error) (V,error) {
//    var i interface{};
//    if v,ok:=m[_var]; ok {
//        return v,nil;
//    }
//    return V(i),err;
//}
