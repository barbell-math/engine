package util;

import (
    "fmt"
    "reflect"
)

func GetStructName[S any](s *S) (string,error) {
    err:=checkForStructVal(s);
    if err==nil {
        val:=reflect.ValueOf(s).Elem();
        return val.Type().Name(),err;
    }
    return "",err;
}

func GetStructFieldNames[S any](
        s *S,
        filter Filter[string]) ([]string,error) {
    rv:=make([]string,0);
    err:=checkForStructVal(s);
    if err==nil {
        val:=reflect.ValueOf(s).Elem();
        for i:=0; i<val.NumField(); i++ {
            colName:=val.Type().Field(i).Name;
            if filter(colName) {
                rv=append(rv,colName);
            }
        }
    }
    return rv,err;
}

func GetStructVals[S any](s *S, filter Filter[string]) ([]reflect.Value,error) {
    rv:=make([]reflect.Value,0);
    err:=checkForStructVal(s);
    if err==nil {
        val:=reflect.ValueOf(s).Elem();
        for i:=0; i<val.NumField(); i++ {
            if filter(val.Type().Field(i).Name) {
                rv=append(rv,reflect.ValueOf(val.Field(i).Interface()));
            }
        }
    }
    return rv,err;
}

func GetStructFieldPntrs[S any](
        s *S,
        filter Filter[string]) ([]reflect.Value,error) {
    rv:=make([]reflect.Value,0);
    err:=checkForStructVal(s);
    if err==nil {
        val:=reflect.ValueOf(s).Elem();
        for i:=0; i<val.NumField(); i++ {
            valField:=val.Field(i);
            if filter(val.Type().Field(i).Name) {
                rv=append(rv,valField.Addr());
            }
        }
    }
    return rv,err;
}

func checkForStructVal[S any](s *S) error {
    if reflect.ValueOf(s).Elem().Kind()!=reflect.Struct {
        return NonStructValue(fmt.Sprintf(
            "Function requires a struct as target. | Got: %s",
            reflect.ValueOf(s).Kind().String(),
        ));
    }
    return nil;
}
