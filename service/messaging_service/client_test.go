package messaging_service

import (
	"github.com/glide-im/glideim/im/message"
	"github.com/glide-im/glideim/im/messaging"
	"github.com/glide-im/glideim/service"
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
