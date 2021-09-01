package api

import (
	"go_im/im"
	"go_im/service/rpc"
)

type Server struct {
	*rpc.BaseServer
	*im.ApiRouter
}

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
		ApiRouter:  im.NewApiRouter(),
	}
	_ = s.Register(s)
	return s
}
