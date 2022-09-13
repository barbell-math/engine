package db;

import (
    "os"
    "fmt"
    "bufio"
    "database/sql"
    "github.com/carmichaeljr/powerlifting-engine/util"

    _ "github.com/lib/pq"
)

type CRUD struct {
    db *sql.DB;
};

func NewCRUD(host string, port int, user string, name string) (CRUD,error) {
    var rv CRUD;
    var err error;
    rv.db,err=sql.Open("postgres",
        fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
            user,os.Getenv("DB_PSWD"),host,port,name,
        ),
    );
    if err==nil {
        err=rv.db.Ping();
    }
    if err==nil {
        err=rv.implicitDataConversion(true);
    }
    return rv,err;
}

func (c *CRUD)RunDataConversion() error {
    return c.implicitDataConversion(false);
}

func (c *CRUD)implicitDataConversion(check bool) error {
    cont:=true;
    dbDataVersion,err:=c.getDataVersion();
    if err!=nil {
        return util.DataVersionNotAvailable;
    }
    for i:=dbDataVersion+1;
        dbDataVersion!=-1 && i<=CURRENT_DATA_VERSION && err==nil && cont; i++ {
        if check {
            prompt:=fmt.Sprintf(
                "Moving data from version v%d to v%d, continue",i-1,i,
            );
            cont=util.YNQuestion(prompt);
        }
        if cont {
            err=c.execDataConversion(i,i-1);
        }
    }
    return err;
}

func (c *CRUD)execDataConversion(toVersion int, fromVersion int) error {
    var err error=nil;
    if f,exist:=DataVersionOps[toVersion]; exist {
        if err=f(c); err==nil {
            err=c.setDataVersion(toVersion);
        } else {
            err=util.DataConversion(
                fmt.Sprintf("From: v%d To: v%d",fromVersion,toVersion),
            );
        }
    } else {
        err=util.NoKnownDataConversion(
            fmt.Sprintf("From: v%d To: v%d",fromVersion,toVersion),
        );
    }
    return err;
}

func (c *CRUD)getDataVersion() (int,error) {
    var rv int;
    err:=c.db.QueryRow("SELECT * FROM Version;").Scan(&rv);
    return rv,err;
}
func (c *CRUD)setDataVersion(v int) error {
    _,err:=c.db.Exec("UPDATE Version SET num=$1",v);
    return err;
}
//TODO - check to make sure one is not already present
func (c *CRUD)addDataVersion(v int) error {
    _,err:=c.db.Exec("INSERT INTO Version(num) VALUES ($1);",v);
    return err;
}

//TODO - read file locations from json file??
func (c *CRUD)ResetDB() error {
    err:=c.execSQLScript("./sql/globalInit.sql");
    return err;
}

func (c *CRUD)execSQLScript(src string) error {
    var err error=nil;
    var globalInit *os.File=nil;
    if globalInit,err=os.Open(src); err==nil {
        defer globalInit.Close();
        scanner:=bufio.NewScanner(globalInit);
        scanner.Split(util.Splitter(";"));
        for scanner.Scan() {
            _,err=c.db.Exec(scanner.Text()+";");
        }
    } else {
        return util.SqlScriptNotFound(fmt.Sprintf("Given file: %s",src));
     }
    return err;
}

func (c *CRUD)CreateExerciseType(e ExerciseType) (int,error) {
    var rv int;
    stmt:="INSERT INTO ExerciseTypes(_type,description) VALUES ($1,$2) RETURNING id;";
    err:=c.db.QueryRow(stmt,e._type,e.description).Scan(&rv);
    return rv,err;
}

func (c *CRUD)CreateExerciseFocus(e ExerciseFocus) (int,error) {
    var rv int;
    stmt:="INSERT INTO ExerciseFocus(focus) VALUES ($1) RETURNING id;";
    err:=c.db.QueryRow(stmt,e.focus).Scan(&rv);
    return rv,err;
}

func (c *CRUD)CreateExercise(e Exercise) (int,error) {
    var rv int;
    stmt:="INSERT INTO Exercises(name,typeID,focusID) VALUES ($1,$2,$3) RETURNING id;";
    err:=c.db.QueryRow(stmt,e.name,e.typeID,e.focusID).Scan(&rv);
    return rv,err;
}

func (c *CRUD)Close(){
    if c.db!=nil {
        c.db.Close();
    }
}
