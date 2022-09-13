package util;

import (
    "os"
    "fmt"
    "bytes"
    "bufio"
    "strings"
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

func CSVFileSplitter(src string, delim byte, hasHeaders bool, rowCallback func(columns []string)){
    if f,err:=os.Open(src); err==nil {
        defer f.Close();
        scanner:=bufio.NewScanner(f);
        if hasHeaders {
            scanner.Scan();
        }
        for scanner.Scan() {
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
            rowCallback(columns);
        }
    }
}