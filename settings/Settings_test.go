package settings;

import (
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
)

func TestReadSettings(t *testing.T){
    ReadSettings("../testData/dbTestSettings.json");
    testUtil.BasicTest(0 ,DataVersion(),"Data version not parsed correctly.",t);
    testUtil.BasicTest("localhost",DBHost(),"DB host not parsed correctly.",t);
    testUtil.BasicTest(5432,DBPort(),"DB port not parsed correctly.",t);
    testUtil.BasicTest("dbTest",DBName(),"DB name not parsed correctly.",t);
    testUtil.BasicTest("../sql/globalInit.sql",SQLGlobalInitScript(),
        "Sql scripts not parsed correctly.",t,
    );
}

func TestModify(t *testing.T){
    //prev:=s.DataVersion;
    err:=Modify(func(s *Settings){ s.DBInfo.DataVersion=10; });
    testUtil.BasicTest(nil,err,"Modifying the data version was no successful.",t);
    testUtil.BasicTest(10,s.DBInfo.DataVersion,"Data version was not updated.",t);
    err=Modify(func(s *Settings){ s.DBInfo.DataVersion=-1; });
    if !util.IsDataVersionMalformed(err) {
        testUtil.FormatError(
            util.DataVersionMalformed(""),err,
            "Malformed data version was not caught.",t,
        );
    }
    testUtil.BasicTest(10,s.DBInfo.DataVersion,"Data version was updated.",t);
}
