package db;

import (
    "reflect"
    "database/sql"
    "github.com/carmichaeljr/powerlifting-engine/util"
)

func OnlyIDFilter(col string) bool { return col=="Id"; }
func AllButIDFilter(col string) bool { return col!="Id"; }

func readRows[S any](rows *sql.Rows, callback func(r *S)) error {
    cntr:=0;
    var err error=nil;
    var rowPntrs []reflect.Value=nil;
    for ; err==nil && rows.Next(); cntr++ {
        var s S;
        rowPntrs,err=util.GetStructFieldPntrs(&s,util.NoFilter[string]);
        potErr:=reflect.ValueOf(rows).MethodByName("Scan").Call(rowPntrs);
        err=util.GetErrorFromReflectValue(&potErr[0]);
        if err==nil {
            callback(&s);
        }
    }
    if cntr==0 {
        return sql.ErrNoRows;
    }
    return err;
}
