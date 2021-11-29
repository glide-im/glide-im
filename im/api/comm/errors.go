package comm

import "strconv"

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
	return strconv.FormatInt(int64(e.Code), 10) + "," + e.msg
}

type ErrUnexpected struct {
	Code int
	msg  string
	e    error
}

func NewUnexpectedErr(msg string, origin error) *ErrUnexpected {
	return &ErrUnexpected{
		Code: 1000,
		msg:  msg,
		e:    origin,
	}
}

func NewDbErr(origin error) *ErrUnexpected {
	return &ErrUnexpected{
		Code: 1001,
		msg:  "internal server error",
		e:    origin,
	}
}

func (u *ErrUnexpected) Error() string {
	return strconv.FormatInt(int64(u.Code), 10) + "," + u.msg
}
