package conn

import (
	"github.com/gorilla/websocket"
	"testing"
	"time"
)

func TestWsServer_Start(t *testing.T) {

}

func TestConnect(t *testing.T) {

	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, _, err := dialer.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		t.Error(err)
	}
	go func() {
		for {
			messageType, p, e := conn.ReadMessage()
			if e != nil {
				break
			}
			t.Log(messageType, string(p))
		}

	}()

	order := 0
	for true {
		time.Sleep(time.Second * 3)
		msg := struct {
			Order int
			Msg   string
		}{Order: order, Msg: "msg here"}
		order++
		e := conn.WriteJSON(msg)
		if e != nil {
			t.Error(e)
		}
	}

	time.Sleep(time.Hour)
}
