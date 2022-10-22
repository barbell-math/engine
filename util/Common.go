package util;

import (
    "os"
    "fmt"
    "bytes"
    "bufio"
    "errors"
)

func Splitter(token string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
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

func AppendWithPreallocation[T any](slices ...[]T) []T {
    var totLen,i int=0, 0;
    for _,s:=range(slices) {
        totLen+=len(s);
    }
    rv:=make([]T,totLen);
    for _,tmp:=range(slices){
        i+=copy(rv[i:],tmp);
    }
    return rv;
}

func ZipSlices[K comparable, V any](keys []K, vals []V) (map[K]V,error) {
    rv:=make(map[K]V,len(keys));
    if len(keys)==len(vals) {
        for i,k:=range(keys) {
            if _,ok:=rv[k]; !ok {
                rv[k]=vals[i];
            } else {
                return rv,SliceZippingError(fmt.Sprintf(
                    "Keys have duplicate values | %v",k,
                ));
            }
        }
    } else {
        return rv,SliceZippingError(fmt.Sprintf(
            "Lengths are not equal. | K: %d V: %d",len(keys),len(vals),
        ));
    }
    return rv,nil;
}

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
