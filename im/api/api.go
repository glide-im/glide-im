package api

import (
	"errors"
	"go_im/im/api/apidep"
	"go_im/im/api/app"
	"go_im/im/api/auth"
	"go_im/im/api/groups"
	"go_im/im/api/http_srv"
	"go_im/im/api/router"
	"go_im/im/api/test"
	"go_im/im/api/user"
	"go_im/im/message"
	"go_im/pkg/logger"
)

var Handler ApiHandler = NewDefaultRouter()

type ApiHandler interface {
	Handle(uid int64, device int64, message *message.Message)
}

// Handle 处理一个 api 消息
func Handle(uid int64, device int64, message *message.Message) {
	Handler.Handle(uid, device, message)
}

// RunHttpServer 启动 http 服务器, 以 HTTP 服务方式访问 api
func RunHttpServer(addr string, port int) error {
	return http_srv.Run(addr, port)
}

// Routers 默认 api 路由, 将不同 action 交给相应的方法处理
type Routers struct {
	*user.UserApi
	*groups.GroupApi
	*auth.AuthApi
	*app.AppApi
	*test.TestApi
	router *route.Router
}

func NewDefaultRouter() *Routers {
	ret := new(Routers)
	ret.init()
	return ret
}

func (a *Routers) init() {
	rt := route.NewRouter()
	rt.Add(
		route.Group("api",
			route.Group("app",
				route.Route("echo", a.Echo),
			),
			route.Group("user",
				route.Route("login", a.Login),
				route.Route("logout", a.Logout),
				route.Route("auth", a.Auth),
				route.Route("register", a.Register),
				route.Route("online", a.GetOnlineUser),
				route.Group("info",
					route.Route("get", a.GetUserInfo),
					route.Route("me", a.UserInfo),
				),
			),
			route.Group("contacts",
				route.Route("get", a.GetAndInitRelationList),
				route.Route("put", a.AddFriend),
			),
			route.Group("chat",
				route.Route("list", a.GetUserChatList),
				route.Route("new", a.NewChat),
			),
			route.Group("group",
				route.Route("create", a.CreateGroup),
				route.Route("info", a.GetGroupInfo),
				route.Route("join", a.JoinGroup),
				route.Route("exit", a.ExitGroup),
				route.Group("member",
					route.Route("get", a.GetGroupMember),
					route.Route("add", a.AddGroupMember),
					route.Route("remove", a.RemoveMember),
				),
			),
			route.Group("test",
				route.Route("login", a.TestLogin),
				route.Route("signout", a.TestSignOut),
			),
		),
	)
	a.router = rt
}

func (a *Routers) Handle(uid int64, device int64, message *message.Message) {

	err := a.handle(uid, device, message)
	if err != nil {
		a.onError(uid, device, message, err)
	}
}

func (a *Routers) handle(uid int64, device int64, message *message.Message) error {

	if err := a.intercept(uid, device, message); err != nil {
		return err
	}

	return a.router.Handle(uid, device, message)
}

const (
	actionLogin    message.Action = "api.user.login"
	actionRegister                = "api.user.register"
	actionAuth                    = "api.user.auth"
	actionEcho                    = "api.app.echo"
)

func (a *Routers) intercept(uid int64, device int64, message *message.Message) error {

	if message.Action.Contains("api.test") {
		return nil
	}

	doNotNeedAuth := message.Action == actionLogin || message.Action == actionRegister || message.Action == actionAuth || message.Action == actionEcho
	if uid <= 0 && !doNotNeedAuth {
		return errors.New("unauthorized")
	}

	// verify fields

	// something else
	return nil
}

func (a *Routers) onError(uid int64, device int64, msg *message.Message, err error) {
	logger.D("a.onError: uid=%d, Action=%s, err=%s", uid, msg.Action, err.Error())

	errMsg := message.NewMessage(msg.Seq, message.ActionFailed, err.Error())
	apidep.ClientManager.EnqueueMessage(uid, device, errMsg)
}
