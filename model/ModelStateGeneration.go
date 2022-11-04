package model;

import (
    "fmt"
    "time"
    "github.com/carmichaeljr/powerlifting-engine/db"
    "github.com/carmichaeljr/powerlifting-engine/mathUtil"
)

//Generates the model state given the date and exercise specified in the
//training log. Uses the training log data as the data that is being predicted,
//which means it needs to have all **VALID** values.
func GenerateModelState(c *db.CRUD, tl *db.TrainingLog) (db.ModelState,error) {
    //Note - THE ORDER OF THE STRUCT FIELDS MUST MATCH THE ORDER OF THE VALUES
    //IN THE QUERY. Otherwise the values returned will be all jumbled up.
    type DataPoint struct {
        DatePerformed time.Time;
        Sets float64;
        Reps float64;
        Effort float64;
        Intensity float64;
        FatigueIndex float64;
    };
    rv:=db.ModelState{
        ClientID: tl.ClientID,
        ExerciseID: tl.ExerciseID,
        Date: tl.DatePerformed,
    };
    var err error=nil;
    var curDate time.Time;
    lr:=fatigueAwareModel();
    query:=modelStateQuery(&tl.DatePerformed);
    if e1:=db.CustomReadQuery(c,query,[]any{tl.ExerciseID},func(d *DataPoint){
        if !curDate.Equal(d.DatePerformed) {
            err=calcAndSetModelState(&lr,&rv,tl,d.DatePerformed);
            curDate=d.DatePerformed;
        }
        lr.UpdateSummations(map[string]float64{
            "F": d.FatigueIndex, "I": d.Intensity, "R": d.Reps,
            "E": d.Effort, "S": d.Sets,
        });
        //fmt.Printf("%+v\n",d);
    }); e1!=nil {
        //fmt.Println(e1);
        return rv,e1;
    }
    //fmt.Println(err);
    return rv,err;
}

func calcAndSetModelState(
        lr *mathUtil.LinearReg[float64],
        cur *db.ModelState,
        tl *db.TrainingLog,
        date time.Time) error {
    diff:=date.Sub(tl.DatePerformed);
    diffDays:=int(diff.Hours()/24);
    res,rcond,err:=lr.Run();
    newSe:=mathUtil.SquareError(tl.Intensity,getPredFromLinRegResult(res,tl));
    oldSe:=mathUtil.SquareError(tl.Intensity,IntensityPrediction(cur,tl));
    //fmt.Printf("%+v\n",res);
    //fmt.Println(rcond);
    if newSe<oldSe && diffDays<=-40 {
        cur.A=res.GetConstant(1);
        cur.B=res.GetConstant(2);
        cur.C=res.GetConstant(3);
        cur.D=res.GetConstant(0);
        cur.Eps=res.GetConstant(4);
        cur.Eps2=res.GetConstant(5);
        cur.TimeFrame=diffDays;
        cur.Rcond=rcond;
        cur.Difference=newSe;
    }
    return err;
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
            AND TrainingLog.DatePerformed<'%s'::date
        ORDER BY TrainingLog.DatePerformed DESC;`,
        t.Format("01/02/2006"),
    );
}
