package stateGenerator;

import (
    "fmt"
    "time"
    "testing"
	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/util/dataStruct"
	"github.com/barbell-math/block/util/io/csv"
	"github.com/barbell-math/block/util/algo/iter"
	potSurf "github.com/barbell-math/block/model/PotentialSurface"
)

func bookData_modelStateGeneratorHelper(surfFactory func() []potSurf.Surface, f string) {
    db.DeleteAll[db.ModelState](&testDB);
    timeFrame:=dataStruct.Pair[int,int]{A: 1, B: 5000};
    window:=dataStruct.Pair[int,int]{A: 1, B: 30};
    sw,_:=NewSlidingWindowStateGen(timeFrame,window,1);
    c,_:=db.GetClientByEmail(&testDB,"testing@testing.com");
    // Earilest data point is 8/10/2021, this date is small enough to get all values
    sw.GenerateClientModelStates(&testDB,c,time.Date(
        2020,time.Month(1),1,0,0,0,0,time.UTC,
    ),surfFactory);
    csv.Flatten(csv.StructToCSV(iter.Map(db.CustomReadQuery[db.ModelState](&testDB,
        "SELECT * FROM ModelState ORDER BY Date;",
        []any{},
    ),func(index int, val *db.ModelState) (db.ModelState, error) {
        return *val,nil;
    }),true,"01/02/2006"),",").ToFile(
        fmt.Sprintf("../../../data/generatedData/%s.csv",f),true,
    );
}

func TestBook_SlidingWindow(t *testing.T) {
    bookData_modelStateGeneratorHelper(func() []potSurf.Surface {
        return []potSurf.Surface{
            potSurf.NewBasicSurface().ToGenericSurf(),
            potSurf.NewVolumeBaseSurface().ToGenericSurf(),
        };
    },"Client1.ms");
}
