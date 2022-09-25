package util;

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

