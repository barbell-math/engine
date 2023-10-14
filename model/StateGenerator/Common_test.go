package stateGenerator

import (
	"fmt"
	"testing"

	"github.com/barbell-math/engine/db"
	"github.com/barbell-math/engine/model/testSetup"
	"github.com/barbell-math/engine/settings"
	customerr "github.com/barbell-math/engine/util/err"
	logUtil "github.com/barbell-math/engine/util/io/log"
)

var testDB db.DB;
const (
    DP_DEBUG int=1<<iota
    MS_DEBUG
    MS_PARALLEL_RESULT_DEBUG
)

func setupLog[T any](dest *logUtil.Logger[T], path string){
    customerr.PanicOnError(func() error {
        var err error;
        *dest,err=logUtil.NewLog[T](logUtil.Debug,path,false);
        return err;
    });
}
func setupLogs(debugFile string, logFlags int) (func()) {
    if DP_DEBUG&logFlags==DP_DEBUG {
        setupLog(&SLIDING_WINDOW_DP_DEBUG,
            fmt.Sprintf("%s.dataPoint.log",debugFile),
        );
    }
    if MS_DEBUG&logFlags==MS_DEBUG {
        setupLog(&SLIDING_WINDOW_MS_DEBUG,
            fmt.Sprintf("%s.modelState.log",debugFile),
        );
    }
    if MS_PARALLEL_RESULT_DEBUG&logFlags==MS_PARALLEL_RESULT_DEBUG {
        setupLog(&SLIDING_WINDOW_MS_PARALLEL_RESULT_DEBUG,
            fmt.Sprintf("%s.modelState.log",debugFile),
        );
    }
    return func(){
        if DP_DEBUG&logFlags==DP_DEBUG {
            SLIDING_WINDOW_DP_DEBUG.Close();
        }
        if MS_DEBUG&logFlags==MS_DEBUG {
            SLIDING_WINDOW_MS_DEBUG.Close();
        }
        if MS_PARALLEL_RESULT_DEBUG&logFlags==MS_PARALLEL_RESULT_DEBUG {
            SLIDING_WINDOW_MS_PARALLEL_RESULT_DEBUG.Close();
        }
    }
}

func TestMain(m *testing.M){
    settings.ReadSettings("testData/testSettings.json");
    testDB=testSetup.SetupDB();
    m.Run();
    testSetup.TeardownDB(&testDB);
}
