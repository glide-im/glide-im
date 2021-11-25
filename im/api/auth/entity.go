package auth

type AuthTokenRequest struct {
	Token  string
	Device int64
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

// AuthResponse login or register result
type AuthResponse struct {
	Token string
	Uid   int64
}
