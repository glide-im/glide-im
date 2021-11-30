package msg

import (
	"go_im/im/api/apidep"
	"go_im/im/api/router"
	"go_im/im/message"
	"go_im/pkg/db"
	"go_im/pkg/logger"
	"testing"
)

var api = MsgApi{}

func init() {
	db.Init()
	apidep.ClientManager = apidep.MockClientManager{}
}

func getContext(uid int64, device int64) *route.Context {
	return &route.Context{
		Uid:    uid,
		Device: device,
		Seq:    1,
		Action: "",
		R: func(message *message.Message) {
			logger.D("Response=%v", message.Data)
		},
	}
}

func TestMsgApi_GetRecentChatMessages(t *testing.T) {

}

func TestMsgApi_GetRecentSessions(t *testing.T) {
	err := api.GetRecentSessions(getContext(1, 1))
	if err != nil {
		t.Error(err)
	}
}

func TestMsgApi_CreateSession(t *testing.T) {
	err := api.GetOrCreateSession(getContext(1, 1), &SessionRequest{To: 3})
	if err != nil {
		t.Error(err)
	}
}
