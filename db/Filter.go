package db;

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
