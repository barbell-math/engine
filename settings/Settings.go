package settings;

import (
    "os"
    "log"
    "fmt"
    "sync"
    "io/ioutil"
    "encoding/json"
    customIO "github.com/barbell-math/engine/util/io"
    customerr "github.com/barbell-math/engine/util/err"
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
};
type SetupData struct {
    ExerciseFocusInit string `json:"exerciseFocusInit"`;
    ExerciseTypeInit string `json:"exerciseTypeInit"`;
    ExerciseInit string `json:"exerciseInit"`;
    ClientInit string `json:"clientInit"`;
    StateGeneratorInit string `json:"stateGeneratorInit"`;
    PotentialSurfaceInit string `json:"potentialSurfaceInit"`;
    RotationInit string `json:"rotationInit"`;
    TrainingLogInit string `json:"trainingLogInit"`;
};

type Settings struct {
    DBInfo DatabaseInfo `json:"database"`;
    SqlFiles SqlScripts `json:"sqlScripts"`;
    InitData SetupData `json:"setupData"`;
};
var s Settings;

func ReadSettings(src string){
    var rv Settings;
    err:=customerr.ChainedErrorOps(
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
    return s.InitData.ExerciseFocusInit;
}
func ExerciseTypeInitData() string {
    return s.InitData.ExerciseTypeInit;
}
func ExerciseInitData() string {
    return s.InitData.ExerciseInit;
}
func ClientInitData() string {
    return s.InitData.ClientInit;
}
func StateGeneratorInitData() string {
    return s.InitData.StateGeneratorInit;
}
func PotentialSurfaceInitData() string {
    return s.InitData.PotentialSurfaceInit;
}
func RotationInitData() string {
    return s.InitData.RotationInit;
}
func TrainingLogInitData() string {
    return s.InitData.TrainingLogInit;
}

func valid(set *Settings) (bool,error) {
    var rv bool=true;
    err:=customerr.ChainedErrorOps(
        func(r ...any) (any,error) {
            rv=(set.DBInfo.DataVersion>=0);
            return rv,customerr.ErrorOnBool(rv,DataVersionMalformed("Should be >=0."));
        }, func(r ...any) (any,error) {
            rv,err:=customIO.FileExists(set.SqlFiles.GlobalInit);
            return rv,customerr.ErrorOnBool(
                rv,SettingsFileNotFound(fmt.Sprintf("GlobalInit | %v",err)),
            );
        }, func(r ...any) (any,error) {
            rv,err:=customIO.FileExists(set.InitData.ExerciseFocusInit);
            return rv,customerr.ErrorOnBool(
                rv,SettingsFileNotFound(fmt.Sprintf("ExerciseFocusInit | %v",err)),
            );
        }, func(r ...any) (any,error) {
            rv,err:=customIO.FileExists(set.InitData.ExerciseTypeInit);
            return rv,customerr.ErrorOnBool(
                rv,SettingsFileNotFound(fmt.Sprintf("ExerciseTypeInit | %v",err)),
            );
        }, func(r ...any) (any,error) {
            rv,err:=customIO.FileExists(set.InitData.ExerciseInit);
            return rv,customerr.ErrorOnBool(
                rv,SettingsFileNotFound(fmt.Sprintf("ExerciseInit | %v",err)),
            );
        }, func(r ...any) (any,error) {
            rv,err:=customIO.FileExists(set.InitData.ClientInit);
            return rv,customerr.ErrorOnBool(
                rv,SettingsFileNotFound(fmt.Sprintf("ClientInit | %v",err)),
            );
        }, func(r ...any) (any,error) {
            rv,err:=customIO.FileExists(set.InitData.StateGeneratorInit);
            return rv,customerr.ErrorOnBool(
                rv,SettingsFileNotFound(fmt.Sprintf("StateGeneratorInit | %v",err)),
            );
        }, func(r ...any) (any,error) {
            rv,err:=customIO.FileExists(set.InitData.PotentialSurfaceInit);
            return rv,customerr.ErrorOnBool(
                rv,SettingsFileNotFound(fmt.Sprintf("PotentialSurfaceInit | %v",err)),
            );
        }, func(r ...any) (any,error) {
            rv,err:=customIO.FileExists(set.InitData.RotationInit);
            return rv,customerr.ErrorOnBool(
                rv,SettingsFileNotFound(fmt.Sprintf("RotationInit | %v",err)),
            );
        }, func(r ...any) (any,error) {
            rv,err:=customIO.FileExists(set.InitData.TrainingLogInit);
            return rv,customerr.ErrorOnBool(
                rv,SettingsFileNotFound(fmt.Sprintf("TrainingLogInit | %v",err)),
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
    };
    rv.InitData=SetupData{
        ExerciseFocusInit: s.InitData.ExerciseFocusInit,
        ExerciseTypeInit: s.InitData.ExerciseTypeInit,
        ExerciseInit: s.InitData.ExerciseInit,
        ClientInit: s.InitData.ClientInit,
        StateGeneratorInit: s.InitData.StateGeneratorInit,
        PotentialSurfaceInit: s.InitData.PotentialSurfaceInit,
        RotationInit: s.InitData.RotationInit,
        TrainingLogInit: s.InitData.TrainingLogInit,
    };
    return rv;
}
