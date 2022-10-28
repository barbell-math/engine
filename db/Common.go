package db;

import (
    "fmt"
    "reflect"
    "database/sql"
    "github.com/carmichaeljr/powerlifting-engine/util"
)

func OnlyIDFilter(col string) bool { return col=="Id"; }
func AllButIDFilter(col string) bool { return col!="Id"; }

func readRows[S any](rows *sql.Rows, callback func(r *S)) error {
    var err error=nil;
    var rowPntrs []reflect.Value=nil;
    var s S;
    for err==nil && rows.Next() {
        fmt.Print("Before: %v  ",s);
        var s S;
        fmt.Println("Adter: %v",s);
        rowPntrs,err=util.GetStructFieldPntrs(&s,util.NoFilter[string]);
        potErr:=reflect.ValueOf(rows).MethodByName("Scan").Call(rowPntrs);
        err=util.GetErrorFromReflectValue(&potErr[0]);
        if err==nil {
            callback(&s);
        }
    }
    //for err==nil && rows.Next() {
    //    s:=new(S);
    //    rowPntrs,err:=util.GetStructFieldPntrs(s,util.NoFilter[string]);
    //    potErr:=reflect.ValueOf(rows).MethodByName("Scan").Call(rowPntrs);
    //    err=util.GetErrorFromReflectValue(&potErr[0]);
    //    if err==nil {
    //        callback(s);
    //    }
    //}
    return err;
}
