package bookDataGen;

import (
    "fmt"
    "time"
    "testing"
    stdMath "math"
	"github.com/barbell-math/engine/db"
	"github.com/barbell-math/engine/util/dataStruct"
	"github.com/barbell-math/engine/util/io/csv"
	"github.com/barbell-math/engine/util/algo/iter"
	math "github.com/barbell-math/engine/util/math/numeric"
	potSurf "github.com/barbell-math/engine/model/potentialSurface"
    stateGen "github.com/barbell-math/engine/model/stateGenerator" 
    "github.com/barbell-math/engine/model" 
)

func TestBook_GenerateSlidingWindowData(t *testing.T) {
    db.DeleteAll[db.ModelState](&testDB);
    timeFrame:=dataStruct.Pair[int,int]{A: 1, B: 5000};
    window:=dataStruct.Pair[int,int]{A: 1, B: 30};
    sw,_:=stateGen.NewSlidingWindowStateGen(timeFrame,window,1);
    c,_:=db.GetClientByEmail(&testDB,"one");
    // Earilest data point is 8/10/2021, this date is small enough to get all values
    sw.GenerateClientModelStates(&testDB,c,time.Date(
        2020,time.Month(1),1,0,0,0,0,time.UTC,
    ),func() []potSurf.Surface {
        return []potSurf.Surface{
            potSurf.NewBasicSurface().ToGenericSurf(),
            potSurf.NewVolumeBaseSurface().ToGenericSurf(),
        };
    });
    db.ReadAll[db.TrainingLog](&testDB).ForEach(
    func(index int, val *db.TrainingLog) (iter.IteratorFeedback, error) {
        if pred,err:=model.GeneratePrediction(&testDB,
            val,
            stateGen.SlidingWindowStateGenId,
            potSurf.BasicSurfaceId,
        ); err==nil {
            db.Create(&testDB,pred);
        }
        return iter.Continue,nil;
    });
}

func TestBook_SaveGeneratedData(t *testing.T) {
    csv.Flatten(csv.StructToCSV(iter.Map(db.CustomReadQuery[db.ModelState](&testDB,
        fmt.Sprintf(`SELECT * 
            FROM ModelState
            WHERE StateGeneratorID=%d 
                AND ClientID=%d 
            ORDER BY Date;`,
        stateGen.SlidingWindowStateGenId,1),
        []any{},
    ),func(index int, val *db.ModelState) (db.ModelState, error) {
        return *val,nil;
    }),true,"01/02/2006"),",").ToFile(
        "../../data/generatedData/Client1.ms.slidingWindow.basic.csv",true,
    );
    csv.Flatten(csv.StructToCSV(iter.Map(db.CustomReadQuery[db.Prediction](&testDB,
        fmt.Sprintf(`SELECT Prediction.*
            FROM Prediction 
            JOIN TrainingLog 
            ON TrainingLog.Id=Prediction.TrainingLogID
            WHERE StateGeneratorID=%d 
                AND ClientID=%d 
            ORDER BY DatePerformed;`,
        stateGen.SlidingWindowStateGenId,1),
        []any{},
    ),func(index int, val *db.Prediction) (db.Prediction, error) {
        return *val,nil;
    }),true,"01/02/2006"),",").ToFile(
        "../../data/generatedData/Client1.pred.slidingWindow.basic.csv",true,
    );
}

func TestBook_SaveSlidingWindowBasicSurfaceEstimated1RM(t *testing.T){
    type temp struct {
        DatePerformed time.Time;
        Exercise string;
        Effort float32;
        PredictedIntensity float32;
        Intensity float32;
        Difference float32;
    };
    csv.Flatten(csv.StructToCSV(iter.Map(db.CustomReadQuery[temp](&testDB,
        fmt.Sprintf(`SELECT 
                TrainingLog.DatePerformed,
                Exercise.Name,
                TrainingLog.Effort,
                Prediction.IntensityPred,
                TrainingLog.Intensity,
                0.0 AS Difference
            FROM Prediction 
            JOIN TrainingLog 
            ON TrainingLog.Id=Prediction.TrainingLogID
            JOIN Exercise
            ON TrainingLog.ExerciseID=Exercise.Id
            WHERE (StateGeneratorID=%d
                AND PotentialSurfaceID=%d
                AND ClientID=%d 
                AND Sets=1
                AND Reps=1
                AND Intensity>0)
                AND (Intensity>=1 OR Effort=10)
            ORDER BY DatePerformed;`,
        stateGen.SlidingWindowStateGenId,potSurf.BasicSurfaceId,1),
        []any{},
    ),func(index int, val *temp) (temp, error) {
        val.Difference=float32(stdMath.Abs(float64(val.Intensity-val.PredictedIntensity)));
        return *val,nil;
    }).Filter(func(index int, val temp) bool {
        return val.PredictedIntensity>0.1;
    }),true,"01/02/2006"),",").ToFile(
        "../../data/generatedData/Client1.1RMPred.slidingWindow.basic.csv",true,
    );
}

// func TestBook_SaveVolumeSkew(t *testing.T){
//     type tempStruct struct {
//         Date time.Time;
//         ClientID int;
//         ExerciseID int;
//         PotentalSurfaceID int;
//         StateGeneratorID int;
//         Eps,Eps1,Eps2,Eps3 float64;
//         Eps4,Eps5,Eps6,Eps7 float64;
//         Sets float64;
//         Reps float64;
//         Effort float64;
//         InterWorkoutFatigue int;
//         InterExerciseFatigue int;
//         VolumeSkew float64;
//         ApproxVolumeSkew float64;
//     };
//     csv.Flatten(csv.StructToCSV[tempStruct](iter.Map[*tempStruct,tempStruct](db.CustomReadQuery[tempStruct](&testDB,
//         `SELECT ModelState.Date,
//                 ModelState.ClientID,
//                 ModelState.ExerciseID,
//                 ModelState.PotentialSurfaceID,
//                 ModelState.StateGeneratorID,
//                 ModelState.Eps,
//                 ModelState.Eps1,
//                 ModelState.Eps2,
//                 ModelState.Eps3,
//                 ModelState.Eps4,
//                 ModelState.Eps5,
//                 ModelState.Eps6,
//                 ModelState.Eps7,
//                 TrainingLog.Sets,
//                 TrainingLog.Reps,
//                 TrainingLog.Effort,
//                 TrainingLog.InterWorkoutFatigue,
//                 TrainingLog.InterExerciseFatigue,
//                 0.0 AS VolumeSkew,
//                 0.0 AS ApproxVolumeSkew
//         FROM ModelState 
//         JOIN TrainingLog
//         ON ModelState.ClientID=TrainingLog.ClientID AND
//             ModelState.Date=TrainingLog.DatePerformed AND
//             ModelState.ExerciseID=TrainingLog.ExerciseID
//         WHERE ModelState.Date>TO_DATE('05/01/2021','MM/DD/YYYY')
//         ORDER BY Date ASC;`,[]any{},
//     ), func(index int, val *tempStruct) (tempStruct, error) {
//         fmt.Print("Working: ",index,"\r");
//         predictor:=potSurf.CalculationsFromSurfaceId(potSurf.PotentialSurfaceId(val.PotentalSurfaceID));
//         tmpMs:=db.ModelState{
//             Eps: val.Eps, Eps1: val.Eps1, Eps2: val.Eps2, Eps3: val.Eps3,
//             Eps4: val.Eps4, Eps5: val.Eps5, Eps6: val.Eps6, Eps7: val.Eps7,
//         };
//         tmpTl:=db.TrainingLog{
//             InterWorkoutFatigue: val.InterWorkoutFatigue,
//             InterExerciseFatigue: val.InterExerciseFatigue,
//             Effort: val.Effort,
//         };
//         val.VolumeSkew=predictor.VolumeSkew(&tmpMs,&tmpTl);
//         val.ApproxVolumeSkew=predictor.VolumeSkewApprox(&tmpMs,&tmpTl);
//         return *val,nil;
//     }),true,"01/02/2006"),",").ToFile(
//         "../../data/generatedData/Client1.volSkew.slidingWindow.basic.csv",true,
//     );
//     fmt.Println();
// }

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
