package broker

import (
	"context"
	"github.com/smallnest/rpcx/client"
	hash2 "go_im/pkg/hash"
	"go_im/pkg/logger"
	"go_im/pkg/rpc"
	"strconv"
)

type brokerSelector struct {
	srv map[string]string

	hash        *hash2.ConsistentHash
	roundRobbin client.Selector
}

func newBrokerSelector() *brokerSelector {
	ret := &brokerSelector{
		srv:         map[string]string{},
		hash:        hash2.NewConsistentHash(),
		roundRobbin: rpc.NewRoundRobinSelector(),
	}
	return ret
}

func (s *brokerSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	value, ok := ctx.Value(ctxKeyGid).(int64)
	if ok {
		gid := strconv.FormatInt(value, 10)
		node, err := s.hash.Get(gid)
		if err != nil {
			logger.E("consistent hash selector get error: %v", err)
			return ""
		}
		return node.Val
	}

	return s.roundRobbin.Select(ctx, servicePath, serviceMethod, args)
}

func (s *brokerSelector) UpdateServer(servers map[string]string) {

	// update node added
	for k, v := range servers {
		_, ok := s.srv[k]
		if !ok {
			s.srv[k] = v
			err := s.hash.Add(k)
			if err != nil {
				logger.E("consistent hash selector add node error:%v", err)
			}
		}
	}

	// update node removed
	for k := range s.srv {
		_, ok := servers[k]
		if !ok {
			delete(s.srv, k)
			err := s.hash.Remove(k)
			if err != nil {
				logger.E("consistent hash selector remove node error:%v", err)
			}
		}
	}
	s.roundRobbin.UpdateServer(servers)
}
