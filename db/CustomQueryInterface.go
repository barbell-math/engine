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
        callback func(r *S) bool) error {
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

func CustomDeleteQuery(c *DB, sqlStmt string, vals []any) (int64,error) {
    if DeleteStmt.isQueryType(sqlStmt) {
        res,err:=c.db.Exec(sqlStmt,vals...);
        if err==nil {
            return res.RowsAffected();
        }
        return 0, err;
    }
    return 0, UnsupportedQueryType(
        "CustomDeleteQuery only accepts 'DELETE' query's.",
    );
}
