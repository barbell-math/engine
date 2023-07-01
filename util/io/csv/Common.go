package csv

import (
	"encoding/csv"
	"io"
	"os"
	"strings"

	"github.com/barbell-math/block/util/algo/iter"
)

func CSVGenerator(sep string, callback func(iter int) (string,bool)) string {
    var sb strings.Builder;
    var temp string;
    cont:=true;
    for i:=0; cont; i++ {
        temp,cont=callback(i);
        sb.WriteString(temp);
        if cont {
            sb.WriteString(sep);
        }
    }
    return sb.String();
}

func Flatten(elems iter.Iter[[]string], sep string) iter.Iter[string] {
    return iter.Next(elems,
    func(index int,
        val []string,
        status iter.IteratorFeedback,
    ) (iter.IteratorFeedback, string, error) {
        if status==iter.Break {
            return iter.Break,"",nil;
        }
        var sb strings.Builder;
        for i,v:=range(val) {
            sb.WriteString(v);
            if i!=len(val)-1 {
                sb.WriteString(sep);
            }
        }
        return iter.Continue,sb.String(),nil;
    });
}

func CSVFileSplitter(src string, delim rune, comment rune) iter.Iter[[]string] {
    var reader *csv.Reader=nil;
    file,err:=os.Open(src);
    if err==nil {
        reader=csv.NewReader(file);
        reader.Comma=delim;
        reader.Comment=comment;
    }
    return func(f iter.IteratorFeedback) ([]string, error, bool) {
        if f==iter.Break || err!=nil {
            file.Close();
            return []string{},err,false;
        }
        cols,readerErr:=reader.Read();
        if readerErr!=nil {
            if readerErr==io.EOF {
                return cols,nil,false;
            }
            return []string{},readerErr,false;
        }
        return cols,readerErr,true;
    }
}
