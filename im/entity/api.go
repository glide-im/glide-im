package entity

type LoginEntity struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterEntity struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
