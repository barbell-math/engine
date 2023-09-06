package bookDataGen;

import (
    "fmt"
    "time"
    "testing"
	"github.com/barbell-math/block/db"
	"github.com/barbell-math/block/util/dataStruct"
	"github.com/barbell-math/block/util/io/csv"
	"github.com/barbell-math/block/util/algo/iter"
	math "github.com/barbell-math/block/util/math/numeric"
	potSurf "github.com/barbell-math/block/model/potentialSurface"
    stateGen "github.com/barbell-math/block/model/stateGenerator" 
)

func TestBook_GenerateSlidingWindowData(t *testing.T) {
    db.DeleteAll[db.ModelState](&testDB);
    timeFrame:=dataStruct.Pair[int,int]{A: 1, B: 5000};
    window:=dataStruct.Pair[int,int]{A: 1, B: 30};
    sw,_:=stateGen.NewSlidingWindowStateGen(timeFrame,window,1);
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

func TestBook_SaveSlidingWindow(t *testing.T) {
    csv.Flatten(csv.StructToCSV(iter.Map(db.CustomReadQuery[db.ModelState](&testDB,
        "SELECT * FROM ModelState ORDER BY Date;", []any{},
    ),func(index int, val *db.ModelState) (db.ModelState, error) {
        return *val,nil;
    }),true,"01/02/2006"),",").ToFile(
        "../../data/generatedData/Client1.ms.csv",true,
    );
}

//func TestBook_SaveVolumeSkew(t *testing.T){
//    type tempStruct struct {
//        Date time.Time;
//        ClientID int;
//        ExerciseID int;
//        PotentalSurfaceID int;
//        StateGeneratorID int;
//        Eps,Eps1,Eps2,Eps3 float64;
//        Eps4,Eps5,Eps6,Eps7 float64;
//        Sets float64;
//        Reps float64;
//        Effort float64;
//        InterWorkoutFatigue int;
//        InterExerciseFatigue int;
//        VolumeSkew float64;
//        ApproxVolumeSkew float64;
//    };
//    csv.Flatten(csv.StructToCSV[tempStruct](iter.Map[*tempStruct,tempStruct](db.CustomReadQuery[tempStruct](&testDB,
//        `SELECT ModelState.Date,
//                ModelState.ClientID,
//                ModelState.ExerciseID,
//                ModelState.PotentialSurfaceID,
//                ModelState.StateGeneratorID,
//                ModelState.Eps,
//                ModelState.Eps1,
//                ModelState.Eps2,
//                ModelState.Eps3,
//                ModelState.Eps4,
//                ModelState.Eps5,
//                ModelState.Eps6,
//                ModelState.Eps7,
//                TrainingLog.Sets,
//                TrainingLog.Reps,
//                TrainingLog.Effort,
//                TrainingLog.InterWorkoutFatigue,
//                TrainingLog.InterExerciseFatigue,
//                0.0 AS VolumeSkew,
//                0.0 AS ApproxVolumeSkew
//        FROM ModelState 
//        JOIN TrainingLog
//        ON ModelState.ClientID=TrainingLog.ClientID AND
//            ModelState.Date=TrainingLog.DatePerformed AND
//            ModelState.ExerciseID=TrainingLog.ExerciseID
//        WHERE ModelState.Date>TO_DATE('05/01/2021','MM/DD/YYYY')
//        ORDER BY Date ASC;`,[]any{},
//    ), func(index int, val *tempStruct) (tempStruct, error) {
//        fmt.Print("Working: ",index,"\r");
//        predictor:=potSurf.CalculationsFromSurfaceId(potSurf.PotentialSurfaceId(val.PotentalSurfaceID));
//        tmpMs:=db.ModelState{
//            Eps: val.Eps, Eps1: val.Eps1, Eps2: val.Eps2, Eps3: val.Eps3,
//            Eps4: val.Eps4, Eps5: val.Eps5, Eps6: val.Eps6, Eps7: val.Eps7,
//        };
//        tmpTl:=db.TrainingLog{
//            InterWorkoutFatigue: val.InterWorkoutFatigue,
//            InterExerciseFatigue: val.InterExerciseFatigue,
//            Effort: val.Effort,
//        };
//        val.VolumeSkew=predictor.VolumeSkew(&tmpMs,&tmpTl);
//        val.ApproxVolumeSkew=predictor.VolumeSkewApprox(&tmpMs,&tmpTl);
//        return *val,nil;
//    }),true,"01/02/2006"),",").ToFile(
//        "../../data/generatedData/Client1.volSkew.csv",true,
//    );
//    fmt.Println();
//}

func Test_BookBasicSurfaceVolumeSkewApproxEps1(t *testing.T){
    bookBasicSurfaceVolumeSkewApprox(func(ms *db.ModelState,
        tl *db.TrainingLog,
        val float64,
    ) {
        ms.Eps1=val;
    },"Eps1",t);
}

func Test_BookBasicSurfaceVolumeSkewApproxEps5(t *testing.T){
    bookBasicSurfaceVolumeSkewApprox(func(ms *db.ModelState,
        tl *db.TrainingLog,
        val float64,
    ) {
        ms.Eps4=val;
        ms.Eps5=0.5;
    },"Eps5",t);
}

func Test_BookBasicSurfaceVolumeSkewApproxEps6(t *testing.T){
    bookBasicSurfaceVolumeSkewApprox(func(ms *db.ModelState,
        tl *db.TrainingLog,
        val float64,
    ) {
        ms.Eps5=val;
    },"Eps6",t);
}

func Test_BookBasicSurfaceVolumeSkewApproxEps7(t *testing.T){
    bookBasicSurfaceVolumeSkewApprox(func(ms *db.ModelState,
        tl *db.TrainingLog,
        val float64,
    ) {
        ms.Eps6=val;
    },"Eps7",t);
}

func bookBasicSurfaceVolumeSkewApprox(
        op func(ms *db.ModelState, tl *db.TrainingLog, val float64),
        fileName string,
        t *testing.T){
    ms:=db.ModelState{
        Eps: 5, Eps1: 0, Eps2: 1, Eps3: 1,
        Eps4: 2, Eps5: 1, Eps6: 1,
    };
    tl:=db.TrainingLog{
        Weight: 0, Sets: 0, Reps: 0, Intensity: 1, Effort: 10,
        InterWorkoutFatigue: 1, InterExerciseFatigue: 1,
    };
    type temp struct {
        Eps1 float64;
        Eps5 float64;
        Eps6 float64;
        Eps7 float64;
        VolSkew float64;
        ApproxVolSkew float64;
    };
    csv.Flatten(csv.StructToCSV[temp](iter.Map(
        math.Range[float64](0,5,0.1),
        func(index int, val float64) (temp,error) {
            fmt.Print("Working: ",index,"\r");
            op(&ms,&tl,val);
            return temp{
                VolSkew: potSurf.BasicSurfaceCalculation.VolumeSkew(&ms,&tl),
                ApproxVolSkew: potSurf.BasicSurfaceCalculation.VolumeSkewApprox(&ms,&tl),
                Eps1: ms.Eps1,
                Eps5: ms.Eps4,
                Eps6: ms.Eps5,
                Eps7: ms.Eps6,
            },nil;
        },
    ),true,"01/02/2006"),",").ToFile(
        fmt.Sprintf("../../data/generatedData/basicSurfVolSkewApprox/%s.csv",fileName),
        true,
    );
    fmt.Println();
}
