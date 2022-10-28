package db;

//import (
//    "database/sql"
//)

//func CustomSelectQuery[R DBTable](whereStmt string, whereVals []any) R
//func CustomUpdateQuery[R DBTable](
//        updateVals R,
//        updateFilter
//        ColFilter,
//        whereStmt string,
//        whereVals []any) R
//func CustomDeleteQuery[R DBTable](whereStmt string, whereVals []any) R

//func CustomReadQuery[S any](
//        c *CRUD,
//        sqlStmt string,
//        vals ...any,
//        callback func(r *S)) error {
//    var iter R;
//    if rows,err:=c.db.Query(sqlStmt,vals); err==nil {
//        defer rows.Close();
//        rowPntrs:=getTablePntrs(&iter,NoFilter);
//        for rows.Next() {
//            err=rows.Scan();
//        }
//    } else {
//        return err;
//    }
//}
