// +build testFrame

//The build option above turns off logging when the 'testFrame' flag is not set
//  during compile time. This allows for creating production and test releases.

package log;

import (
    "os"
    "log"
)

func NewLog(status LogStatus, file string) Logger {
    var logger *log.Logger;
    if f,err:=os.OpenFile(
        file,
        os.O_APPEND | os.O_CREATE | os.O_WRONLY,
        0644,
    ); err==nil {
        defer f.Close();
        logger=log.New(f,status.String(),log.LstdFlags);
    }
    return func(message string, args ...any){
        logger.Printf(message,args...);
    }
}
