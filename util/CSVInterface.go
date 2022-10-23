package util;

import (
    "os"
    "fmt"
    "time"
    "bufio"
    "strings"
    "reflect"
    "strconv"
)

//Only basic types are supported through the CSV interface (which mirrors the
//limitations of a CSV file). For things like arrays and maps JSON is far more
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
//Note that quotes around strings are optional, and will be removed before setting
//the value of the struct variable. Also note that only double quotes (") are
//recognized, single quotes (') are not.
//If any columns are missing or there are blank values the corresponding values
//in the structs that are generated will be zero-value initialized.
func CSVToStruct[R any](
        src string,
        delim byte,
        timeDateFormat string,
        callback func(tab *R)) error {
    cntr:=1;
    var iter R;
    var err error=nil;
    headers:=make([]string,0);
    return ChainedErrorOps(
        func(r ...any) (any,error) {
            return nil,ErrorOnBool(reflect.ValueOf(iter).Kind()==reflect.Struct,
                NonStructValue(fmt.Sprintf(
                    "CSVToStruct requires a struct as target. | Got: %s",
                    reflect.ValueOf(iter).Kind().String(),
                )),
            );
        },func(r ...any) (any,error) {
            if e1:=CSVFileSplitter(src,delim,false,func(c []string) bool {
                if cntr!=1 {
                    if err=convFromCSV[R](&iter,headers,c,timeDateFormat); err==nil {
                        callback(&iter);
                    } else {
                        err=MalformedCSVFile(
                            fmt.Sprintf("File: %s: Line %d: %s",src,cntr,err),
                        );
                    }
                } else {
                    headers=c;
                }
                cntr++;
                return err==nil;
            }); e1!=nil {
                return nil,e1;
            } else {
                return nil,err;
            }
    });
}

func convFromCSV[R any](
        v *R,
        headers []string,
        columns []string,
        timeDateFormat string) error {
    tmp:=fmt.Errorf(fmt.Sprintf(
        "Expected %d cols, have %d",len(headers),len(columns),
    ));
    return ChainedErrorOps(
        func(r ...any) (any,error) {
            return nil,ErrorOnBool(len(headers)==len(columns),tmp);
        }, func(r ...any) (any,error) {
            var err error=nil;
            for i:=0; err==nil && i<len(headers); i++ {
                if len(columns[i])>0 {
                    err=setTableValue[R](v,headers[i],columns[i],timeDateFormat);
                }
            }
            return nil,err;
        },
    );
}

//Only basic types are supported
func setTableValue[R any](
        row *R,
        name string,
        val string,
        timeDateFormat string) error {
    var err error=nil;
    s:=reflect.ValueOf(row).Elem();
    f:=s.FieldByName(name);
    if f.IsValid() && f.CanSet() {
        switch f.Interface().(type) {
            case time.Time: var tmp time.Time;
                tmp,err=time.Parse(timeDateFormat,val);
                f.Set(reflect.ValueOf(tmp));
            case bool: var tmp bool;
                tmp,err=strconv.ParseBool(val);
                f.SetBool(tmp);
            case uint: err=setUint[uint](f,val);
            case uint8: err=setUint[uint8](f,val);
            case uint16: err=setUint[uint16](f,val);
            case uint32: err=setUint[uint32](f,val);
            case uint64: err=setUint[uint64](f,val);
            case int: err=setInt[int](f,val);
            case int8: err=setInt[int8](f,val);
            case int16: err=setInt[int16](f,val);
            case int32: err=setInt[int32](f,val);
            case int64: err=setInt[int64](f,val);
            case float32: err=setFloat[float32](f,val);
            case float64: err=setFloat[float32](f,val);
            case string: err=setString(f,val);
            default: err=fmt.Errorf(
                "The type '%s' is not a supported type.",f.Kind().String(),
            );
        }
    } else {
        err=fmt.Errorf(
            "Requested header value not in struct or is not settable. | Header: '%s'",
            name,
        );
    }
    return err;
}

func setUint[N ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](
        f reflect.Value,
        v string) error {
    tmp,err:=strconv.ParseUint(v,10,64);
    f.SetUint(tmp);
    return err;
}
func setInt[N ~int | ~int8 | ~int16 | ~int32 | ~int64](
        f reflect.Value,
        v string) error {
    tmp,err:=strconv.ParseInt(v,10,64);
    f.SetInt(tmp);
    return err;
}
func setFloat[N ~float32 | ~float64](f reflect.Value, v string) error {
    tmp,err:=strconv.ParseFloat(v,64);
    f.SetFloat(tmp);
    return err;
}
func setString(f reflect.Value, v string) error {
    s,e:=0, len(v);
    if len(v)>0 && v[0]=='"' {
        s++;
    }
    if len(v)>0 && v[len(v)-1]=='"' {
        e--;
    }
    f.SetString(v[s:e]);
    return nil;
}

func CSVFileSplitter(
        src string,
        delim byte,
        hasHeaders bool,
        rowCallback func(columns []string) bool) error {
    if f,err:=os.Open(src); err==nil {
        defer f.Close();
        scanner:=bufio.NewScanner(f);
        if hasHeaders {
            scanner.Scan();
        }
        cont:=true;
        for cont && scanner.Scan() {
            cont=splitLineForCSV(scanner.Text(),delim,rowCallback);
        }
        return nil;
    } else {
        return err;
    }
}

func splitLineForCSV(
        row string,
        delim byte,
        rowCallback func(columns[]string) bool) bool {
    prevIndex:=0;
    inQuotes:=false;
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
    return rowCallback(columns);
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
