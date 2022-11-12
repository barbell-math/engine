package model;

import (
    "fmt"
    "time"
    "math"
    "github.com/carmichaeljr/powerlifting-engine/db"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/mathUtil"
)

type PredictionState struct {
    window int;
    minTimeFrame int;
    maxTimeFrame int;
    optimalMs db.ModelState;
    startDate time.Time;
    actualVals []db.TrainingLog;
    lr mathUtil.LinearReg[float64];
};

//The lr and actual values are not created until the generate prediction method
//in order to make sure the slice values are unique. (ie. It makes sure multiple
//threads don't access the same slices because they are not deep copied even
//though the method does not have a pointer receiver.)
func NewPredictionState(
        minTimeFrame int,
        maxTimeFrame int,
        window int) PredictionState {
    return PredictionState{
        window: -mathUtil.Abs(window),
        minTimeFrame: -mathUtil.Abs(minTimeFrame),
        maxTimeFrame: -mathUtil.Abs(maxTimeFrame),
        optimalMs: db.ModelState{
            Mse: math.Inf(1),
        },
    };
}

//Generates all missing model states for the given client across all exercises.
//This method does not have a pointer receiver because it is meant to be run in
//parallel. Not having a pointer receiver ensures values are copied to the new
//thread.
func (p PredictionState)UpdateModelStates(
        c *db.CRUD,
        clientID int,
        rv chan<- error) {
    type needsModelState struct {
        Date time.Time;
        ExerciseID int;
    };
    rv<-formatModelDataError(db.CustomReadQuery(c,p.msMissingQuery(),[]any{
        clientID,
    }, func(n *needsModelState){
        fmt.Printf("%+v\n",n);
    }));
}

type ModelStateGenerationRes struct {
    Ms db.ModelState;
    Err error;
};

//Generates the model state given the date and exercise specified in the
//training log. The following values need to be valid in the training log arg:
//  - DatePerformed
//  - ClientID
//  - ExerciseID
//This method does not have a pointer receiver because it is meant to be run in
//parallel. Not having a pointer receiver ensures values are copied to the new
//thread.
func (p PredictionState)GenerateModelState(
        c *db.CRUD,
        tl db.TrainingLog,
        ch chan<- ModelStateGenerationRes){
    var curDate time.Time;
    p.lr=fatigueAwareModel();
    p.startDate=tl.DatePerformed;
    p.optimalMs.ClientID=tl.ClientID;
    p.optimalMs.ExerciseID=tl.ExerciseID;
    p.optimalMs.Date=tl.DatePerformed;
    rv:=ModelStateGenerationRes{ Err: nil };
    if rv.Err=p.getActualVals(c,&tl); rv.Err!=nil {
        ch <- rv;
        return;
    }
    rv.Err=db.CustomReadQuery(c,p.msStateSelectionQuery(&tl.DatePerformed),[]any{
        tl.ExerciseID,tl.ClientID,
    },func(d *dataPoint){
        if !curDate.Equal(d.DatePerformed) {
            rv.Err=p.calcAndSetModelState(curDate);
            curDate=d.DatePerformed;
        }
        p.lr.UpdateSummations(map[string]float64{
            "F": d.FatigueIndex, "I": d.Intensity, "R": d.Reps,
            "E": d.Effort, "S": d.Sets,
        });
    });
    rv.Err=formatModelDataError(rv.Err);
    rv.Ms=p.optimalMs;
    ch <- rv;
    return;
}

func (p *PredictionState)getActualVals(
        c *db.CRUD,
        tl *db.TrainingLog) error {
    p.actualVals=make([]db.TrainingLog,0);
    return formatModelDataError(db.CustomReadQuery(c,fmt.Sprintf(`
        SELECT *
        FROM TrainingLog
        WHERE ClientID=$1
            AND ExerciseID=$2
            AND DatePerformed<='%s'::date
            AND DatePerformed>='%s'::date`,
        tl.DatePerformed.Format("01/02/2006"),
        tl.DatePerformed.AddDate(0, 0, p.window).Format("01/02/2006"),
    ), []any{
        tl.ClientID,tl.ExerciseID,
    },func(iterTl *db.TrainingLog){
        p.actualVals=append(p.actualVals,*iterTl);
    }));
}

func (p *PredictionState)calcAndSetModelState(date time.Time) error {
    var newMse float64=0.0;
    var diff time.Duration;
    diff=date.Sub(p.startDate);
    diffDays:=int(diff.Hours()/24);
    res,rcond,err:=p.lr.Run();
    for _,iter:=range(p.actualVals) {
        newMse+=mathUtil.SquareError(
            iter.Intensity,getPredFromLinRegResult(res,&iter),
        );
    }
    //fmt.Printf("At %d days the diff is %0.8f",int(diff.Hours()/24),newSe);
    newMse/=float64(len(p.actualVals));
    //fmt.Printf(" and the mse is: %0.8f\n",newMse);
    //fmt.Printf("With A=%0.8f, B=%0.8f, C=%0.8f, D=%0.8f, Eps=%0.8f, Eps2=%0.8f\n",
    //    res.GetConstant(1),res.GetConstant(2),res.GetConstant(3),res.GetConstant(0),
    //    res.GetConstant(4),res.GetConstant(5),
    //);
    if newMse<p.optimalMs.Mse && diffDays<=p.minTimeFrame {
        p.optimalMs.A=math.Max(res.GetConstant(1),0);
        p.optimalMs.B=math.Max(res.GetConstant(2),0);
        p.optimalMs.C=math.Max(res.GetConstant(3),0);
        p.optimalMs.D=math.Max(res.GetConstant(0),0);
        p.optimalMs.Eps=res.GetConstant(4);
        p.optimalMs.Eps2=res.GetConstant(5);
        p.optimalMs.TimeFrame=diffDays;
        p.optimalMs.Rcond=rcond;
        p.optimalMs.Mse=newMse;
    }
    return err;
}

func (p *PredictionState)msStateSelectionQuery(t *time.Time) string {
    return fmt.Sprintf(`SELECT TrainingLog.DatePerformed,
            TrainingLog.Sets,
            TrainingLog.Reps,
            TrainingLog.Effort,
            TrainingLog.Intensity,
            TrainingLog.FatigueIndex
        FROM TrainingLog
        WHERE TrainingLog.ExerciseID=$1
            AND ClientID=$2
            AND TrainingLog.DatePerformed<'%s'::date
            AND TrainingLog.DatePerformed>'%s'::date
        ORDER BY TrainingLog.DatePerformed DESC;`,
        t.Format("01/02/2006"),
        t.AddDate(0, 0, p.maxTimeFrame).Format("01/02/2006"),
    );
}

func (p *PredictionState)msMissingQuery() string {
    return `SELECT TrainingLog.DatePerformed,
        TrainingLog.ExerciseID
    FROM TrainingLog
    LEFT JOIN ModelState
    ON TrainingLog.ExerciseID=ModelState.ExerciseID
        AND ModelState.ClientID=TrainingLog.ClientID
        AND TrainingLog.DatePerformed=ModelState.Date
    JOIN Exercise
    ON Exercise.Id=TrainingLog.ExerciseID
    JOIN ExerciseType
    ON ExerciseType.Id=Exercise.TypeID
    WHERE TrainingLog.ClientID=$1
        AND ModelState.Id IS NULL
        AND (ExerciseType.T='Main Compound'
        OR ExerciseType.T='Main Compound Accessory')
    GROUP BY TrainingLog.DatePerformed,
        TrainingLog.ExerciseID;`;
}

func formatModelDataError(err error) error {
    if err!=nil {
        return util.ModelDataError(fmt.Sprintf(
            "Could not read model data at given date. | %s",err,
        ));
    }
    return nil;
}
