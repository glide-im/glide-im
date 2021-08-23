package im

import (
	"errors"
	"go_im/im/api"
	"go_im/im/comm"
	"go_im/im/entity"
)

var (
	ApiManager = newApiRouter()
)

var ActionDoNotNeedToken = map[entity.Action]int8{
	entity.ActionUserAuth:     0,
	entity.ActionUserLogin:    0,
	entity.ActionUserRegister: 0,
}

func init() {
	rt := api.NewRouter()
	rt.Add(
		api.Group("api",
			api.Group("user",
				api.Route("login", ApiManager.Login),
				api.Route("auth", ApiManager.Auth),
				api.Route("register", ApiManager.Register),
				api.Route("online", ApiManager.GetOnlineUser),
				api.Group("info",
					api.Route("get", ApiManager.GetUserInfo),
					api.Route("me", ApiManager.UserInfo),
				),
			),
			api.Group("contacts",
				api.Route("get", ApiManager.GetAndInitRelationList),
				api.Route("add", ApiManager.AddFriend),
			),
			api.Group("chat",
				api.Route("list", ApiManager.GetUserChatList),
				api.Route("new", ApiManager.NewChat),
				api.Route("info", ApiManager.GetChatInfo),
				api.Route("history", ApiManager.GetChatHistory),
			),
			api.Group("group",
				api.Route("create", ApiManager.CreateGroup),
				api.Route("info", ApiManager.GetGroupInfo),
				api.Route("join", ApiManager.JoinGroup),
				api.Route("exit", ApiManager.ExitGroup),
				api.Group("member",
					api.Route("get", ApiManager.GetGroupMember),
					api.Route("add", ApiManager.AddGroupMember),
					api.Route("remove", ApiManager.RemoveMember),
				),
			),
		),
	)
	SetRouter(rt)
}

type ApiRouter struct {
	*userApi
	*groupApi
	router *api.Router
}

func newApiRouter() *ApiRouter {
	ret := new(ApiRouter)
	return ret
}

func SetRouter(router *api.Router) {
	ApiManager.router = router
}

func (a *ApiRouter) Handle(uid int64, message *entity.Message) {

	// TODO async
	err := a.handle(uid, message)
	if err != nil {
		a.onError(uid, message, err)
	}
}

func (a *ApiRouter) handle(uid int64, message *entity.Message) error {

	if err := a.intercept(uid, message); err != nil {
		return err
	}

	return a.router.Handle(uid, message)
}

func (a *ApiRouter) intercept(uid int64, message *entity.Message) error {

	_, ok := ActionDoNotNeedToken[message.Action]
	if uid <= 0 && !ok {
		return errors.New("unauthorized")
	}

	// verify fields

	// something else
	return nil
}

func (a *ApiRouter) onError(uid int64, message *entity.Message, err error) {
	comm.Slog.D("ApiManager.onError: uid=%d, Action=%s, err=%s", uid, message.Action, err.Error())

	msg := entity.NewMessage(message.Seq, entity.ActionNotify, err.Error())
	ClientManager.EnqueueMessage(uid, msg)
}
