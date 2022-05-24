package comm

import (
	"go_im/im/dao/common"
	"runtime"
	"strconv"
)

type ErrApiBiz struct {
	Code int
	msg  string
}

func NewApiBizError(code int, msg string) *ErrApiBiz {
	return &ErrApiBiz{
		Code: code,
		msg:  msg,
	}
}

func (e *ErrApiBiz) Error() string {
	return e.msg
}

type ErrUnexpected struct {
	Code   int
	Msg    string
	Origin error
	Line   string
}

func NewUnexpectedErr(msg string, origin error) *ErrUnexpected {
	return &ErrUnexpected{
		Code:   1000,
		Msg:    msg,
		Origin: origin,
		Line:   getLine(),
	}
}

func NewDbErr(origin error) *ErrUnexpected {
	msg := "internal server error"
	if origin == common.ErrNoRecordFound {
		msg = "not found"
	}
	return &ErrUnexpected{
		Code:   1001,
		Msg:    msg,
		Origin: origin,
		Line:   getLine(),
	}
}

func (u *ErrUnexpected) Error() string {
	return u.Msg
}

func getLine() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return file + ":" + strconv.FormatInt(int64(line), 10)
	}
	return ""
}
