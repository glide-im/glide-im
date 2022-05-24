package dispatch

import (
	"context"
	"errors"
	"go_im/pkg/rpc"
	"go_im/protobuf/gen/pb_rpc"
)

var cache *routeCache
var nsq *nsqMsgProducer

type Server struct {
	*rpc.BaseServer
}

func NewServer(msgNsqdAddr string, options *rpc.ServerOptions) (*rpc.BaseServer, error) {
	s := rpc.NewBaseServer(options)
	s.Register(options.Name, &Server{})
	cache = newRouteCache()
	var err error
	nsq, err = newNsqMsgProducer(msgNsqdAddr)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Dispatch(ctx context.Context, param *pb_rpc.DispatchRequest, replay *pb_rpc.Response) error {
	if param.Direct {
		return nsq.publish(param.RouteVal, param.Data)
	}
	node := cache.getRoute(param.SrvName, param.Id)
	if node == "" {
		return errors.New("route not fund")
	}
	err := nsq.publish(node, param.Data)
	return err
}

func (s *Server) UpdateRoute(ctx context.Context, param *pb_rpc.UpdateRouteRequest, replay *pb_rpc.Response) error {
	cache.updateRoute(param.GetSrvName(), param.GetId(), param.GetNode())
	return nil
}

func (s *Server) GetUserGateway(ctx context.Context, param *pb_rpc.UidRequest, replay *pb_rpc.UserGatewayResponse) error {
	node := cache.getRoute("gateway", param.GetUid())
	if node == "" {
		return errors.New("route not fund")
	}
	replay.Node = node
	return nil
}
