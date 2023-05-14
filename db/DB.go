package db;

import (
    "os"
    "fmt"
    "bufio"
    "strings"
    "database/sql"
    "github.com/barbell-math/block/settings"
    customio "github.com/barbell-math/block/util/io"
    customerr "github.com/barbell-math/block/util/err"

    _ "github.com/lib/pq"
)

type DB struct {
    db *sql.DB;
};

func NewDB(host string, port int, name string) (DB,error) {
    var rv DB;
    err:=customerr.ChainedErrorOps(
        func(r ...any) (any,error) {
            return sql.Open("postgres",
                fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
                    os.Getenv("DB_USER"),os.Getenv("DB_PSWD"),host,port,name,
            ));
        },
        func(r ...any) (any,error) { return nil,r[0].(*sql.DB).Ping(); },
        func(r ...any) (any,error) {
            rv.db=r[0].(*sql.DB);
            //rv.db.SetMaxOpenConns(100);
            //rv.db.SetMaxIdleConns(100);
            //rv.db.SetConnMaxLifetime(time.Minute*3);
            return nil,rv.implicitDataConversion(false);
    });
    return rv,err;
}

func (c *DB)RunDataConversion() error {
    return c.implicitDataConversion(false);
}

func (c *DB)implicitDataConversion(check bool) error {
    return customerr.ChainedErrorOps(
        func(r ...any) (any,error) {
            if v,err:=c.getDataVersion(); err!=nil {
                return v,DataVersionNotAvailable;
            } else {
                return v,err;
            }
        }, func(r ...any) (any,error) {
            cont:=true;
            var err error=nil;
            for i:=r[0].(int)+1;
                r[0].(int)>=0 && i<=settings.DataVersion() && err==nil && cont;
                i++ {
                if check {
                    prompt:=fmt.Sprintf(
                        "Moving data from version v%d to v%d, continue",i-1,i,
                    );
                    cont=customio.YNQuestion(prompt);
                }
                if cont {
                    err=c.execDataConversion(i,i-1);
                }
            }
            return nil,err;
    });
}

func (c *DB)execDataConversion(toVersion int, fromVersion int) error {
    return customerr.ChainedErrorOps(
        func(r ...any) (any,error) {
            if f,e:=DataVersionOps[toVersion]; e {
                return f,nil;
            } else {
                return f,NoKnownDataConversion(
                    fmt.Sprintf("From: v%d To: v%d",fromVersion,toVersion),
                );
            }
        }, func(r ...any) (any,error) {
            if err:=r[0].(DataVersionConversion)(c); err!=nil {
                return nil,DataConversion(fmt.Sprintf(
                    "From: v%d To: v%d | %v",fromVersion,toVersion,err,
                ));
            }
            return nil,nil;
        }, func(r ...any) (any,error) {
            return nil,c.setDataVersion(toVersion);
    });
}

func (c *DB)getDataVersion() (int,error) {
    var rv int;
    err:=c.db.QueryRow("SELECT * FROM Version;").Scan(&rv);
    return rv,err;
}
func (c *DB)setDataVersion(v int) error {
    val,err:=c.getDataVersion();
    if err==sql.ErrNoRows {
        _,err=c.db.Exec("INSERT INTO Version(num) VALUES ($1);",v);
    } else if err==nil && val!=v {
        _,err=c.db.Exec("UPDATE Version SET num=$1",v);
    }
    return err;
}

func (c *DB)ResetDB() error {
    err:=c.ExecSQLScript(settings.SQLGlobalInitScript());
    return err;
}

func (c *DB)ExecSQLScript(src string) error {
    var err error=nil;
    var globalInit *os.File=nil;
    if globalInit,err=os.Open(src); err==nil {
        defer globalInit.Close();
        scanner:=bufio.NewScanner(globalInit);
        scanner.Split(customio.Splitter(";"));
        for err==nil && scanner.Scan() {
            _,err=c.db.Exec(strings.TrimSpace(scanner.Text())+";");
        }
    } else {
        return SqlScriptNotFound(fmt.Sprintf("Given file: %s",src));
    }
    return err;
}

func (c *DB)Stats() sql.DBStats {
    return c.db.Stats();
}

func (c *DB)Close(){
    if c.db!=nil {
        c.db.Close();
    }
}
