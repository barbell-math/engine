package db;

import (
    "fmt"
    "github.com/carmichaeljr/powerlifting-engine/util"
)

//Only basic types are supported through the CSV interface (which mirrors what
//an SQL DB can contain). For things like slices and maps JSON is far more
//expressive and less error prone.
//Supported types:
//  - All integer types (int,int8,int16,int32,int64)
//  - All unsigned integer types (uint,uint8,uint16,uint32,uint64)
//  - All float types (float32,float64)
//  - Strings
//  - Booleans
//  - TimeDate formats
//The CSV file **MUST** have headers. Without this the structs fields cannot
//be set properly.
func CSVToDBTable[R DBTable](
        src string,
        delim byte,
        timeDateFormat string,
        callback func(tab *R)) error {
    cntr:=0;
    var iter R;
    var err error=nil;
    headers:=make([]string,0);
    util.CSVFileSplitter(src,delim,false,func(c []string) bool {
        if cntr!=0 {
            if err=convFromCSV[R](&iter,headers,c,timeDateFormat); err==nil {
                callback(&iter);
            } else {
                err=util.MalformedCSVToDBTableFile(
                    fmt.Sprintf("Line %d: %s",cntr,err),
                );
            }
        } else {
            headers=c;
        }
        cntr++;
        return err==nil;
    });
    return err;
}

func convFromCSV[R DBTable](
        v *R,
        headers []string,
        columns []string,
        timeDateFormat string) error {
    tmp:=util.MalformedCSVToDBTableFile(fmt.Sprintf(
        "Expected %d cols, have %d",len(headers),len(columns),
    ));
    return util.ChainedErrorOps(
        func(r ...any) (any,error) {
            return nil,util.ErrorOnBool(len(headers)==len(columns),tmp);
        }, func(r ...any) (any,error) {
            var err error=nil;
            for i:=0; err==nil && i<len(headers); i++ {
                err=setTableValue[R](v,headers[i],columns[i],timeDateFormat);
            }
            return nil,err;
        },
    );
}
