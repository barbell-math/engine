package db;

import (
    "os"
    "fmt"
    "bufio"
    "errors"
    "strings"
    "database/sql"
    "github.com/carmichaeljr/powerlifting-engine/util"
    "github.com/carmichaeljr/powerlifting-engine/settings"

    _ "github.com/lib/pq"
)

type CRUD struct {
    db *sql.DB;
};

func NewCRUD(host string, port int, name string) (CRUD,error) {
    var rv CRUD;
    err:=util.ChainedErrorOps(
        func(r ...any) (any,error) {
            return sql.Open("postgres",
                fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
                    os.Getenv("DB_USER"),os.Getenv("DB_PSWD"),host,port,name,
            ));
        },
        func(r ...any) (any,error) { return nil,r[0].(*sql.DB).Ping(); },
        func(r ...any) (any,error) {
            rv.db=r[0].(*sql.DB);
            return nil,rv.implicitDataConversion(false);
    });
    return rv,err;
}

func (c *CRUD)RunDataConversion() error {
    return c.implicitDataConversion(false);
}

func (c *CRUD)implicitDataConversion(check bool) error {
    cont:=true;
    return util.ChainedErrorOpsWithCustomErrors([]error{
            util.DataVersionNotAvailable,
        }, func(r ...any) (any,error) {
            return c.getDataVersion();
        }, func(r ...any) (any,error) {
            var err error=nil;
            for i:=r[0].(int)+1;
                r[0].(int)>=0 && i<=settings.DataVersion() && err==nil && cont;
                i++ {
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
            return nil,err;
    });
    //dbDataVersion,err:=c.getDataVersion();
    //if err!=nil {
    //    return util.DataVersionNotAvailable;
    //}
    //for i:=dbDataVersion+1;
    //    dbDataVersion>=0 && i<=CURRENT_DATA_VERSION && err==nil && cont; i++ {
    //    if check {
    //        prompt:=fmt.Sprintf(
    //            "Moving data from version v%d to v%d, continue",i-1,i,
    //        );
    //        cont=util.YNQuestion(prompt);
    //    }
    //    if cont {
    //        err=c.execDataConversion(i,i-1);
    //    }
    //}
    //return err;
}

func (c *CRUD)execDataConversion(toVersion int, fromVersion int) error {
    return util.ChainedErrorOpsWithCustomErrors(
        []error{
            util.NoKnownDataConversion(
                fmt.Sprintf("From: v%d To: v%d",fromVersion,toVersion),
            ), util.DataConversion(
                fmt.Sprintf("From: v%d To: v%d",fromVersion,toVersion),
            ),
        }, func(r ...any) (any,error) {
            if f,e:=DataVersionOps[toVersion]; e {
                return f,nil;
            } else {
                return f,errors.New("");
            }
        }, func(r ...any) (any,error) {
            return nil,r[0].(DataVersionConversion)(c);
        }, func(r ...any) (any,error) {
            return nil,c.setDataVersion(toVersion);
    });
}

func (c *CRUD)getDataVersion() (int,error) {
    var rv int;
    err:=c.db.QueryRow("SELECT * FROM Version;").Scan(&rv);
    return rv,err;
}
func (c *CRUD)setDataVersion(v int) error {
    val,err:=c.getDataVersion();
    if err==sql.ErrNoRows {
        _,err=c.db.Exec("INSERT INTO Version(num) VALUES ($1);",v);
    } else if err==nil && val!=v {
        _,err=c.db.Exec("UPDATE Version SET num=$1",v);
    }
    return err;
}

func (c *CRUD)ResetDB() error {
    err:=c.ExecSQLScript(settings.SQLGlobalInitScript());
    return err;
}

func (c *CRUD)ExecSQLScript(src string) error {
    var err error=nil;
    var globalInit *os.File=nil;
    if globalInit,err=os.Open(src); err==nil {
        defer globalInit.Close();
        scanner:=bufio.NewScanner(globalInit);
        scanner.Split(util.Splitter(";"));
        for err==nil && scanner.Scan() {
            _,err=c.db.Exec(strings.TrimSpace(scanner.Text())+";");
        }
    } else {
        return util.SqlScriptNotFound(fmt.Sprintf("Given file: %s",src));
    }
    return err;
}

func (c *CRUD)Close(){
    if c.db!=nil {
        c.db.Close();
    }
}
