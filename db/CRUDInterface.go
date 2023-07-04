package db

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/barbell-math/block/util/algo"
	"github.com/barbell-math/block/util/algo/iter"
	"github.com/barbell-math/block/util/io/csv"
	customReflect "github.com/barbell-math/block/util/reflect"
)

func Create[R DBTable](c *DB, rows ...R) ([]int,error) {
    if len(rows)==0 {
        return []int{},sql.ErrNoRows;
    }
    columns,_:=iter.SliceElems(getTableColumns(&rows[0],AllButIDFilter)).Map(
    func(index int, val string) (string, error) {
        return val,nil;
    }).Collect();
    if len(columns)==0 {
        return []int{},FilterRemovedAllColumns("Row was not added to database.");
    }
    intoStr,_,_:=csv.Flatten(iter.SliceElems([][]string{columns}),",").Nth(0);
    valuesStr:=csv.CSVGenerator(",",func(iter int) (string,bool) {
        return fmt.Sprintf("$%d",iter+1), iter+1<len(columns);
    });
    sqlStmt:=fmt.Sprintf(
        "INSERT INTO %s(%s) VALUES (%s) RETURNING Id;",
        getTableName(&rows[0]),intoStr,valuesStr,
    );
    var err error=nil;
    rv:=make([]int,len(rows));
    for i:=0; err==nil && i<len(rows); i++ {
        rv[i],err=getQueryRowReflectResults(c,algo.AppendWithPreallocation(
                []reflect.Value{reflect.ValueOf(sqlStmt)},
                getTableVals(&rows[i],AllButIDFilter),
        ));
    }
    return rv,err;
}

func Read[R DBTable](
        c *DB,
        rowVals R,
        filter algo.Filter[string]) iter.Iter[*R] {
    columns,_:=iter.SliceElems(getTableColumns(&rowVals,filter)).Map(
    func(index int, val string) (string, error) {
        return fmt.Sprintf("%s=$%d",val,index+1),nil;
    }).Collect();
    if len(columns)==0 {
        return iter.ValElem[*R](nil,
            FilterRemovedAllColumns("No value rows were selected."),1,
        );
    }
    valuesStr,_,_:=csv.Flatten(iter.SliceElems([][]string{columns})," AND ").Nth(0);
    sqlStmt:=fmt.Sprintf(
        "SELECT * FROM %s WHERE %s;",getTableName(&rowVals),valuesStr,
    );
    return getQueryReflectResults[R](c,
        algo.AppendWithPreallocation(
            []reflect.Value{reflect.ValueOf(sqlStmt)},
            getTableVals(&rowVals,filter),
        ),
    );
}

func ReadAll[R DBTable](c *DB) iter.Iter[*R] {
    var tmp R;
    sqlStmt:=fmt.Sprintf("SELECT * FROM %s;",getTableName(&tmp));
    return getQueryReflectResults[R](c,
        []reflect.Value{reflect.ValueOf(sqlStmt)},
    );
}

func Update[R DBTable](
        c *DB,
        searchVals R,
        searchValsFilter algo.Filter[string],
        updateVals R,
        updateValsFilter algo.Filter[string]) (int64,error) {
    updateColumns,_:=iter.SliceElems(getTableColumns(&updateVals,updateValsFilter)).Map(
    func(index int, val string) (string, error) {
        return fmt.Sprintf("%s=$%d",val,index+1),nil;
    }).Collect();
    searchColumns,_:=iter.SliceElems(getTableColumns(&searchVals,searchValsFilter)).Map(
    func(index int, val string) (string, error) {
        return fmt.Sprintf("%s=$%d",val,index+1+len(updateColumns)),nil;
    }).Collect();
    if len(updateColumns)==0 || len(searchColumns)==0 {
        return 0, FilterRemovedAllColumns("No rows were updated.");
    }
    setStr,_,_:=csv.Flatten(iter.SliceElems([][]string{updateColumns}),", ").Nth(0);
    whereStr,_,_:=csv.Flatten(iter.SliceElems([][]string{searchColumns})," AND ").Nth(0);
    sqlStmt:=fmt.Sprintf(
        "UPDATE %s SET %s WHERE %s;",getTableName(&searchVals),setStr,whereStr,
    );
    return getExecReflectResults(c,
        algo.AppendWithPreallocation(
            []reflect.Value{reflect.ValueOf(sqlStmt)},
            getTableVals(&updateVals,updateValsFilter),
            getTableVals(&searchVals,searchValsFilter),
        ),
    );
}

func UpdateAll[R DBTable](
        c *DB,
        updateVals R,
        updateValsFilter algo.Filter[string]) (int64,error) {
    updateColumns,_:=iter.SliceElems(getTableColumns(&updateVals,updateValsFilter)).Map(
    func(index int, val string) (string, error) {
        return fmt.Sprintf("%s=$%d",val,index+1),nil;
    }).Collect();
    if len(updateColumns)==0 {
        return 0, FilterRemovedAllColumns("No rows were updated.");
    }
    setStr,_,_:=csv.Flatten(iter.SliceElems([][]string{updateColumns}),", ").Nth(0);
    sqlStmt:=fmt.Sprintf("UPDATE %s SET %s;",getTableName(&updateVals),setStr);
    return getExecReflectResults(c,
        algo.AppendWithPreallocation(
            []reflect.Value{reflect.ValueOf(sqlStmt)},
            getTableVals(&updateVals,updateValsFilter),
        ),
    );
}

func Delete[R DBTable](
        c *DB,
        searchVals R,
        searchValsFilter algo.Filter[string]) (int64,error) {
    columns,_:=iter.SliceElems(getTableColumns(&searchVals,searchValsFilter)).Map(
    func(index int, val string) (string, error) {
        return fmt.Sprintf("%s=$%d",val,index+1),nil;
    }).Collect();
    if len(columns)==0 {
        return 0, FilterRemovedAllColumns("No rows were deleted.");
    }
    whereStr,_,_:=csv.Flatten(iter.SliceElems([][]string{columns})," AND ").Nth(0);
    sqlStmt:=fmt.Sprintf(
        "DELETE FROM %s WHERE %s;",getTableName(&searchVals),whereStr,
    );
    return getExecReflectResults(c,
        algo.AppendWithPreallocation(
            []reflect.Value{reflect.ValueOf(sqlStmt)},
            getTableVals(&searchVals,searchValsFilter),
        ),
    );
}

func DeleteAll[R DBTable](c *DB) (int64,error) {
    var tmp R;
    sqlStmt:=fmt.Sprintf("DELETE FROM %s;",getTableName(&tmp));
    return getExecReflectResults(c,[]reflect.Value{reflect.ValueOf(sqlStmt)});
}

func getQueryReflectResults[R DBTable](
        c *DB,
        vals []reflect.Value) iter.Iter[*R] {
    reflectVals:=reflect.ValueOf(c.db).MethodByName("Query").Call(vals);
    err:=customReflect.GetErrorFromReflectValue(&reflectVals[1]);
    if err==nil {
        rows:=reflectVals[0].Interface().(*sql.Rows);
        return readRows[R](rows);
    }
    return iter.ValElem[*R](nil,err,1);
}

func getExecReflectResults(c *DB, vals []reflect.Value) (int64,error) {
    reflectVals:=reflect.ValueOf(c.db).MethodByName("Exec").Call(vals);
    err:=customReflect.GetErrorFromReflectValue(&reflectVals[1]);
    if err==nil {
        res:=reflectVals[0].Interface().(sql.Result);
        return res.RowsAffected();
    }
    return 0 ,err;
}

func getQueryRowReflectResults(c *DB, vals []reflect.Value) (int,error) {
    var rv int;
    reflectVal:=reflect.ValueOf(c.db).MethodByName("QueryRow").Call(vals)[0]
    rowVal:=reflectVal.Interface().(*sql.Row);
    err:=rowVal.Scan(&rv);
    return rv,err;
}

//These are convenience functions that allow for inline function calls to be used
func getTableName[R DBTable](row *R) string {
    //It is safe to ignore the err this because row is guaranteed to be a struct
    n,_:=customReflect.GetStructName(row);
    return n;
}

func getTableColumns[R DBTable](row *R, filter algo.Filter[string]) []string {
    //It is safe to ignore the err this because row is guaranteed to be a struct
    rv,_:=customReflect.GetStructFieldNames(row,filter);
    return rv;
}

func getTableVals[R DBTable](row *R, filter algo.Filter[string]) []reflect.Value {
    //It is safe to ignore the err this because row is guaranteed to be a struct
    rv,_:=customReflect.GetStructVals(row,filter);
    return rv;
}

func getTablePntrs[R DBTable](row *R,filter algo.Filter[string]) []reflect.Value {
    //It is safe to ignore the err this because row is guaranteed to be a struct
    rv,_:=customReflect.GetStructFieldPntrs(row,filter);
    return rv;
}
