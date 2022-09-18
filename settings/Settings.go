package settings;

import (
    "os"
    "log"
    //"fmt"
    "io/ioutil"
    "encoding/json"
    "github.com/carmichaeljr/powerlifting-engine/util"
)

var S settings=readSettings();

type settings struct {
    dataVersion int `json:"dataVersion"`;
};

func readSettings() settings {
    var rv settings;
    err:=util.ChainedErrorOps(
        func(r ...any) (any,error) { return os.Open("./settings.json"); },
        func(r ...any) (any,error) {
            defer r[0].(*os.File).Close();
            return ioutil.ReadAll(r[0].(*os.File));
        },
        func(r ...any) (any,error) {
            return nil,json.Unmarshal(r[1].([]byte),&rv);
        },
    );
    if err!=nil {
        log.Fatalf("An error occurred reading the settings file.\n %s",err);
    }
    return rv;
}
