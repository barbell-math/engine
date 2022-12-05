package settings;

import (
    "testing"
    "github.com/barbell-math/block/util/test"
)

func TestReadSettings(t *testing.T){
    ReadSettings("../testData/dbTestSettings.json");
    test.BasicTest(0 ,DataVersion(),"Data version not parsed correctly.",t);
    test.BasicTest("localhost",DBHost(),"DB host not parsed correctly.",t);
    test.BasicTest(5432,DBPort(),"DB port not parsed correctly.",t);
    test.BasicTest("dbTest",DBName(),"DB name not parsed correctly.",t);
    test.BasicTest("../sql/globalInit.sql",SQLGlobalInitScript(),
        "Sql scripts not parsed correctly.",t,
    );
}

func TestModify(t *testing.T){
    //prev:=s.DataVersion;
    err:=Modify(func(s *Settings){ s.DBInfo.DataVersion=10; });
    test.BasicTest(nil,err,"Modifying the data version was no successful.",t);
    test.BasicTest(10,s.DBInfo.DataVersion,"Data version was not updated.",t);
    err=Modify(func(s *Settings){ s.DBInfo.DataVersion=-1; });
    if !IsDataVersionMalformed(err) {
        test.FormatError(
            DataVersionMalformed(""),err,
            "Malformed data version was not caught.",t,
        );
    }
    test.BasicTest(10,s.DBInfo.DataVersion,"Data version was updated.",t);
}
