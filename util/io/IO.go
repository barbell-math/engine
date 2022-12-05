package io;

import (
    "os"
    "errors"
)

func FileExists(f string) (bool,error) {
    info,err:=os.Stat(f);
    if err==nil {
        return !info.IsDir(),nil;
    }
    if errors.Is(err,os.ErrNotExist) {
        return false,nil;
    }
    return false,err;
}
