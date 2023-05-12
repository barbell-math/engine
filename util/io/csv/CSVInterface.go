package csv;

import (
	"io"
	"os"
	"fmt"
    "strings"
	"strconv"
	"reflect"
	"time"
	"encoding/csv"
	"github.com/barbell-math/block/util/algo/iter"
)

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

func CSVFileSplitter(src string, delim rune, comment rune) iter.Iter[[]string] {
    var reader *csv.Reader=nil;
    file,err:=os.Open(src);
    if err==nil {
        reader=csv.NewReader(file);
        reader.Comma=delim;
        reader.Comment=comment;
    }
    return func(f iter.IteratorFeedback) ([]string, error, bool) {
        if f==iter.Break || err!=nil {
            file.Close();
            return []string{},err,false;
        }
        cols,readerErr:=reader.Read();
        if readerErr!=nil {
            if readerErr==io.EOF {
                return cols,nil,false;
            }
            return []string{},readerErr,false;
        }
        return cols,readerErr,true;
    }
}

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
//If any columns are missing or there are blank values the corresponding values
//in the structs that are generated will be zero-value initialized.
func CSVToStruct[R any](src iter.Iter[[]string], timeDateFormat string) iter.Iter[R] {
    var tmp R;
    if reflect.ValueOf(tmp).Kind()!=reflect.Struct {
        return iter.ValElem(tmp,NonStructValue(fmt.Sprintf(
            "CSVToStruct requires a struct as target. | Got: %s",
            reflect.ValueOf(tmp).Kind().String(),
        )),1);
    }
    headers:=make([]string,0);
    return iter.Next(src,
    func(index int, val []string, status iter.IteratorFeedback) (iter.IteratorFeedback, R, error) {
        if status==iter.Break {
            return iter.Break,tmp,nil;
        }
        if index==0 {
            headers=val;
            return iter.Iterate,tmp,nil;
        } 
        if i,err:=convFromCSV[R](headers,val,timeDateFormat); err==nil {
            return iter.Continue,i,err;
        } else {
            err=MalformedCSVFile(
                //fmt.Sprintf("File '%s': Line %d: %s",src,index+1,err),
                fmt.Sprintf("Line %d: %s",index+1,err),
            );
            return iter.Continue,i,err;
        }
    });
}

func convFromCSV[R any](
        headers []string,
        columns []string,
        timeDateFormat string) (R,error) {
    var rv R;
    if len(headers)==len(columns) {
        var err error=nil;
        for i:=0; err==nil && i<len(headers); i++ {
            if len(columns[i])>0 {
                err=setTableValue(&rv,headers[i],columns[i],timeDateFormat);
            }
        }
        return rv,err;
    }
    return rv,fmt.Errorf(fmt.Sprintf(
        "Expected %d cols, have %d",len(headers),len(columns),
    ));
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
