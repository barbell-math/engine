package settings;

import (
    "os"
    "log"
    //"fmt"
    "sync"
    "io/ioutil"
    "encoding/json"
    "github.com/carmichaeljr/powerlifting-engine/util"
)

var mu sync.Mutex;
var settingsSrc string="./settings.json";

type Settings struct {
    DataVersion int `json:"dataVersion"`;
};
var s Settings;

func ReadSettings(){
    var rv Settings;
    err:=util.ChainedErrorOps(
        func(r ...any) (any,error) { return os.Open(settingsSrc); },
        func(r ...any) (any,error) {
            defer r[0].(*os.File).Close();
            return ioutil.ReadAll(r[0].(*os.File));
        }, func(r ...any) (any,error) {
            return nil,json.Unmarshal(r[1].([]byte),&rv);
        }, func(r ...any) (any,error) {
            return valid(&rv);
    });
    if err!=nil {
        log.Fatalf("An error occurred reading the settings file.\n %s",err);
    }
    mu.Lock();
    s=rv;
    mu.Unlock();
}

type SettingsMod func(s *Settings);
func Modify(ops ...SettingsMod) error {
    working:=_copy();
    for _,op:=range(ops) {
        op(&working);
    }
    _,err:=valid(&working);
    if err==nil {
        mu.Lock();
        s=working;
        mu.Unlock();
    }
    return err;
}

func DataVersion() int {
    return s.DataVersion;
}

func valid(set *Settings) (bool,error) {
    var rv bool=true;
    err:=util.ChainedErrorOps(
        func(r ...any) (any,error) {
            rv=(set.DataVersion>0);
            return rv,util.ErrorOnBool(rv,util.DataVersionMalformed(""));
        },
    );
    return rv,err;
}

func _copy() Settings {
    var rv Settings;
    rv.DataVersion=s.DataVersion;
    return rv;
}
