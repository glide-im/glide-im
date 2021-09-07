package api

import (
	"context"
	"github.com/smallnest/rpcx/share"
	"go_im/im/api"
	"go_im/im/message"
	"go_im/service/api/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	*rpc.BaseServer
}

func (s *Server) Handle(ctx context.Context, r *pb.HandleRequest, _ *emptypb.Empty) error {
	res := ctx.Value(share.ResMetaDataKey).(map[string]string)
	res["from_server"] = "value_2"
	rpc.LogContext(ctx)

	api.Handle(r.Uid, &message.Message{
		Seq:    r.GetMessage().GetSeq(),
		Action: message.Action(r.GetMessage().GetAction()),
		Data:   r.GetMessage().GetData(),
	})
	return nil
}

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
	}
	s.Register(options.Name, s)
	return s
}
