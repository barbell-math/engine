// +build !testFrame

//The build option above turns off logging when the 'testFrame' flag is not set
//  during compile time. This allows for creating production and test releases.

package log;

func NewLog(status LogStatus, file string) Logger {
    return Logger{};
}
