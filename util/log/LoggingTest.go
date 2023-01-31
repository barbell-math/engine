// +build testFrame

//The build option above turns off logging when the 'testFrame' flag is not set
//  during compile time. This allows for creating production and test releases.

package log;

import (
    "os"
    "log"
)

func NewLog(status LogStatus, file string) Logger {
    var rv Logger;
    var err error;
    if rv.logFile,err=os.OpenFile(
        file,
        os.O_APPEND | os.O_CREATE | os.O_WRONLY,
        0644,
    ); err==nil {
        rv.logger=log.New(rv.logFile,status.String(),log.LstdFlags);
    }
    rv.Log=func(message string, args ...any){
        rv.logger.Printf(message,args...);
    }
    return rv;
}
