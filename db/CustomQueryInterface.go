package db;

//func CustomSelectQuery[R DBTable](whereStmt string, whereVals []any) R
//func CustomUpdateQuery[R DBTable](
//        updateVals R,
//        updateFilter
//        ColFilter,
//        whereStmt string,
//        whereVals []any) R
//func CustomDeleteQuery[R DBTable](whereStmt string, whereVals []any) R

func CustomReadQuery[S any](
        c *DB,
        sqlStmt string,
        vals []any,
        callback func(r *S)) error {
    if SelectStmt.isQueryType(sqlStmt) {
        rows,err:=c.db.Query(sqlStmt,vals...);
        if err==nil {
            err=readRows(rows,callback);
            rows.Close();
        }
        return err;
    }
    return UnsupportedQueryType(
        "CustomReadQuery only accepts 'SELECT' query's.",
    );
}
