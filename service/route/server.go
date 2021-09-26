package route

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go_im/pkg/logger"
	"go_im/service/pb"
	"go_im/service/rpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	ServiceName = "route"
)

const (
	ExtraTag        = "rt_extra_tag"
	ExtraSrvUrl     = "rt_extra_srv_url"
	ExtraFrom       = "rt_extra_from"
	ExtraSelectMode = "rt_extra_select_mode"
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

func (s *Server) SetTag(ctx context.Context, req *pb.SetTagReq, _ *emptypb.Empty) error {
	rt, ok := s.rts[req.SrvId]
	if !ok {
		return fmt.Errorf("service not found: srvId=%s", req.SrvId)
	}
	rt.addTag(req.GetTag(), req.GetValue())
	return nil
}

func (s *Server) RemoveTag(ctx context.Context, req *pb.ClearTagReq, _ *emptypb.Empty) error {
	rt, ok := s.rts[req.SrvId]
	if !ok {
		return fmt.Errorf("service not found: srvId=%s", req.SrvId)
	}
	rt.removeTag(req.GetTag())
	return nil
}

func (s *Server) GetAllTag(ctx context.Context, req *pb.AllTagReq, reply *pb.AllTagResp) error {
	rt, ok := s.rts[req.SrvId]
	if !ok {
		return nil
	}
	reply.Tags = map[string]string{}
	for k, v := range rt.selector.tags {
		reply.Tags[k] = v
	}
	return nil
}

func (s *Server) Route(ctx context.Context, param *pb.RouteReq, reply *pb.RouteReply) error {
	rt, ok := s.rts[param.GetSrvId()]
	if !ok {
		return fmt.Errorf("service not register: srvId=%s", param.GetSrvId())
	}
	reply.Success = true
	reply.Msg = "success"
	reply.Reply = &anypb.Any{}

	err := rt.route(ctx, param.Fn, param, reply)
	if err != nil {
		reply.Success = false
		reply.Msg = err.Error()
		return errors.Wrap(err, "service route error")
	}
	return nil
}

func (s *Server) Unregister(ctx context.Context, param *pb.UnRegisterReq, _ *emptypb.Empty) error {
	rv, ok := s.rts[param.SrvId]
	if ok {
		return rv.Close()
	}
	return errors.New("service not register")
}

func (s *Server) Register(ctx context.Context, param *pb.RegisterRtReq, _ *emptypb.Empty) error {
	sv, err := newService(&rpc.ClientOptions{
		Name:        param.GetSrvId(),
		EtcdServers: param.GetDiscoverySrvUrl(),
	})
	if err != nil {
		return err
	}
	err = sv.BaseClient.Run()
	if err != nil {
		return err
	}
	old, ok := s.rts[param.GetSrvId()]
	if ok {
		err := old.Close()
		if err != nil {
			logger.E("route register error", err)
		}
	}
	s.rts[param.GetSrvId()] = sv
	logger.D("service registered: %s", param.GetSrvId())
	return nil
}

func (s *Server) Run() error {
	// TODO sync tags from redis, init service
	return s.BaseServer.Run()
}
