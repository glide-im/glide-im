package im

import (
	"testing"
)

func TestWsServer_Start(t *testing.T) {
	ws := NewWsServer(nil)
	ws.Run()
}

func TestConnect(t *testing.T) {

}
