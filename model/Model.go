package model;

import (
    "fmt"
    "time"
    "math"
    "github.com/carmichaeljr/powerlifting-engine/db"
    "github.com/carmichaeljr/powerlifting-engine/mathUtil"
)

//The model equation is as follows:
//  I=d-a(s-1)^2*(r-1)^2-b(s-1)^2-c(r-1)^2-eps_1*E-eps_2*F
//Where:
//  d,a,b,c,eps_1,eps_2 are the constants linear reg will find
//  s is sets
//  r is reps
//  E is effort (RPE)
//  F is the fatigue index
func MakeIntensityPrediction(ms *db.ModelState, tl *db.TrainingLog) float64 {
    return (ms.D-
            ms.A*math.Pow(float64(tl.Sets-1),2)*math.Pow(float64(tl.Reps-1),2)-
            ms.B*math.Pow(float64(tl.Sets-1),2)-
            ms.C*math.Pow(float64(tl.Reps-1),2)-
            ms.Eps*tl.Effort-
            ms.Eps2*float64(tl.FatigueIndex));
}

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
    query:=genModelStateQuery(&tl.DatePerformed);
    if e1:=db.CustomReadQuery(c,query,[]any{tl.ExerciseID},func(d *DataPoint){
        if !curDate.Equal(d.DatePerformed) {
            err=calcAndSetModelState(&lr,&rv,tl);
            curDate=d.DatePerformed;
        }
        lr.UpdateSummations(map[string]float64{
            "F": d.FatigueIndex, "I": d.Intensity, "R": d.Reps,
            "E": d.Effort, "S": d.Sets,
        });
        fmt.Printf("%+v\n",d);
    }); e1!=nil {
        fmt.Println(e1);
        return rv,e1;
    }
    fmt.Println(err);
    return rv,err;
}

func calcAndSetModelState(
        lr *mathUtil.LinearReg[float64],
        cur *db.ModelState,
        tl *db.TrainingLog) error {
    //put rcond in database later
    res,_,err:=lr.Run();
    newSe:=mathUtil.SquareError(tl.Intensity,getPred(res,tl));
    oldSe:=mathUtil.SquareError(tl.Intensity,MakeIntensityPrediction(cur,tl));
    if newSe<oldSe {
        //Set constants here... BUT HOW TO GET THEM
    }
    return err;
}

func getPred(
        res mathUtil.LinRegResult[float64],
        tl *db.TrainingLog) float64 {
    rv,_:=res(map[string]float64{
        "F": float64(tl.FatigueIndex), "I": tl.Intensity, "R": float64(tl.Reps),
        "E": tl.Effort, "S": float64(tl.Sets),
    });
    return rv;
}

func genModelStateQuery(t *time.Time) string {
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

//Returns non-standard linear regression for the model according to the
//model equation.
func fatigueAwareModel() mathUtil.LinearReg[float64] {
    return mathUtil.NewLinearReg[float64](fatigueAwareSumOpGen());
}
func fatigueAwareSumOpGen() ([]mathUtil.SummationOp[float64],
        mathUtil.SummationOp[float64]) {
    return []mathUtil.SummationOp[float64]{
        func(vals map[string]float64) (float64,error) {
            s,err:=mathUtil.VarAcc(vals,"S");
            if err!=nil {
                return 0, err;
            }
            r,err:=mathUtil.VarAcc(vals,"R");
            if err!=nil {
                return 0, err;
            }
            return math.Pow(s-1,2)*math.Pow(r-1,2),nil;
        }, func(vals map[string]float64) (float64,error) {
            s,err:=mathUtil.VarAcc(vals,"S");
            if err!=nil {
                return 0, err;
            }
            return math.Pow(s-1,2),nil;
        }, func(vals map[string]float64) (float64,error) {
            r,err:=mathUtil.VarAcc(vals,"R");
            if err!=nil {
                return 0, err;
            }
            return math.Pow(r-1,2),nil;
        }, mathUtil.LinearSummationOp[float64]("E"),
        mathUtil.LinearSummationOp[float64]("F"),
    },mathUtil.LinearSummationOp[float64]("I");
}
