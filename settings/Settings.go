package settings;

import (
    "os"
    "log"
    "fmt"
    "sync"
    "io/ioutil"
    "encoding/json"
    "github.com/carmichaeljr/powerlifting-engine/util"
)

var mu sync.Mutex;

type DatabaseInfo struct {
    DataVersion int `json:"dataVersion"`;
    Host string `json:"host"`;
    Port int `json:"port"`;
    Name string `json:"name"`;
};
type SqlScripts struct {
    GlobalInit string `json:"globalInit"`;
    ExerciseFocusInit string `json:"exerciseFocusInit"`;
    ExerciseTypeInit string `json:"exerciseTypeInit"`;
    ExerciseInit string `json:"exerciseInit"`;
};

type Settings struct {
    DBInfo DatabaseInfo `json:"database"`;
    SqlFiles SqlScripts `json:"sqlScripts"`;
};
var s Settings;

func ReadSettings(src string){
    var rv Settings;
    err:=util.ChainedErrorOps(
        func(r ...any) (any,error) { return os.Open(src); },
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
    return s.DBInfo.DataVersion;
}
func DBHost() string {
    return s.DBInfo.Host;
}
func DBPort() int {
    return s.DBInfo.Port;
}
func DBName() string {
    return s.DBInfo.Name;
}
func SQLGlobalInitScript() string {
    return s.SqlFiles.GlobalInit;
}
func ExerciseFocusInitData() string {
    return s.SqlFiles.ExerciseFocusInit;
}
func ExerciseTypeInitData() string {
    return s.SqlFiles.ExerciseTypeInit;
}
func ExerciseInitData() string {
    return s.SqlFiles.ExerciseInit;
}

func valid(set *Settings) (bool,error) {
    var rv bool=true;
    err:=util.ChainedErrorOps(
        func(r ...any) (any,error) {
            rv=(set.DBInfo.DataVersion>=0);
            return rv,util.ErrorOnBool(rv,util.DataVersionMalformed("Should be >=0."));
        }, func(r ...any) (any,error) {
            rv,err:=util.FileExists(set.SqlFiles.GlobalInit);
            return rv,util.ErrorOnBool(
                rv,util.SettingsFileNotFound(fmt.Sprintf("GlobalInit | %s",err)),
            );
        }, func(r ...any) (any,error) {
            rv,err:=util.FileExists(set.SqlFiles.ExerciseFocusInit);
            return rv,util.ErrorOnBool(
                rv,util.SettingsFileNotFound(
                    fmt.Sprintf(" ExerciseFocusInit %s",err),
                ),
            );
        }, func(r ...any) (any,error) {
            rv,err:=util.FileExists(set.SqlFiles.ExerciseTypeInit);
            return rv,util.ErrorOnBool(
                rv,util.SettingsFileNotFound(
                    fmt.Sprintf(" ExerciseTypeInit | %s",err),
                ),
            );
        }, func(r ...any) (any,error) {
            rv,err:=util.FileExists(set.SqlFiles.ExerciseInit);
            return rv,util.ErrorOnBool(
                rv,util.SettingsFileNotFound(
                    fmt.Sprintf(" ExerciseInit | %s",err),
                ),
            );
        },
    );
    return rv,err;
}

func _copy() Settings {
    var rv Settings;
    rv.DBInfo=DatabaseInfo{
        DataVersion: s.DBInfo.DataVersion,
        Host: s.DBInfo.Host,
        Port: s.DBInfo.Port,
        Name: s.DBInfo.Name,
    };
    rv.SqlFiles=SqlScripts{
        GlobalInit: s.SqlFiles.GlobalInit,
        ExerciseFocusInit: s.SqlFiles.ExerciseFocusInit,
        ExerciseTypeInit: s.SqlFiles.ExerciseTypeInit,
        ExerciseInit: s.SqlFiles.ExerciseInit,
    };
    return rv;
}
