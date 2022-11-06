package db;

import (
    "github.com/carmichaeljr/powerlifting-engine/util"
)

//func CustomSelectQuery[R DBTable](whereStmt string, whereVals []any) R
//func CustomUpdateQuery[R DBTable](
//        updateVals R,
//        updateFilter
//        ColFilter,
//        whereStmt string,
//        whereVals []any) R
//func CustomDeleteQuery[R DBTable](whereStmt string, whereVals []any) R

func CustomReadQuery[S any](
        c *CRUD,
        sqlStmt string,
        vals []any,
        callback func(r *S)) error {
    if SelectStmt.isQueryType(sqlStmt) {
        rows,err:=c.db.Query(sqlStmt,vals...);
        if err==nil {
            defer rows.Close();
            err=readRows(rows,callback);
        }
        return err;
    }
    return util.UnsupportedQueryType(
        "CustomReadQuery only accepts 'SELECT' query's.",
    );
}
