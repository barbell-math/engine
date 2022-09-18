package util;

import (
    "os"
    "fmt"
    "bytes"
    "bufio"
    "strings"
    "reflect"
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

func CSVFileSplitter(
        src string,
        delim byte,
        hasHeaders bool,
        rowCallback func(columns []string) bool){
    if f,err:=os.Open(src); err==nil {
        defer f.Close();
        scanner:=bufio.NewScanner(f);
        if hasHeaders {
            scanner.Scan();
        }
        cont:=true;
        for cont && scanner.Scan() {
            prevIndex:=0;
            inQuotes:=false;
            row:=scanner.Text();
            columns:=make([]string,0);
            for i,char:=range(row) {
                if char=='"' {
                    inQuotes=!inQuotes;
                } else if byte(char)==delim && !inQuotes {
                    columns=append(columns,strings.TrimSpace(row[prevIndex:i]));
                    prevIndex=i+1;
                }
            }
            columns=append(columns,strings.TrimSpace(row[prevIndex:]));
            cont=rowCallback(columns);
        }
    }
}

func CSVGenerator(sep string, callback func(iter int) (string,bool)) string {
    var sb strings.Builder;
    var temp string;
    cont:=true;
    for i:=0; cont; i++ {
        temp,cont=callback(i);
        sb.WriteString(temp);
        if cont {
            sb.WriteString(sep);
        }
    }
    return sb.String();
}

func GetErrorFromReflectValue(in *reflect.Value) error {
    switch in.Interface().(type) {
        case error: return in.Interface().(error);
        default: return nil;
    }
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
