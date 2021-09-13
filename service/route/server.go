package route

import (
	"context"
	"go_im/service/pb"
	"go_im/service/rpc"
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

func (s *Server) Route(ctx context.Context, param *pb.RouteReq, reply *pb.Any) error {
	rt, ok := s.rts[param.SrvId]
	if !ok {
		//
	}

	apiReq := &pb.HandleRequest{
		Uid:     0,
		Message: nil,
	}
	_ = rt.route(ctx, param.Fn, apiReq, reply)
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
	return nil
}
