package db;

import (
    "fmt"
    "reflect"
    "database/sql"
    "github.com/carmichaeljr/powerlifting-engine/util"
)

type FilterType func(col string) bool;

func NoFilter(col string) bool { return true; }
func OnlyIDFilder(col string) bool { return col=="Id"; }
func AllButIDFilter(col string) bool { return col!="Id"; }

func Create[R DBTable](c *CRUD, row R) (int,error) {
    var rv int=0;
    columns:=getTableColumns(&row,AllButIDFilter);
    if len(columns)==0 {
        return -1,util.FilterRemovedAllColumns("Row was not added to database.");
    }
    intoStr:=util.CSVGenerator(",",func(iter int) (string,bool) {
        return columns[iter], iter+1<len(columns);
    });
    valuesStr:=util.CSVGenerator(",",func(iter int) (string,bool) {
        return fmt.Sprintf("$%d",iter+1), iter+1<len(columns);
    });
    sqlStmt:=fmt.Sprintf(
        "INSERT INTO %s(%s) VALUES (%s) RETURNING Id;",
        getTableName(&row),intoStr,valuesStr,
    );
    vals:=append(
        []reflect.Value{reflect.ValueOf(sqlStmt)},
        getTableVals(&row,AllButIDFilter)...
    );
    rowValue:=reflect.ValueOf(c.db).MethodByName("QueryRow").Call(vals)[0];
    err:=rowValue.Interface().(*sql.Row).Scan(&rv);
    return rv,err;
}

func Read[R DBTable](
        c *CRUD,
        rowVals R,
        filter FilterType,
        callback func(val *R)) error {
    columns:=getTableColumns(&rowVals,filter);
    if len(columns)==0 {
        return util.FilterRemovedAllColumns("No value rows were selected.");
    }
    valuesStr:=util.CSVGenerator(" AND ",func(iter int) (string,bool) {
        return fmt.Sprintf("%s=$%d",columns[iter],iter+1), iter+1<len(columns);
    });
    sqlStmt:=fmt.Sprintf(
        "SELECT * FROM %s WHERE %s;",getTableName(&rowVals),valuesStr,
    );
    vals:=append(
        []reflect.Value{reflect.ValueOf(sqlStmt)},
        getTableVals(&rowVals,filter)...
    );
    callVals:=reflect.ValueOf(c.db).MethodByName("Query").Call(vals);
    rows,err:=callVals[0].Interface().(*sql.Rows),
        util.GetErrorFromReflectValue(&callVals[1]);
    defer rows.Close();
    var iter R;
    rowPntrs:=getTablePntrs(&iter,NoFilter);
    for err==nil && rows.Next() {
        potErr:=reflect.ValueOf(rows).MethodByName("Scan").Call(rowPntrs);
        err=util.GetErrorFromReflectValue(&potErr[0]);
        callback(&iter);
    }
    return err;
}

func Update[R DBTable](
        c *CRUD,
        rowVals R,
        rowValsFilter FilterType,
        updateVals R,
        updateValsFilter FilterType,
        callback func(val *R)) (int,error) {
    return 0 ,nil;
}

func getTableName[R DBTable](row *R) string {
    val:=reflect.ValueOf(row).Elem();
    return val.Type().Name();
}

func getTableColumns[R DBTable](row *R, filter func(col string) bool) []string {
    val:=reflect.ValueOf(row).Elem();
    rv:=make([]string,0);
    for i:=0; i<val.NumField(); i++ {
        colName:=val.Type().Field(i).Name;
        if filter(colName) {
            rv=append(rv,colName);
        }
    }
    return rv;
}

func getTableVals[R DBTable](row *R, filter func(col string) bool) []reflect.Value {
    val:=reflect.ValueOf(row).Elem();
    rv:=make([]reflect.Value,0);
    for i:=0; i<val.NumField(); i++ {
        if filter(val.Type().Field(i).Name) {
            rv=append(rv,reflect.ValueOf(val.Field(i).Interface()));
        }
    }
    return rv;
}

func getTablePntrs[R DBTable](row *R,filter func(col string) bool) []reflect.Value {
    val:=reflect.ValueOf(row).Elem();
    rv:=make([]reflect.Value,0);
    for i:=0; i<val.NumField(); i++ {
        valField:=val.Field(i);
        if filter(val.Type().Field(i).Name) {
            rv=append(rv,valField.Addr());
        }
    }
    return rv;
}
