package model;

import (
    "fmt"
    "time"
    "math"
    "github.com/carmichaeljr/powerlifting-engine/db"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/mathUtil"
)

//Note - THE ORDER OF THE STRUCT FIELDS MUST MATCH THE ORDER OF THE VALUES
//IN THE QUERY. Otherwise the values returned will be all jumbled up.
type dataPoint struct {
    DatePerformed time.Time;
    Sets float64;
    Reps float64;
    Effort float64;
    Intensity float64;
    FatigueIndex float64;
};

//Generates all missing model states for the given client across all exercises
func UpdateModelStates(c *db.CRUD, clientID int) error {
    return nil;
}

//Generates the model state given the date and exercise specified in the
//training log. Uses the training log data as the data that is being predicted,
//which means it needs to have all **VALID** values.
func GenerateModelState(c *db.CRUD, tl *db.TrainingLog) (db.ModelState,error) {
    rv:=db.ModelState{
        ClientID: tl.ClientID,
        ExerciseID: tl.ExerciseID,
        Date: tl.DatePerformed,
    };
    var mse float64=math.Inf(1);
    var curDate time.Time;
    lr:=fatigueAwareModel();
    query:=modelStateQuery(&tl.DatePerformed);
    actualVals,err:=getActualVals(c,tl);
    if err!=nil {
        return rv,formatModelDataError(err);
    }
    if e1:=db.CustomReadQuery(c,query,[]any{
        tl.ExerciseID,tl.ClientID,
    },func(d *dataPoint){
        if !curDate.Equal(d.DatePerformed) {
            mse,err=calcAndSetModelState(&lr,&rv,actualVals,curDate,mse);
            curDate=d.DatePerformed;
        }
        lr.UpdateSummations(map[string]float64{
            "F": d.FatigueIndex, "I": d.Intensity, "R": d.Reps,
            "E": d.Effort, "S": d.Sets,
        });
        //fmt.Println("Date: ",curDate);
    }); e1!=nil {
        return rv,formatModelDataError(e1);
    }
    return rv,formatModelDataError(err);
}

func getActualVals(c *db.CRUD, tl *db.TrainingLog) ([]db.TrainingLog,error) {
    rv:=make([]db.TrainingLog,0);
    err:=db.Read(c,*tl,util.GenFilter(
        false,"DatePerformed","ExerciseID","ClientID",
    ),func(tl *db.TrainingLog){
        rv=append(rv,*tl);
    });
    return rv,err;
}

func calcAndSetModelState(
        lr *mathUtil.LinearReg[float64],
        cur *db.ModelState,
        tl []db.TrainingLog,
        date time.Time,
        oldSe float64) (float64,error) {
    var newSe float64=0.0;
    var diff time.Duration;
    if len(tl)>0 {
        diff=date.Sub(tl[0].DatePerformed);
    }
    diffDays:=int(diff.Hours()/24);
    res,rcond,err:=lr.Run();
    for _,iter:=range(tl) {
        newSe+=mathUtil.SquareError(
            iter.Intensity,getPredFromLinRegResult(res,&iter),
        );
    }
    //fmt.Printf("At %d days the diff is %0.8f",int(diff.Hours()/24),newSe);
    //newSe/=float64(len(tl));
    //fmt.Printf(" and the mse is: %0.8f\n",newSe);
    //fmt.Printf("With A=%0.8f, B=%0.8f, C=%0.8f, D=%0.8f, Eps=%0.8f, Eps2=%0.8f\n",
    //    res.GetConstant(1),res.GetConstant(2),res.GetConstant(3),res.GetConstant(0),
    //    res.GetConstant(4),res.GetConstant(5),
    //);
    if newSe<oldSe && diffDays<=-30 {
        cur.A=math.Max(res.GetConstant(1),0);
        cur.B=math.Max(res.GetConstant(2),0);
        cur.C=math.Max(res.GetConstant(3),0);
        cur.D=math.Max(res.GetConstant(0),0);
        cur.Eps=res.GetConstant(4);
        cur.Eps2=res.GetConstant(5);
        cur.TimeFrame=diffDays;
        cur.Rcond=rcond;
        cur.Difference=newSe;
        return newSe,err;
    }
    return oldSe,err;
}

func modelStateQuery(t *time.Time) string {
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
        ORDER BY TrainingLog.DatePerformed DESC;`,
        t.Format("01/02/2006"),
    );
}

func formatModelDataError(err error) error {
    if err!=nil {
        return util.ModelDataError(fmt.Sprintf(
            "Could not read model data at given date | %s",err,
        ));
    }
    return nil;
}
