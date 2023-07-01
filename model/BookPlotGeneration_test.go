package model;

import (
    "time"
    "testing"
	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/util/dataStruct"
)

func TestGenerateAllModelStates(t *testing.T) {
    timeFrame:=dataStruct.Pair[int,int]{A: 1, B: 5000};
    window:=dataStruct.Pair[int,int]{A: 1, B: 30};
    sw,_:=NewSlidingWindowStateGen(timeFrame,window,10);
    c,_:=db.GetClientByEmail(&testDB,"testing@testing.com")
    // Earilest data point is 8/10/2021, this date is small enough to get all values
    sw.GenerateClientModelStates(&testDB,c,time.Date(
        2020,time.Month(1),1,0,0,0,0,time.UTC),
    );
    db.CustomReadQuery[db.ModelState](&testDB,
        "SELECT * FROM ModelState ORDER BY Date",
        []any{},
    ).ForEach()
}
