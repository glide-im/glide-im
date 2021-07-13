package im

import "fmt"

var logger = Logger{}

type Logger struct {
}

func (l *Logger) E(msg string, err error) {
	fmt.Println(msg)
	fmt.Printf(err.Error())
}

func (l *Logger) I(format string, args ...string) {
	fmt.Println(fmt.Sprintf(format, args))
}

func (l *Logger) D(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args))
}

func (l *Logger) W(msg string) {
	fmt.Println(msg)
}
