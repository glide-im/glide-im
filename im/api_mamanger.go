package im

import (
	"errors"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/comm"
	"go_im/im/message"
)

type ApiRouter struct {
	*api.UserApi
	*api.GroupApi
	router *api.Router
}

func newApiRouter() *ApiRouter {
	ret := new(ApiRouter)
	ret.init()
	return ret
}

func (a *ApiRouter) init() {
	rt := api.NewRouter()
	rt.Add(
		api.Group("api",
			api.Group("user",
				api.Route("login", a.Login),
				api.Route("auth", a.Auth),
				api.Route("register", a.Register),
				api.Route("online", a.GetOnlineUser),
				api.Group("info",
					api.Route("get", a.GetUserInfo),
					api.Route("me", a.UserInfo),
				),
			),
			api.Group("contacts",
				api.Route("get", a.GetAndInitRelationList),
				api.Route("add", a.AddFriend),
			),
			api.Group("chat",
				api.Route("list", a.GetUserChatList),
				api.Route("new", a.NewChat),
				api.Route("info", a.GetChatInfo),
				api.Route("history", a.GetChatHistory),
			),
			api.Group("group",
				api.Route("create", a.CreateGroup),
				api.Route("info", a.GetGroupInfo),
				api.Route("join", a.JoinGroup),
				api.Route("exit", a.ExitGroup),
				api.Group("member",
					api.Route("get", a.GetGroupMember),
					api.Route("add", a.AddGroupMember),
					api.Route("remove", a.RemoveMember),
				),
			),
		),
	)
	a.router = rt
}

func (a *ApiRouter) Handle(uid int64, message *message.Message) {

	// TODO async
	err := a.handle(uid, message)
	if err != nil {
		a.onError(uid, message, err)
	}
}

func (a *ApiRouter) handle(uid int64, message *message.Message) error {

	if err := a.intercept(uid, message); err != nil {
		return err
	}

	return a.router.Handle(uid, message)
}

const (
	actionLogin    message.Action = "api.user.login"
	actionRegister                = "api.user.register"
	actionAuth                    = "api.user.auth"
)

func (a *ApiRouter) intercept(uid int64, message *message.Message) error {

	doNotNeedAuth := message.Action == actionLogin || message.Action == actionRegister || message.Action == actionAuth
	if uid <= 0 && !doNotNeedAuth {
		return errors.New("unauthorized")
	}

	// verify fields

	// something else
	return nil
}

func (a *ApiRouter) onError(uid int64, msg *message.Message, err error) {
	comm.Slog.D("a.onError: uid=%d, Action=%s, err=%s", uid, msg.Action, err.Error())

	errMsg := message.NewMessage(msg.Seq, message.ActionNotify, err.Error())
	client.EnqueueMessage(uid, errMsg)
}
