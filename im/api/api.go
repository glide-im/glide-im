package api

import (
	"errors"
	"go_im/im/api/apidep"
	"go_im/im/api/app"
	"go_im/im/api/auth"
	"go_im/im/api/groups"
	"go_im/im/api/msg"
	"go_im/im/api/router"
	"go_im/im/api/test"
	"go_im/im/api/user"
	"go_im/im/message"
	"go_im/pkg/logger"
	"strings"
)

// Routers 默认 api 路由, 将不同 action 交给相应的方法处理
type Routers struct {
	*user.UserApi
	*groups.GroupApi
	*auth.AuthApi
	*app.AppApi
	*msg.MsgApi
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
				route.Route("login", a.SignIn),
				route.Route("logout", a.Logout),
				route.Route("auth", a.AuthToken),
				route.Route("register", a.Register),
				route.Route("guest", a.GuestRegister),
				route.Route("online", a.GetOnlineUser),
				route.Group("info",
					route.Route("get", a.GetUserInfo),
					route.Route("me", a.UserProfile),
				),
			),
			route.Group("contacts",
				route.Route("get", a.GetContactList),
				route.Route("add", a.AddContact),
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

func (a *Routers) Handle(uid int64, device int64, message *message.Message) (*message.Message, error) {

	logger.D("%v", message)
	m, err := a.handle(uid, device, message)
	return m, err
}

func (a *Routers) handle(uid int64, device int64, message *message.Message) (*message.Message, error) {

	if err := a.intercept(uid, device, message); err != nil {
		return nil, err
	}

	return a.router.Handle(uid, device, message)
}

const (
	actionLogin    string = "api.user.login"
	actionRegister        = "api.user.register"
	actionAuth            = "api.user.auth"
	actionEcho            = "api.app.echo"
)

func (a *Routers) intercept(uid int64, device int64, message *message.Message) error {

	if strings.HasPrefix(message.GetAction(), "api.test") {
		return nil
	}

	action := message.GetAction()
	doNotNeedAuth := action == actionLogin || action == actionRegister || action == actionAuth || action == actionEcho
	if uid <= 0 && !doNotNeedAuth {
		return errors.New("unauthorized")
	}

	// verify fields

	// something else
	return nil
}

func (a *Routers) onError(uid int64, device int64, msg *message.Message, err error) {
	logger.D("a.onError: uid=%d, Action=%s, err=%s", uid, msg.GetAction(), err.Error())

	errMsg := message.NewMessage(msg.GetSeq(), message.ActionApiFailed, err.Error())
	apidep.SendMessageIfOnline(uid, device, errMsg)
}
