package logger

import (
	"go.uber.org/zap"
)

var Zap *zap.Logger

var sugar *zap.SugaredLogger

func init() {
	var err error
	Zap, err = zap.NewDevelopment(
		zap.Development(),
		zap.AddCallerSkip(1),
		zap.WithCaller(true),
		zap.AddCaller(),
	)
	if err != nil {
		panic(err)
	}
	sugar = Zap.Sugar()
}

func E(format string, logs ...interface{}) {
	sugar.Errorf(format, logs...)
}

func I(format string, args ...interface{}) {
	sugar.Infof(format, args...)
}

func D(format string, args ...interface{}) {
	sugar.Debugf(format, args...)
}

func W(format string, args ...interface{}) {
	sugar.Warnf(format, args)
}

func ErrE(msg string, e error) {
	Zap.Error(msg, zap.Error(e))
}

func ErrStr(msg string, k string, v string) {
	Zap.Error(msg, zap.String(k, v))
}

func ErrInt(msg string, k string, v int64) {
	Zap.Error(msg, zap.Int64(k, v))
}

func DebugStr(msg string, k string, v string) {
	Zap.Debug(msg, zap.String(k, v))
}
