package comm

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

var Slog = Logger{}

type Logger struct {
}

func (l *Logger) E(msg string, log ...interface{}) {
	f := strings.Repeat("%v, ", len(log))
	l.log("E", msg+":"+fmt.Sprintf(f, log))
}

func (l *Logger) I(format string, args ...interface{}) {
	return
	l.log("I", fmt.Sprintf(format, args...))
}

func (l *Logger) D(format string, args ...interface{}) {
	lg := fmt.Sprintf(format, args...)
	l.log("D", lg)
}

func (l *Logger) W(msg string) {
	l.log("W", msg)
}

func (l *Logger) log(level string, log string) {
	fmt.Printf("%s: %s\n", level, trace(log))
}

func trace(log string) string {
	if true {
		return log
	}
	line, _ := callerInfo()
	l := strings.ReplaceAll(log, "\n", "\n\t")
	return fmt.Sprintf("%s\n\t%s", line, l)
}

func callerInfo() (string, string) {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(5, rpc[:])
	if n < 1 {
		return "-", "-"
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	filePath := strings.ReplaceAll(frame.File, projectRootPath(), "")
	funcName := strings.Split(frame.Function, ".")[1]
	return fmt.Sprintf("%s:%d %s", filePath, frame.Line, funcName), funcName
}

func projectRootPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir+"/", "\\", "/", -1)
}
