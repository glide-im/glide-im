package im

import "fmt"

var logger = Logger{}

type Logger struct {
}

func (l *Logger) E(msg string, err error) {
	fmt.Println(msg)
	fmt.Printf(err.Error())
}

func (l *Logger) I(msg string) {
	fmt.Println(msg)
}

func (l *Logger) D(msg string) {
	fmt.Println(msg)
}

func (l *Logger) W(msg string) {
	fmt.Println(msg)
}
