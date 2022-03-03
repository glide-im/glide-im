package api_service

import (
	"go_im/im/message"
)

type GatewayInterface struct {
}

func (g GatewayInterface) ClientSignIn(oldUid int64, uid int64, device int64) error {
	return nil
}

func (g GatewayInterface) ClientLogout(uid int64, device int64) error {
	return nil
}

func (g GatewayInterface) EnqueueMessage(uid int64, device int64, message *message.Message) error {
	return nil
}
