package model;

import (
    "time"
    "testing"
	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/util/dataStruct"
	"github.com/barbell-math/block/util/io/csv"
	"github.com/barbell-math/block/util/algo/iter"
)

func TestSaveAllModelStates(t *testing.T) {
    timeFrame:=dataStruct.Pair[int,int]{A: 1, B: 5000};
    window:=dataStruct.Pair[int,int]{A: 1, B: 30};
    sw,_:=NewSlidingWindowStateGen(timeFrame,window,10);
    c,_:=db.GetClientByEmail(&testDB,"testing@testing.com")
    // Earilest data point is 8/10/2021, this date is small enough to get all values
    sw.GenerateClientModelStates(&testDB,c,time.Date(
        2020,time.Month(1),1,0,0,0,0,time.UTC),
    );
    csv.Flatten(csv.StructToCSV(iter.Map(db.CustomReadQuery[db.ModelState](&testDB,
        "SELECT * FROM ModelState ORDER BY Date;",
        []any{},
    ),func(index int, val *db.ModelState) (db.ModelState, error) {
        return *val,nil;
    }),true,"01/02/2006"),",").ToFile(
        "../../data/generatedData/Client1AllModelStates.csv",true,
    );
}
