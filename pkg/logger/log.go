package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

func E(msg string, logs ...interface{}) {
	f := strings.Repeat("%v, ", len(logs))
	log("E", msg+":"+fmt.Sprintf(f, logs))
}

func I(format string, args ...interface{}) {
	log("I", fmt.Sprintf(format, args...))
}

func D(format string, args ...interface{}) {
	lg := fmt.Sprintf(format, args...)
	log("D", lg)
}

func W(msg string) {
	log("W", msg)
}

func log(level string, log string) {
	t := time.Now().Format("2006/01/02 15:04:05")
	line, _ := callerInfo()
	fmt.Printf("%s %s %s: %s\n", level, t, line, log)
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
	n := runtime.Callers(4, rpc[:])
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

	}
	return strings.Replace(dir+"/", "\\", "/", -1)
}
