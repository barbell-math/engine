package util;

import "reflect"

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

func GetErrorFromReflectValue(in *reflect.Value) error {
    switch in.Interface().(type) {
        case error: return in.Interface().(error);
        default: return nil;
    }
}

//func safeMapAcc[K comparable, V any](m map[K]V, _var K, err error) (V,error) {
//    var i interface{};
//    if v,ok:=m[_var]; ok {
//        return v,nil;
//    }
//    return V(i),err;
//}
