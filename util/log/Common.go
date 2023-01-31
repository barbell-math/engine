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
    logFile *os.File;
    logger *log.Logger;
    Log func(message string, args ...any);
};

func (l *Logger)Close(){
    l.logFile.Close();
}
