package messaging

import (
	"github.com/glide-im/glideim/im/auth"
	"github.com/glide-im/glideim/im/client"
	"github.com/glide-im/glideim/im/message"
)

func handleAuth(from int64, device int64, msg *message.Message) {

	t := auth.Token{}
	e := msg.DeserializeData(&t)
	if e != nil {
		resp := message.NewMessage(0, message.ActionApiFailed, "invalid token")
		_ = client.EnqueueMessageToDevice(from, device, resp)
		return
	}
	result, err := auth.Auth(from, device, &t)

	if err == nil {
		resp := message.NewMessage(msg.GetSeq(), message.ActionApiSuccess, result)
		_ = client.SignIn(from, result.Uid, device)
		_ = client.EnqueueMessageToDevice(result.Uid, device, resp)
	} else {
		resp := message.NewMessage(0, message.ActionApiFailed, err.Error())
		_ = client.EnqueueMessageToDevice(from, device, resp)
	}
}
