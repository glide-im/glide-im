package api

import (
	"errors"
	"go_im/im/message"
	"go_im/pkg/logger"
)

var MessageHandleFunc func(uid int64, message *message.Message)

var apiHandler IApiHandler

type IApiHandler interface {
	Handle(uid int64, message *message.Message)
}

func SetHandler(api IApiHandler) {
	apiHandler = api
}

func Handle(uid int64, message *message.Message) {
	apiHandler.Handle(uid, message)
}

func respond(uid int64, seq int64, action message.Action, data interface{}) {
	resp := message.NewMessage(seq, action, data)
	respondMessage(uid, resp)
}

func respondMessage(uid int64, msg *message.Message) {
	MessageHandleFunc(uid, msg)
}

type Routers struct {
	*UserApi
	*GroupApi
	*AppApi
	*TestApi
	router *Router
}

func NewApiRouter() *Routers {
	ret := new(Routers)
	ret.init()
	return ret
}

func (a *Routers) init() {
	rt := NewRouter()
	rt.Add(
		Group("api",
			Group("app",
				Route("echo", a.Echo),
			),
			Group("user",
				Route("login", a.Login),
				Route("logout", a.Logout),
				Route("auth", a.Auth),
				Route("register", a.Register),
				Route("online", a.GetOnlineUser),
				Group("info",
					Route("get", a.GetUserInfo),
					Route("me", a.UserInfo),
				),
			),
			Group("contacts",
				Route("get", a.GetAndInitRelationList),
				Route("put", a.AddFriend),
			),
			Group("chat",
				Route("list", a.GetUserChatList),
				Route("new", a.NewChat),
				Route("info", a.GetChatInfo),
				Route("history", a.GetChatHistory),
			),
			Group("group",
				Route("create", a.CreateGroup),
				Route("info", a.GetGroupInfo),
				Route("join", a.JoinGroup),
				Route("exit", a.ExitGroup),
				Group("member",
					Route("get", a.GetGroupMember),
					Route("add", a.AddGroupMember),
					Route("remove", a.RemoveMember),
				),
			),
			Group("test",
				Route("login", a.TestLogin),
				Route("signout", a.TestSignOut),
			),
		),
	)
	a.router = rt
}

func (a *Routers) Handle(uid int64, message *message.Message) {

	// TODO async
	err := a.handle(uid, message)
	if err != nil {
		a.onError(uid, message, err)
	}
}

func (a *Routers) handle(uid int64, message *message.Message) error {

	if err := a.intercept(uid, message); err != nil {
		return err
	}

	return a.router.Handle(uid, message)
}

const (
	actionLogin    message.Action = "api.user.login"
	actionRegister                = "api.user.register"
	actionAuth                    = "api.user.auth"
	actionEcho                    = "api.app.echo"
)

func (a *Routers) intercept(uid int64, message *message.Message) error {

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

func (a *Routers) onError(uid int64, msg *message.Message, err error) {
	logger.D("a.onError: uid=%d, Action=%s, err=%s", uid, msg.Action, err.Error())

	errMsg := message.NewMessage(msg.Seq, message.ActionFailed, err.Error())
	MessageHandleFunc(uid, errMsg)
}
