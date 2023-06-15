package db

import (
	"fmt"
	"testing"

	customerr "github.com/barbell-math/block/util/err"
	"github.com/barbell-math/block/util/test"
)

func TestCreateBufferedCreate(t *testing.T) {
	tmp,err:=NewBufferedCreate[Client](0);
    if !customerr.IsValOutsideRange(err) {
        test.FormatError(customerr.ValOutsideRange(""),err,
            "Incorrect error was returned.",t,
        );
    }
	tmp,err=NewBufferedCreate[Client](1);
    test.BasicTest(nil,err,
        "An error was returned when it should not have been.",t,
    );
    test.BasicTest(1,len(tmp.buf),"Buffer was created with incorrect length.",t);
}

func TestBufferedCreateWrite(t *testing.T) {
    setup();
	tmp,_:=NewBufferedCreate[Client](10);
    for i:=0; i<9; i++ {
        tmp.Write(&testDB,Client{Id: i, Email: fmt.Sprintf("%d",i)});
    }
    test.BasicTest(9,tmp.bufCntr,
        "Buffer counter was not incremented properly.",t,
    );
    for i:=0; i<9; i++ {
        test.BasicTest(i,tmp.buf[i].Id,
            "The correct values were not in the buffer.",t,
        );
    }
    err:=tmp.Write(&testDB,Client{Id: 10, Email: "10"});
    test.BasicTest(0,tmp.bufCntr,
        "Buffer counter was not incremented properly.",t,
    );
    test.BasicTest(nil,err,
        "An error was returned when it should not have been.",t,
    );
    cnt,_:=ReadAll[Client](&testDB).Count();
    test.BasicTest(10,cnt,
        "The correct number of values were not created.",t,
    );
}

func TestBufferedCreateWriteRollover(t *testing.T) {
    setup();
	tmp,_:=NewBufferedCreate[Client](10);
    for i:=0; i<20; i++ {
        tmp.Write(&testDB,Client{Id: i, Email: fmt.Sprintf("%d",i)});
    }
    test.BasicTest(0,tmp.bufCntr,
        "Buffer counter was not incremented properly.",t,
    );
    cnt,_:=ReadAll[Client](&testDB).Count();
    test.BasicTest(20,cnt,
        "The correct number of values were not created.",t,
    );
}

func TestBufferedCreateFlush(t *testing.T) {
    setup();
	tmp,_:=NewBufferedCreate[Client](10);
    for i:=0; i<5; i++ {
        tmp.Write(&testDB,Client{Id: i, Email: fmt.Sprintf("%d",i)});
    }
    test.BasicTest(5,tmp.bufCntr,
        "Buffer counter was not incremented properly.",t,
    );
    for i:=0; i<5; i++ {
        test.BasicTest(i,tmp.buf[i].Id,
            "The correct values were not in the buffer.",t,
        );
    }
    err:=tmp.Flush(&testDB);
    test.BasicTest(0,tmp.bufCntr,
        "Buffer counter was not incremented properly.",t,
    );
    test.BasicTest(nil,err,
        "An error was returned when it should not have been.",t,
    );
    cnt,_:=ReadAll[Client](&testDB).Count();
    test.BasicTest(5,cnt,
        "The correct number of values were not created.",t,
    );
}
