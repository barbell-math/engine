package db;

import (
    "fmt"
    "reflect"
    "database/sql"
    "github.com/carmichaeljr/powerlifting-engine/util"
)

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
        rv[i],err=getQueryRowReflectResults(c,util.AppendWithPreallocation(
                []reflect.Value{reflect.ValueOf(sqlStmt)},
                getTableVals(&rows[i],AllButIDFilter),
        ));
    }
    return rv,err;
}

func Read[R DBTable](
        c *CRUD,
        rowVals R,
        filter util.Filter[string],
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
    return getQueryReflectResults(c,
        util.AppendWithPreallocation(
            []reflect.Value{reflect.ValueOf(sqlStmt)},
            getTableVals(&rowVals,filter),
        ), callback,
    );
}

func ReadAll[R DBTable](c *CRUD, callback func(val *R)) error {
    var tmp R;
    sqlStmt:=fmt.Sprintf("SELECT * FROM %s;",getTableName(&tmp));
    return getQueryReflectResults(c,
        []reflect.Value{reflect.ValueOf(sqlStmt)},
        callback,
    );
}

func Update[R DBTable](
        c *CRUD,
        searchVals R,
        searchValsFilter util.Filter[string],
        updateVals R,
        updateValsFilter util.Filter[string]) (int64,error) {
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
    return getExecReflectResults(c,
        util.AppendWithPreallocation(
            []reflect.Value{reflect.ValueOf(sqlStmt)},
            getTableVals(&updateVals,updateValsFilter),
            getTableVals(&searchVals,searchValsFilter),
        ),
    );
}

func UpdateAll[R DBTable](
        c *CRUD,
        updateVals R,
        updateValsFilter util.Filter[string]) (int64,error) {
    updateColumns:=getTableColumns(&updateVals,updateValsFilter);
    if len(updateColumns)==0 {
        return 0, util.FilterRemovedAllColumns("No rows were updated.");
    }
    setStr:=util.CSVGenerator(", ",func(iter int) (string,bool) {
        return fmt.Sprintf("%s=$%d",updateColumns[iter],iter+1),
            iter+1<len(updateColumns);
    });
    sqlStmt:=fmt.Sprintf("UPDATE %s SET %s;",getTableName(&updateVals),setStr);
    return getExecReflectResults(c,
        util.AppendWithPreallocation(
            []reflect.Value{reflect.ValueOf(sqlStmt)},
            getTableVals(&updateVals,updateValsFilter),
        ),
    );
}

func Delete[R DBTable](
        c *CRUD,
        searchVals R,
        searchValsFilter util.Filter[string]) (int64,error) {
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
    return getExecReflectResults(c,
        util.AppendWithPreallocation(
            []reflect.Value{reflect.ValueOf(sqlStmt)},
            getTableVals(&searchVals,searchValsFilter),
        ),
    );
}

func DeleteAll[R DBTable](c *CRUD) (int64,error) {
    var tmp R;
    sqlStmt:=fmt.Sprintf("DELETE FROM %s;",getTableName(&tmp));
    return getExecReflectResults(c,[]reflect.Value{reflect.ValueOf(sqlStmt)});
}

func getQueryReflectResults[R DBTable](
        c *CRUD,
        vals []reflect.Value,
        callback func(val *R)) error {
    reflectVals:=reflect.ValueOf(c.db).MethodByName("Query").Call(vals);
    err:=util.GetErrorFromReflectValue(&reflectVals[1]);
    if err==nil {
        rows:=reflectVals[0].Interface().(*sql.Rows);
        defer rows.Close();
        //var iter R;
        err=readRows(rows,callback);
    }
    return err;
}

func getExecReflectResults(c *CRUD, vals []reflect.Value) (int64,error) {
    reflectVals:=reflect.ValueOf(c.db).MethodByName("Exec").Call(vals);
    err:=util.GetErrorFromReflectValue(&reflectVals[1]);
    if err==nil {
        res:=reflectVals[0].Interface().(sql.Result);
        return res.RowsAffected();
    }
    return 0 ,err;
}

func getQueryRowReflectResults(c *CRUD, vals []reflect.Value) (int,error) {
    var rv int;
    reflectVal:=reflect.ValueOf(c.db).MethodByName("QueryRow").Call(vals)[0]
    rowVal:=reflectVal.Interface().(*sql.Row);
    err:=rowVal.Scan(&rv);
    return rv,err;
}

//These are convenience functions that allow for inline function calls to be used
func getTableName[R DBTable](row *R) string {
    //It is safe to ignore the err this because row is guaranteed to be a struct
    n,_:=util.GetStructName(row);
    return n;
}

func getTableColumns[R DBTable](row *R, filter util.Filter[string]) []string {
    //It is safe to ignore the err this because row is guaranteed to be a struct
    rv,_:=util.GetStructFieldNames(row,filter);
    return rv;
}

func getTableVals[R DBTable](row *R, filter util.Filter[string]) []reflect.Value {
    //It is safe to ignore the err this because row is guaranteed to be a struct
    rv,_:=util.GetStructVals(row,filter);
    return rv;
}

func getTablePntrs[R DBTable](row *R,filter util.Filter[string]) []reflect.Value {
    //It is safe to ignore the err this because row is guaranteed to be a struct
    rv,_:=util.GetStructFieldPntrs(row,filter);
    return rv;
}
