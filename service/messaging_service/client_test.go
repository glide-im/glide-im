package messaging_service

import (
	"go_im/im/message"
	"go_im/im/messaging"
	"go_im/service"
	"testing"
)

func TestNewClientEtcd(t *testing.T) {

	config, err := service.GetConfig()
	if err != nil {
		t.Error(err)
	}

	err = SetupClient(config)
	if err != nil {
		t.Error(err)
	}

	chatMsg := message.NewChatMessage(1, 1, 1, 1, 1, "h", 1)
	err = messaging.HandleMessage(1, 1, message.NewMessage(1, "a", &chatMsg))

	if err != nil {
		t.Error(err)
	}
}
