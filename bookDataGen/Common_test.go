package bookDataGen;

import (
	"testing"

	"github.com/barbell-math/engine/db"
	"github.com/barbell-math/engine/settings"
	"github.com/barbell-math/engine/model/testSetup"
)

var testDB db.DB;

func TestMain(m *testing.M){
    settings.ReadSettings("testData/testSettings.json");
    testDB=testSetup.SetupDB();
    m.Run();
    testSetup.TeardownDB(&testDB);
}
