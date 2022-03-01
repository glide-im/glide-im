package api_service

import (
	"go_im/im/message"
)

type GatewayInterface struct {
}

func (g GatewayInterface) ClientSignIn(oldUid int64, uid int64, device int64) {

}

func (g GatewayInterface) ClientLogout(uid int64, device int64) {

}

func (g GatewayInterface) EnqueueMessage(uid int64, device int64, message *message.Message) {

}
