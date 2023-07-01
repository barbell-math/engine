package iter

import (
	"bufio"
	"os"
    "fmt"
)

func (i Iter[T])Consume() error {
    return i.ForEach(func(index int, val T) (IteratorFeedback, error) { 
        return Continue,nil;
    });
}

//TODO - look into closing channel
func (i Iter[T])ToChan(c chan T) {
    i.ForEach(func(i int, val T) (IteratorFeedback,error) {
        c <- val;
        return Continue,nil;
    });
}

func (i Iter[T])ToFile(src string, addNewLine bool) error {
    f,err:=os.Create(src);
    if err!=nil {
        return err;
    }
    w:=bufio.NewWriter(f);
    err=i.ForEach(func(index int, val T) (IteratorFeedback, error) {
        w.WriteString(fmt.Sprintf("%v",val));
        if addNewLine {
            w.WriteString("\n");
        }
        return Continue,nil;
    });
    w.Flush();
    f.Close();
    return err;
}

func (i Iter[T])Count() (int,error) {
    rv:=0;
    err:=i.ForEach(func(i int, val T) (IteratorFeedback,error) {
        rv++;
        return Continue,nil;
    });
    return rv,err;
}

func (i Iter[T])Collect() ([]T,error) {
    rv:=make([]T,0);
    err:=i.ForEach(func(index int, val T) (IteratorFeedback, error) {
        rv=append(rv,val);
        return Continue,nil;
    });
    return rv,err;
}

func (i Iter[T])CollectInto(buffer []T) (int,error) {
    j:=0;
    err:=i.ForEach(func(index int, val T) (IteratorFeedback, error) {
        buffer[j]=val;
        j++;
        if j<len(buffer) {
            return Continue,nil;
        }
        return Break,nil;
    });
    return j,err;
}

func (i Iter[T])AppendTo(orig *[]T) (int,error) {
    j:=0;
    err:=i.ForEach(func(index int, val T) (IteratorFeedback,error) {
        *orig=append(*orig,val);
        j++;
        return Continue,nil;
    });
    return j,err;
}

func (i Iter[T])All(op func(val T) (bool,error)) (bool,error) {
    rv:=true;
    err:=i.ForEach(func(index int, val T) (IteratorFeedback,error) {
        res,err:=op(val);
        if !res {
            rv=false;
            return Break,err;
        }
        return Continue,err;
    });
    return rv,err;
}

func (i Iter[T])Any(op func(val T) (bool,error)) (bool,error) {
    rv:=false;
    err:=i.ForEach(func(index int, val T) (IteratorFeedback, error) {
        res,err:=op(val);
        if res {
            rv=true;
            return Break,err;
        }
        return Continue,err;
    });
    return rv,err;
}

func (i Iter[T])Find(op func(val T) (bool,error)) (T,error,bool) {
    var rv T;
    iterState:=Continue;
    err:=i.ForEach(func(i int, val T) (IteratorFeedback,error){
        found,err:=op(val);
        if found {
            rv=val;
            iterState=Break;
        }
        return iterState,err;
    });
    if iterState==Break {
        return rv,err,true;
    }
    return rv,err,false;
}

func (i Iter[T])Index(op func(val T) (bool,error)) (int,error) {
    rv:=-1;
    err:=i.ForEach(func(i int, val T) (IteratorFeedback,error) {
        found,err:=op(val);
        if found {
            rv=i;
            return Break,err;
        }
        return Continue,err;
    });
    return rv,err;
}

func (i Iter[T])Nth(idx int) (T,error,bool) {
    var valRv T;
    found:=false;
    err:=i.ForEach(func(i int, val T) (IteratorFeedback,error) {
        if i==idx {
            valRv=val;
            found=true;
            return Break,nil;
        }
        return Continue,nil;
    });
    return valRv,err,found;
}

func (i Iter[T])Reduce(start T, op func(accum *T, iter T) error) (T,error) {
    accum:=start;
    err:=i.ForEach(func(index int, val T) (IteratorFeedback, error) {
        return Continue,op(&accum,val);
    });
    return accum,err;
}
