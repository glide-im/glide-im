package auth

type AuthTokenReq struct {
	Token  string
	Device int64
}

type LoginRequest struct {
	Device   int64
	Account  string
	Password string
}

type LogoutRequest struct {
	Device  int64
	Account string
	Token   string
}

type AuthRequest struct {
	Token    string
	DeviceId int64
}

type RegisterRequest struct {
	Account  string
	Password string
}

// AuthorResponse login or register result
type AuthorResponse struct {
	Token string
	Uid   int64
}
