package io;

import (
    "os"
    "errors"
    "testing"
    "github.com/barbell-math/block/util/test"
)

func TestFileExistsRealFile(t *testing.T){
    v,e:=FileExists("testData/realFile.txt");
    test.BasicTest(true,v,"File that exists was flagged as not.",t);
    test.BasicTest(nil,e,"An error was generated when it shouldn't have been.",t);
}

func TestFileExistsFakeFile(t *testing.T){
    v,e:=FileExists("testData/realFile1.txt");
    test.BasicTest(false,v,"File that doesn't exist was flagged as existing.",t);
    if errors.Is(e,os.ErrNotExist) {
        test.FormatError(os.ErrNotExist,e,"An incorrect error was generated.",t);
    }
}

func TestFileExistsDir(t *testing.T){
    v,e:=FileExists("testData/");
    test.BasicTest(false,v,"File that doesn't exist was flagged as existing.",t);
    test.BasicTest(nil,e,"An error was generated when it shouldn't have been.",t);
}

