package route

import (
	"context"
	"errors"
	"go_im/pkg/logger"
	"go_im/service/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	ExtraTag    = "rt_extra_tag"
	ExtraSrvUrl = "rt_extra_srv_url"
	ExtraFrom   = "rt_extra_from"
)

type Server struct {
	*rpc.BaseServer
	rts map[string]*service
}

func NewServer(options *rpc.ServerOptions) *Server {
	s := &Server{
		BaseServer: rpc.NewBaseServer(options),
		rts:        map[string]*service{},
	}
	s.BaseServer.Register(options.Name, s)
	return s
}

func (s *Server) SetTag(ctx context.Context, req *pb.SetTagReq, empty *emptypb.Empty) {

}

func (s *Server) ClearTag(ctx context.Context, req *pb.ClearTagReq, empty *emptypb.Empty) {

}

func (s *Server) Route(ctx context.Context, param *pb.RouteReq, reply *pb.RouteReply) error {
	rt, ok := s.rts[param.SrvId]
	if !ok {
		return errors.New("service not found: srvId=" + param.SrvId)
	}
	reply.Reset()
	reply.Success = true
	reply.Msg = "success"
	reply.Reply = &anypb.Any{}

	p := &pb.RouteReqParam{Data: reply.GetReply()}
	_ = rt.route(ctx, param.Fn, p, reply.GetReply())
	return nil
}

func (s *Server) Register(ctx context.Context, param *pb.RegisterRtReq, _ *emptypb.Empty) error {
	if param.GetRoutePolicy() == 0 {

	}
	if param.GetDiscoveryType() == 1 {

	}
	sv := newService(&rpc.ClientOptions{
		Name:        param.GetSrvName(),
		EtcdServers: param.GetDiscoverySrvUrl(),
	})
	err := sv.BaseClient.Run()
	if err != nil {
		return err
	}
	s.rts[param.SrvId] = sv
	logger.D("Service register: %s", param.SrvName)
	return nil
}
