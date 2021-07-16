package entity

type LoginRequest struct {
	Device   int64  `json:"device"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthorResponse login or register result
type AuthorResponse struct {
	Token string `json:"token"`
}
