package io;

import (
    "os"
    "fmt"
    "bytes"
    "bufio"
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

func YNQuestion(question string) bool {
    scanner:=bufio.NewScanner(os.Stdin);
    ans:="";
    cont:=true;
    for cont {
        fmt.Printf("%s? (Y/N) ",question);
        scanner.Scan();
        ans=scanner.Text();
        cont=(len(ans)!=1);
        if !cont {
            cont=(ans[0]!='y' && ans[0]!='Y' && ans[0]!='n' && ans[0]!='N');
        }
    }
    return (ans[0]=='y' || ans[0]=='Y');
}

func Splitter(token string) bufio.SplitFunc {
    temp:=[]byte(token);
    return func (data []byte, atEOF bool) (advance int, token []byte, err error){
        if atEOF && len(data)==0 {
            return 0, nil, nil;
        }
        if i:=bytes.Index(data,temp); i>=0 {
            return i+len(temp), data[0:i], nil;
        }
        if atEOF {
            return len(data),data,nil;
        }
        return 0, nil, nil;
    }
}
