package reflect;

import (
    "fmt"
    stdReflect "reflect"
    "github.com/barbell-math/engine/util/algo"
)

func GetStructName[S any](s *S) (string,error) {
    err:=IsStructVal(s);
    if err==nil {
        val:=stdReflect.ValueOf(s).Elem();
        return val.Type().Name(),err;
    }
    return "",err;
}

func GetStructFieldNames[S any](
        s *S,
        filter algo.Filter[string]) ([]string,error) {
    rv:=make([]string,0);
    err:=IsStructVal(s);
    if err==nil {
        val:=stdReflect.ValueOf(s).Elem();
        for i:=0; i<val.NumField(); i++ {
            colName:=val.Type().Field(i).Name;
            if filter(colName) {
                rv=append(rv,colName);
            }
        }
    }
    return rv,err;
}

func GetStructVals[S any](
        s *S,
        filter algo.Filter[string]) ([]stdReflect.Value,error) {
    rv:=make([]stdReflect.Value,0);
    err:=IsStructVal(s);
    if err==nil {
        val:=stdReflect.ValueOf(s).Elem();
        for i:=0; i<val.NumField(); i++ {
            if filter(val.Type().Field(i).Name) {
                rv=append(rv,stdReflect.ValueOf(val.Field(i).Interface()));
            }
        }
    }
    return rv,err;
}

func GetStructFieldPntrs[S any](
        s *S,
        filter algo.Filter[string]) ([]stdReflect.Value,error) {
    rv:=make([]stdReflect.Value,0);
    err:=IsStructVal(s);
    if err==nil {
        val:=stdReflect.ValueOf(s).Elem();
        for i:=0; i<val.NumField(); i++ {
            valField:=val.Field(i);
            if filter(val.Type().Field(i).Name) {
                rv=append(rv,valField.Addr());
            }
        }
    }
    return rv,err;
}

func IsStructVal[S any](s *S) error {
    if stdReflect.ValueOf(s).Elem().Kind()!=stdReflect.Struct {
        return NonStructValue(fmt.Sprintf(
            "Function requires a struct as target. | Got: %s",
            stdReflect.ValueOf(s).Kind().String(),
        ));
    }
    return nil;
}

func IsStructPntr[S any](s **S) error {
    if stdReflect.ValueOf(s).Elem().Elem().Kind()!=stdReflect.Struct {
        return NonStructValue(fmt.Sprintf(
            "Function requires a struct as target. | Got: %s",
            stdReflect.ValueOf(s).Kind().String(),
        ));
    }
    return nil;
}
