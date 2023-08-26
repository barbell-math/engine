package stateGenerator;

import (
    "time"
    "testing"
	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/util/dataStruct"
	"github.com/barbell-math/block/util/io/csv"
	"github.com/barbell-math/block/util/algo/iter"
	potSurf "github.com/barbell-math/block/model/potentialSurface"
)

func TestBook_GenerateBookData(t *testing.T) {
    db.DeleteAll[db.ModelState](&testDB);
    timeFrame:=dataStruct.Pair[int,int]{A: 1, B: 5000};
    window:=dataStruct.Pair[int,int]{A: 1, B: 30};
    sw,_:=NewSlidingWindowStateGen(timeFrame,window,1);
    c,_:=db.GetClientByEmail(&testDB,"testing@testing.com");
    // Earilest data point is 8/10/2021, this date is small enough to get all values
    sw.GenerateClientModelStates(&testDB,c,time.Date(
        2020,time.Month(1),1,0,0,0,0,time.UTC,
    ),func() []potSurf.Surface {
        return []potSurf.Surface{
            potSurf.NewBasicSurface().ToGenericSurf(),
            potSurf.NewVolumeBaseSurface().ToGenericSurf(),
        };
    });
}

func TestBook_SlidingWindow(t *testing.T) {
    csv.Flatten(csv.StructToCSV(iter.Map(db.CustomReadQuery[db.ModelState](&testDB,
        "SELECT * FROM ModelState ORDER BY Date;", []any{},
    ),func(index int, val *db.ModelState) (db.ModelState, error) {
        return *val,nil;
    }),true,"01/02/2006"),",").ToFile(
        "../../../data/generatedData/Client1.ms.csv",true,
    );
}

func TestBook_VolumeSkew(t *testing.T){
//    //generate all model states
//    //iterate over all model states, get volume skew, get approximate
//    //put in csv file with date,client,exercise,sets?,reps?
    type tempStruct struct {
        Date time.Time;
        ClientID int;
        ExerciseID int;
        PotentalSurfaceID int;
        StateGeneratorID int;
        Eps,Eps1,Eps2,Eps3 float64;
        Eps4,Eps5,Eps6,Eps7 float64;
        Sets float64;
        Reps float64;
        Effort float64;
        InterWorkoutFatigue int;
        InterExerciseFatigue int;
        VolumeSkew float64;
        ApproxVolumeSkew float64;
    };
    csv.Flatten(csv.StructToCSV[tempStruct](iter.Map[*tempStruct,tempStruct](db.CustomReadQuery[tempStruct](&testDB,
        `SELECT ModelState.Date,
                ModelState.ClientID,
                ModelState.ExerciseID,
                ModelState.PotentialSurfaceID,
                ModelState.StateGeneratorID,
                ModelState.Eps,
                ModelState.Eps1,
                ModelState.Eps2,
                ModelState.Eps3,
                ModelState.Eps4,
                ModelState.Eps5,
                ModelState.Eps6,
                ModelState.Eps7,
                TrainingLog.Sets,
                TrainingLog.Reps,
                TrainingLog.Effort,
                TrainingLog.InterWorkoutFatigue,
                TrainingLog.InterExerciseFatigue,
                0.0 AS VolumeSkew,
                0.0 AS ApproxVolumeSkew
        FROM ModelState 
        JOIN TrainingLog
        ON ModelState.ClientID=TrainingLog.ClientID AND
            ModelState.Date=TrainingLog.DatePerformed AND
            ModelState.ExerciseID=TrainingLog.ExerciseID
        WHERE ModelState.Date>TO_DATE('05/01/2022','MM/DD/YYYY')
        ORDER BY Date DESC;`,[]any{},
    ), func(index int, val *tempStruct) (tempStruct, error) {
        predictor:=potSurf.PredictorFromSurfaceId(potSurf.PotentialSurfaceId(val.PotentalSurfaceID));
        tmpMs:=db.ModelState{
            Eps: val.Eps, Eps1: val.Eps1, Eps2: val.Eps2, Eps3: val.Eps3,
            Eps4: val.Eps4, Eps5: val.Eps5, Eps6: val.Eps6, Eps7: val.Eps7,
        };
        tmpTl:=db.TrainingLog{
            InterWorkoutFatigue: val.InterWorkoutFatigue,
            InterExerciseFatigue: val.InterExerciseFatigue,
            Effort: val.Effort,
        };
        val.VolumeSkew=predictor.VolumeSkew(&tmpMs,&tmpTl);
        val.ApproxVolumeSkew=predictor.VolumeSkewApprox(&tmpMs,&tmpTl);
        return *val,nil;
    }),true,"01/02/2006"),",").ToFile(
        "../../../data/generatedData/Client1.volSkew.csv",true,
    );
}
