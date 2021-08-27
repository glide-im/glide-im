package main

import (
	"go_im/im/client"
	"testing"
	"time"
)

func TestClientServer(t *testing.T) {
	runClientService(TypeClientService)
}

func TestClientClient(t *testing.T) {

	go runClientService(TypeApiService)
	time.Sleep(time.Second * 2)

	online := client.Manager.IsOnline(1)
	t.Log("online:", online)
}
