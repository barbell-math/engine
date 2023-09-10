package db

import (
	"database/sql"
	"reflect"

	"github.com/barbell-math/block/util/algo"
	"github.com/barbell-math/block/util/algo/iter"
	customReflect "github.com/barbell-math/block/util/reflect"
)

func OnlyIDFilter(col string) bool { return col=="Id"; }
func AllButIDFilter(col string) bool { return col!="Id"; }

func readRows[S any](rows *sql.Rows) iter.Iter[*S] {
    cntr:=0;
    var err error=nil;
    var rowPntrs []reflect.Value=nil;
    return func(f iter.IteratorFeedback) (*S, error, bool) {
        if f==iter.Break {
            rows.Close();
            return nil,nil,false;
        }
        if rows.Next() {
            var s S;
            rowPntrs,err=customReflect.GetStructFieldPntrs(&s,algo.NoFilter[string]);
            if err==nil {
                potErr:=reflect.ValueOf(rows).MethodByName("Scan").Call(rowPntrs);
                err=customReflect.GetErrorFromReflectValue(&potErr[0]);
                if err==nil {
                    cntr++;
                    return &s,nil,true;
                }
            }
        }
        if cntr==0 && err==nil {
            return nil,sql.ErrNoRows,false;
        }
        return nil,err,false;
    }
}
