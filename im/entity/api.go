package entity

type LoginEntity struct {
	Device   int64  `json:"device"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterEntity struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthorDto login or register result
type AuthorDto struct {
	Token string `json:"token"`
}
