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

func ChainedErrorOpsWithCustomErrors(errs []error, ops ...ErrorOp) error {
    var i int=0;
    var err error=nil;
    res:=make([]any,len(ops));
    for ; err==nil && i<len(ops); i++ {
        res[i],err=ops[i](res[:i]...);
    }
    if err!=nil && i!=len(ops) && i-1<len(errs) {
        return errs[i-1];
    }
    return err;
}

func ErrorOnBool(in bool, e error) error {
    if in {
        return nil;
    }
    return e;
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
    } else if second!=nil {
        first=fmt.Errorf("%s | %s",first,second);
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
