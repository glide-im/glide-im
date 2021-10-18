package main

import (
	"go_im/im"
	"go_im/pkg/db"
	_ "net/http/pprof"
)

func init() {
	db.Init()
}

func main() {
	server := im.NewServer(im.Options{
		SvrType:       im.WebSocket,
		ApiImpl:       im.NewApiRouter(),
		ClientMgrImpl: im.NewClientManager(),
		GroupMgrImpl:  im.NewGroupManager(),
	})
	server.Serve("0.0.0.0", 8080)
}
