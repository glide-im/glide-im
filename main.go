package main

import (
	"go_im/im"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/group"
	"go_im/pkg/db"
	_ "net/http/pprof"
)

func init() {
	db.Init()
}

func main() {
	server := im.NewServer(im.Options{
		SvrType:       im.WebSocket,
		ApiImpl:       api.NewApiRouter(),
		ClientMgrImpl: client.NewClientManager(),
		GroupMgrImpl:  group.NewGroupManager(),
	})
	server.Serve("0.0.0.0", 8080)
}
