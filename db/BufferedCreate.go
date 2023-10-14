package db;

import (
    "fmt"
    customerr "github.com/barbell-math/engine/util/err"
)

type BufferedCreate[R DBTable] struct {
    buf []R;
    bufCntr int;
    succeeded int;
    failed int;
};

func NewBufferedCreate[R DBTable](bufSize int) (BufferedCreate[R],error) {
    if bufSize<1 {
        return BufferedCreate[R]{}, customerr.ValOutsideRange(fmt.Sprintf(
            "bufSize needs to be >=1. | %d",bufSize,
        ));
    }
    return BufferedCreate[R]{
        buf: make([]R,bufSize),
    },nil;
}

func (b *BufferedCreate[R])Succeeded() int { return b.succeeded; }
func (b *BufferedCreate[R])Failed() int { return b.failed; }

func (b *BufferedCreate[R])Write(c *DB, rows ...R) error {
    var rv error;
    for _,r:=range(rows) {
        b.buf[b.bufCntr]=r;
        b.bufCntr++;
        if b.bufCntr==len(b.buf) {
            rv=b.Flush(c);
        }
    }
    return rv;
}

func (b *BufferedCreate[R])Flush(c *DB) error {
    // To avoid copying the whole buffer when it is full we only copy it when
    // it is not completely full.
    bufPntr:=&b.buf;
    if b.bufCntr<len(b.buf) {
        tmp:=b.buf[0:b.bufCntr];
        bufPntr=&tmp;
    }
    added,err:=Create(c,*bufPntr...);
    succeeded:=0;
    for _,v:=range(added) {
        if v>0 {
            succeeded+=1;
        }
    }
    b.succeeded+=succeeded;
    b.failed+=(b.bufCntr-succeeded);
    b.bufCntr-=len(*bufPntr);
    return err;
}
