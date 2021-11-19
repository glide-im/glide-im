package comm

type ErrApiBiz struct {
	Code int
	msg  string
}

func NewApiBizError(code int, msg string) ErrApiBiz {
	return ErrApiBiz{
		Code: code,
		msg:  msg,
	}
}

func (e *ErrApiBiz) Error() string {
	return e.msg
}

type ErrUnexpected struct {
	msg string
	e   error
}

func NewUnexpectedErr(msg string, origin error) ErrUnexpected {
	return ErrUnexpected{
		msg: msg,
		e:   origin,
	}
}

func (u *ErrUnexpected) Error() string {
	return u.e.Error()
}
