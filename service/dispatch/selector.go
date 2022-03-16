package dispatch

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"go_im/pkg/logger"
	"go_im/pkg/rpc"
	"reflect"
	"strconv"
)

const ctxKeyCalculateVal = "key_hash_calculate_value"

// dispatchSelector
// 一致性 Hash 节点选择器, 客户端 rpc 调用时在上下文中附带 hash 计算值, selector 通过该值继续路由.
// Consistent has server selector.
type dispatchSelector struct {
	srv map[string]string

	hash        *ConsistentHash
	roundRobbin client.Selector
}

func newSelector() *dispatchSelector {
	ret := &dispatchSelector{
		srv:         map[string]string{},
		hash:        NewConsistentHash(),
		roundRobbin: rpc.NewRoundRobinSelector(),
	}
	return ret
}

func (s *dispatchSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	value, ok := ctx.Value(ctxKeyCalculateVal).(int64)
	if ok {
		uid := strconv.FormatInt(value, 10)
		node, err := s.hash.Get(uid)
		if err != nil {
			logger.E("consistent hash selector get error: %v", err)
			return ""
		}
		return node.val
	}

	return s.roundRobbin.Select(ctx, servicePath, serviceMethod, args)
}

func (s *dispatchSelector) UpdateServer(servers map[string]string) {

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

// contextOfUidHashRoute
// 返回一个通过 uid 路由的 Context
// Create the context with uid route info.
func contextOfUidHashRoute(uid int64) context.Context {
	return context.WithValue(context.Background(), ctxKeyCalculateVal, uid)
}

func reflectMethodName(method interface{}) string {
	typeOf := reflect.TypeOf(method)
	if typeOf.Kind() != reflect.Func {
		// not func
	} else {
		return typeOf.Name()
	}
	return ""
}
