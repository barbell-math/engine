package log;

type LogStatus int;
const (
    Error LogStatus =iota
    Warning
    Deprecation
    Info
    Debug
);

type Logger func(message string, args ...any);

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
