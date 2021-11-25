package api

import (
	"go_im/im/api/apidep"
	"go_im/pkg/db"
	"testing"
)

func TestRunHttpServer(t *testing.T) {
	db.Init()
	apidep.ClientManager = apidep.MockClientManager{}
	_ = RunHttpServer("0.0.0.0", 8080)
}
