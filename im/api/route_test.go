package api

import (
	"encoding/json"
	"go_im/im/entity"
	"testing"
)

type TestData struct {
	Name string
}

func (t *TestData) Validate(data interface{}) error {
	return json.Unmarshal([]byte(data.(string)), t)
}

func TestGroup(t *testing.T) {

	fn := func(info *RequestInfo, ts *TestData) {
		t.Log(info, ts)
	}

	router := NewRouter()
	router.Add(
		Group("api",
			Group("user",
				Route("login", fn),
				Route("register", fn),
			),
			Route("info", fn),
		),
	)

	msg := entity.NewMessage(-1, "api.user.login", &TestData{Name: "1234"})
	err := router.Handle(1, msg)
	t.Log(err, router)
}
