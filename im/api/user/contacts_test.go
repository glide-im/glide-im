package user

import (
	"github.com/glide-im/glideim/im/api/apidep"
	"github.com/glide-im/glideim/im/api/router"
	"github.com/glide-im/glideim/im/message"
	"github.com/glide-im/glideim/pkg/db"
	"github.com/glide-im/glideim/pkg/logger"
	"testing"
)

var api = UserApi{}

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

func TestUserApi_AddContact(t *testing.T) {
	err := api.AddContact(getContext(543603, 1), &AddContacts{
		Uid: 543602,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestUserApi_GetContactList(t *testing.T) {
	err := api.GetContactList(getContext(543603, 1))
	if err != nil {
		t.Error(err)
	}
}
