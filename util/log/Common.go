package log;

import (
    "os"
    "log"
)

type LogStatus int;
const (
    Error LogStatus =iota
    Warning
    Deprecation
    Info
    Debug
);

func (l LogStatus)String() string {
    switch l {
        case Error: return "Error";
        case Warning: return "Warning";
        case Deprecation: return "Deprecation";
        case Info: return "Info";
        case Debug: return "Debug";
        default: return "";
    }
}

type Logger struct {
    file string;
    logFile *os.File;
    logger *log.Logger;
    Log func(message string, args ...any);
};

func NewLog(status LogStatus, file string) Logger {
    rv:=Logger{ file: file };
    var err error=nil;
    if rv.logFile,err=os.OpenFile(
        file,
        os.O_TRUNC | os.O_CREATE | os.O_WRONLY,
        0644,
    ); err==nil {
        rv.logger=log.New(rv.logFile,status.String(),log.LstdFlags);
    }
    rv.Log=func(message string, args ...any){
        rv.logger.Printf(message,args...);
    }
    return rv;
}

func NewBlankLog() Logger {
    return Logger{
        Log: func(message string, args ...any){},
    };
}

func (l *Logger)Close(){
    if l.logFile!=nil {
        l.logFile.Close();
    }
}

func (l *Logger)Clear() error {
    if len(l.file)>0 {
        return os.Truncate(l.file,0);
    }
    return LogFileNotSpecified("Nothing to clear.");
}
