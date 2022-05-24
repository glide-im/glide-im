package auth

import (
	"github.com/glide-im/glideim/im/api/apidep"
	route "github.com/glide-im/glideim/im/api/router"
	"github.com/glide-im/glideim/im/message"
	"github.com/glide-im/glideim/pkg/db"
	"github.com/glide-im/glideim/pkg/logger"
	"testing"
)

var authApi = AuthApi{}

func init() {
	db.Init()
	apidep.ClientInterface = apidep.MockClientManager{}
}

func getContext(uid int64, device int64) *route.Context {
	return &route.Context{
		Uid:    uid,
		Device: device,
		Seq:    1,
		Action: "",
		R: func(message *message.Message) {
			logger.D("Response=%v", message)
		},
	}
}

func logErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func TestAuthApi_AuthToken(t *testing.T) {
	err := authApi.AuthToken(getContext(2, 0), &AuthTokenRequest{
		Token: "RN9fXQtAoplDCX8uSiajitgFgCZlrcpX",
	})
	logErr(t, err)
}

func TestAuthApi_Register(t *testing.T) {
	err := authApi.Register(getContext(2, 0), &RegisterRequest{
		Account:  "bb",
		Password: "bb",
	})
	logErr(t, err)
}

func TestAuthApi_SignIn(t *testing.T) {
	err := authApi.SignIn(getContext(2, 0), &SignInRequest{
		Account:  "aa",
		Password: "1234567",
		Device:   1,
	})
	logErr(t, err)
}

func TestAuthApi_Logout(t *testing.T) {
	err := authApi.Logout(getContext(543603, 1))
	logErr(t, err)
}
