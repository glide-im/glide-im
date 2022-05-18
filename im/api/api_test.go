package api

import (
	"go_im/im/message"
	"testing"
)

func TestRunHttpServer(t *testing.T) {
	//db.Init()
	//apidep.ClientInterface = apidep.MockClientManager{}
	//_ = RunHttpServer("0.0.0.0", 8080)

	handle, err := handler.Handle(1, 1, message.NewMessage(1, "api.app.echo", "echo"))
	t.Log(handle, err)
}
