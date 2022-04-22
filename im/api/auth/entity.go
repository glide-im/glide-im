package auth

type AuthTokenRequest struct {
	Token string
}

type SignInRequest struct {
	Device   int64
	Account  string
	Password string
}

type LogoutRequest struct {
}

type RegisterRequest struct {
	Account  string
	Password string
}

type GuestRegisterRequest struct {
	Avatar   string
	Nickname string
}

// AuthResponse login or register result
type AuthResponse struct {
	Token   string
	Uid     int64
	Servers []string
}
