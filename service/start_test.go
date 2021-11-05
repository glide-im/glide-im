package main

import (
	"github.com/gorilla/websocket"
	"go_im/im/message"
	"testing"
	"time"
)

func TestClientServer(t *testing.T) {
	runClientService(TypeClientService)
}

func TestClientClient(t *testing.T) {

	go runClientService(TypeApiService)
	time.Sleep(time.Second * 2)
}

func TestClientOnline(t *testing.T) {
	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	con, _, err := dialer.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		t.Error(err)
	}
	go func() {
		for {
			messageType, p, e := con.ReadMessage()
			if e != nil {
				t.Log(e)
				break
			}
			t.Log(messageType, string(p))
		}

	}()

	_ = con.WriteJSON(message.NewMessage(1, "api.app.echo", ""))
	time.Sleep(time.Hour)
}
