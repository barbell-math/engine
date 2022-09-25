package settings;

import (
    "testing"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/testUtil"
)

func TestMain(m *testing.M){
    setup();
    m.Run();
    teardown();
}

func setup(){
    settingsSrc="../settings.json";
}
func teardown(){
    settingsSrc="./settings.json";
}

func TestModify(t *testing.T){
    //prev:=s.DataVersion;
    err:=Modify(func(s *Settings){ s.DataVersion=10; });
    testUtil.BasicTest(nil,err,"Modifying the data version was no successful.",t);
    testUtil.BasicTest(10,s.DataVersion,"Data version was not updated.",t);
    err=Modify(func(s *Settings){ s.DataVersion=-1; });
    if !util.IsDataVersionMalformed(err) {
        testUtil.FormatError(
            util.DataVersionMalformed(""),err,
            "Malformed data version was not caught.",t,
        );
    }
    testUtil.BasicTest(10,s.DataVersion,"Data version was updated.",t);
}
