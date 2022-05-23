package api

import (
	"go_im/im/api/apidep"
	"go_im/im/message"
)

var handler Interface = NewDefaultRouter()

// Interface
// Deprecated
type Interface interface {
	Handle(uid int64, device int64, message *message.Message) (*message.Message, error)
}

// SetInterfaceImpl
//
// Deprecated
func SetInterfaceImpl(i Interface) {
	handler = i
}

func SetClientInterfaceImpl(i apidep.ClientManagerInterface) {
	apidep.ClientInterface = i
}

func SetGroupInterfaceImpl(i apidep.GroupManagerInterface) {
	apidep.GroupInterface = i
}

func MockDep() {
	apidep.GroupInterface = &apidep.MockGroupManager{}
	apidep.ClientInterface = &apidep.MockClientManager{}
}

// Handle 处理一个 api 消息
//
// Deprecated
func Handle(uid int64, device int64, message *message.Message) (*message.Message, error) {
	return handler.Handle(uid, device, message)
}

// RunHttpServer 启动 http 服务器, 以 HTTP 服务方式访问 api
func RunHttpServer(addr string, port int) error {
	return run(addr, port)
}
