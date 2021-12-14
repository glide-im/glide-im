package comm

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
}

func NewUnexpectedErr(msg string, origin error) *ErrUnexpected {
	return &ErrUnexpected{
		Code:   1000,
		Msg:    msg,
		Origin: origin,
	}
}

func NewDbErr(origin error) *ErrUnexpected {
	return &ErrUnexpected{
		Code:   1001,
		Msg:    "internal server error",
		Origin: origin,
	}
}

func (u *ErrUnexpected) Error() string {
	return u.Msg
}
