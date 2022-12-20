package db;

import (
    "reflect"
    "database/sql"
    "github.com/barbell-math/block/util/algo"
    customReflect "github.com/barbell-math/block/util/reflect"
)

func OnlyIDFilter(col string) bool { return col=="Id"; }
func AllButIDFilter(col string) bool { return col!="Id"; }

func readRows[S any](rows *sql.Rows, callback func(r *S)) error {
    cntr:=0;
    var err error=nil;
    var rowPntrs []reflect.Value=nil;
    for ; err==nil && rows.Next(); cntr++ {
        var s S;
        rowPntrs,err=customReflect.GetStructFieldPntrs(&s,algo.NoFilter[string]);
        if err==nil {
            potErr:=reflect.ValueOf(rows).MethodByName("Scan").Call(rowPntrs);
            err=customReflect.GetErrorFromReflectValue(&potErr[0]);
            if err==nil {
                callback(&s);
            }
        }
    }
    if cntr==0 && err==nil {
        return sql.ErrNoRows;
    }
    return err;
}
