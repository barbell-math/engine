package db;

import (
    "fmt"
    "reflect"
    "database/sql"
    "github.com/carmichaeljr/powerlifting-engine/util"
)

type ColumnFilter func(col string) bool;

func NoFilter(col string) bool { return true; }
func OnlyIDFilder(col string) bool { return col=="Id"; }
func AllButIDFilter(col string) bool { return col!="Id"; }

func GenColFilter(inverse bool, cols ...string) ColumnFilter {
    return func(col string) bool {
        rv:=inverse;
        for i:=0; (inverse && rv) || (!inverse && !rv) && i<len(cols); i++ {
            if inverse {
                rv=(rv && col!=cols[i]);
            } else {
                rv=(col==cols[i]);
            }
        }
        return rv;
    }
}

func Create[R DBTable](c *CRUD, rows ...R) ([]int,error) {
    if len(rows)==0 {
        return []int{},sql.ErrNoRows;
    }
    columns:=getTableColumns(&rows[0],AllButIDFilter);
    if len(columns)==0 {
        return []int{},util.FilterRemovedAllColumns("Row was not added to database.");
    }
    intoStr:=util.CSVGenerator(",",func(iter int) (string,bool) {
        return columns[iter], iter+1<len(columns);
    });
    valuesStr:=util.CSVGenerator(",",func(iter int) (string,bool) {
        return fmt.Sprintf("$%d",iter+1), iter+1<len(columns);
    });
    sqlStmt:=fmt.Sprintf(
        "INSERT INTO %s(%s) VALUES (%s) RETURNING Id;",
        getTableName(&rows[0]),intoStr,valuesStr,
    );
    var err error=nil;
    rv:=make([]int,len(rows));
    for i:=0; err==nil && i<len(rows); i++ {
        vals:=util.AppendWithPreallocation(
            []reflect.Value{reflect.ValueOf(sqlStmt)},
            getTableVals(&rows[i],AllButIDFilter),
        );
        rowValue:=reflect.ValueOf(c.db).MethodByName("QueryRow").Call(vals)[0];
        err=rowValue.Interface().(*sql.Row).Scan(&rv[i]);
    }
    return rv,err;
}

func Read[R DBTable](
        c *CRUD,
        rowVals R,
        filter ColumnFilter,
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
    vals:=util.AppendWithPreallocation(
        []reflect.Value{reflect.ValueOf(sqlStmt)},getTableVals(&rowVals,filter),
    );
    callVals:=reflect.ValueOf(c.db).MethodByName("Query").Call(vals);
    err:=util.GetErrorFromReflectValue(&callVals[1]);
    if err==nil {
        rows:=callVals[0].Interface().(*sql.Rows);
        defer rows.Close();
        var iter R;
        rowPntrs:=getTablePntrs(&iter,NoFilter);
        for err==nil && rows.Next() {
            potErr:=reflect.ValueOf(rows).MethodByName("Scan").Call(rowPntrs);
            err=util.GetErrorFromReflectValue(&potErr[0]);
            callback(&iter);
        }
    }
    return err;
}

func ReadAll[R DBTable](c *CRUD, callback func(val *R)) error {
    var iter R;
    sqlStmt:=fmt.Sprintf("SELECT * FROM %s;",getTableName(&iter));
    callVals:=reflect.ValueOf(c.db).MethodByName("Query").Call(
        []reflect.Value{reflect.ValueOf(sqlStmt)},
    );
    err:=util.GetErrorFromReflectValue(&callVals[1]);
    if err==nil {
        rows:=callVals[0].Interface().(*sql.Rows);
        defer rows.Close();
        rowPntrs:=getTablePntrs(&iter,NoFilter);
        for err==nil && rows.Next() {
            potErr:=reflect.ValueOf(rows).MethodByName("Scan").Call(rowPntrs);
            err=util.GetErrorFromReflectValue(&potErr[0]);
            callback(&iter);
        }
    }
    return err;
}

func Update[R DBTable](
        c *CRUD,
        searchVals R,
        searchValsFilter ColumnFilter,
        updateVals R,
        updateValsFilter ColumnFilter) (int64,error) {
    updateColumns:=getTableColumns(&updateVals,updateValsFilter);
    searchColumns:=getTableColumns(&searchVals,searchValsFilter);
    if len(updateColumns)==0 || len(searchColumns)==0 {
        return 0, util.FilterRemovedAllColumns("No rows were updated.");
    }
    setStr:=util.CSVGenerator(", ",func(iter int) (string,bool) {
        return fmt.Sprintf("%s=$%d",updateColumns[iter],iter+1),
            iter+1<len(updateColumns);
    });
    whereStr:=util.CSVGenerator(" AND ",func(iter int) (string,bool) {
        return fmt.Sprintf("%s=$%d",searchColumns[iter],iter+1+len(updateColumns)),
            iter+1<len(searchColumns);
    });
    sqlStmt:=fmt.Sprintf(
        "UPDATE %s SET %s WHERE %s;",getTableName(&searchVals),setStr,whereStr,
    );
    vals:=util.AppendWithPreallocation(
        []reflect.Value{reflect.ValueOf(sqlStmt)},
        getTableVals(&updateVals,updateValsFilter),
        getTableVals(&searchVals,searchValsFilter),
    );
    callVals:=reflect.ValueOf(c.db).MethodByName("Exec").Call(vals);
    err:=util.GetErrorFromReflectValue(&callVals[1]);
    if err==nil {
        res:=callVals[0].Interface().(sql.Result);
        return res.RowsAffected();
    }
    return 0 ,err;
}

func Delete[R DBTable](
        c *CRUD,
        searchVals R,
        searchValsFilter ColumnFilter) (int64,error) {
    columns:=getTableColumns(&searchVals,searchValsFilter);
    if len(columns)==0 {
        return 0, util.FilterRemovedAllColumns("No rows were deleted.");
    }
    whereStr:=util.CSVGenerator(" AND ",func(iter int)(string,bool) {
        return fmt.Sprintf("%s=$%d",columns[iter],iter+1),iter+1<len(columns);
    });
    sqlStmt:=fmt.Sprintf(
        "DELETE FROM %s WHERE %s;",getTableName(&searchVals),whereStr,
    );
    vals:=util.AppendWithPreallocation(
        []reflect.Value{reflect.ValueOf(sqlStmt)},
        getTableVals(&searchVals,searchValsFilter),
    );
    callVals:=reflect.ValueOf(c.db).MethodByName("Exec").Call(vals);
    err:=util.GetErrorFromReflectValue(&callVals[1]);
    if err==nil {
        res:=callVals[0].Interface().(sql.Result);
        return res.RowsAffected();
    }
    return 0, err;
}

func DeleteAll[R DBTable](c *CRUD) (int64,error) {
    var tmp R;
    return 0, nil;
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
