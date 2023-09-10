package db

import "github.com/barbell-math/engine/util/algo/iter"

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
        vals []any) iter.Iter[*S] {
    if SelectStmt.isQueryType(sqlStmt) {
        rows,err:=c.db.Query(sqlStmt,vals...);
        if err==nil {
            return readRows[S](rows);
        }
        return iter.ValElem[*S](nil,err,1);
    }
    return iter.ValElem[*S](nil,UnsupportedQueryType(
        "CustomReadQuery only accepts 'SELECT' query's.",
    ),1);
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
