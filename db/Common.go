package db;

import (
    "reflect"
    "database/sql"
    "github.com/carmichaeljr/powerlifting-engine/util"
)

func OnlyIDFilter(col string) bool { return col=="Id"; }
func AllButIDFilter(col string) bool { return col!="Id"; }

func readRows[S any](rows *sql.Rows, s *S, callback func(r *S)) error {
    rowPntrs,err:=util.GetStructFieldPntrs(s,util.NoFilter[string]);
    for err==nil && rows.Next() {
        potErr:=reflect.ValueOf(rows).MethodByName("Scan").Call(rowPntrs);
        err=util.GetErrorFromReflectValue(&potErr[0]);
        if err==nil {
            callback(s);
        }
    }
    return err;
}
