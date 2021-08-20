package im

import (
	"errors"
	"go_im/im/comm"
	"go_im/im/entity"
)

var (
	ErrUnknownAction = errors.New("ErrUnknownAction")
	ApiService       = newRootApiHandler()
)

var ActionDoNotNeedToken = map[entity.Action]int8{
	entity.ActionUserAuth:     0,
	entity.ActionUserLogin:    0,
	entity.ActionUserRegister: 0,
}

type RequestInfo struct {
	uid int64
	seq int64
}

type RootApiHandler struct {
	*userApi
	*groupApi
	rt map[entity.Action]interface{}
}

func newRootApiHandler() *RootApiHandler {
	ret := new(RootApiHandler)
	ret.rt = make(map[entity.Action]interface{})
	return ret
}

func (a *RootApiHandler) AddHandler(action entity.Action, handlerFunc interface{}) {
	a.rt[action] = handlerFunc
}

func (a *RootApiHandler) Handle(uid int64, message *entity.Message) {

	// TODO async
	err := a.handle(uid, message)
	if err != nil {
		a.onError(uid, message, err)
	}
}

func (a *RootApiHandler) handle(uid int64, message *entity.Message) error {

	if err := a.intercept(uid, message); err != nil {
		return err
	}

	en := entity.NewRequestFromAction(message.Action)

	if en != nil {
		e := message.DeserializeData(en)
		if e != nil {
			return e
		}
	}

	msg := &RequestInfo{
		uid: uid,
		seq: message.Seq,
	}

	switch message.Action {
	case entity.ActionUserLogin:
		return a.Login(msg, en.(*entity.LoginRequest))
	case entity.ActionUserAuth:
		return a.Auth(msg, en.(*entity.AuthRequest))
	case entity.ActionUserRegister:
		return a.Register(msg, en.(*entity.RegisterRequest))
	case entity.ActionUserChatList:
		return a.GetUserChatList(msg)
	case entity.ActionUserContacts:
		return a.GetAndInitRelationList(msg)
	case entity.ActionOnlineUser:
		return a.GetOnlineUser(msg)
	case entity.ActionUserNewChat:
		return a.NewChat(msg, en.(*entity.UserNewChatRequest))
	case entity.ActionUserChatHistory:
		return a.GetChatHistory(msg, en.(*entity.ChatHistoryRequest))
	case entity.ActionUserChatInfo:
		return a.GetChatInfo(msg, en.(*entity.ChatInfoRequest))
	case entity.ActionUserLogout:
	case entity.ActionUserEditInfo:
	case entity.ActionUserGetInfo:
		return a.GetUserInfo(msg, en.(*entity.UserInfoRequest))
	case entity.ActionUserAddFriend:
		return a.AddFriend(msg, en.(*entity.AddContacts))
	case entity.ActionUserInfo:
		return a.UserInfo(msg)
	case entity.ActionGroupCreate:
		return a.CreateGroup(msg, en.(*entity.CreateGroupRequest))
	case entity.ActionGroupInfo:
		return a.GetGroupInfo(msg, en.(*entity.GroupInfoRequest))
	case entity.ActionGroupExit:
		return a.ExitGroup(msg, en.(*entity.ExitGroupRequest))
	case entity.ActionGroupJoin:
		return a.JoinGroup(msg, en.(*entity.JoinGroupRequest))
	case entity.ActionGroupAddMember:
		return a.AddGroupMember(msg, en.(*entity.AddMemberRequest))
	case entity.ActionGroupGetMember:
		return a.GetGroupMember(msg, en.(*entity.GetGroupMemberRequest))
	default:
		return ErrUnknownAction
	}

	return ErrUnknownAction
}

func (a *RootApiHandler) intercept(uid int64, message *entity.Message) error {

	_, ok := ActionDoNotNeedToken[message.Action]
	if uid <= 0 && !ok {
		return errors.New("unauthorized")
	}

	// verify fields

	// something else
	return nil
}

func (a *RootApiHandler) onError(uid int64, message *entity.Message, err error) {
	comm.Slog.D("ApiService.onError: uid=%d, Action=%s, err=%s", uid, message.Action, err.Error())

	msg := entity.NewMessage(message.Seq, entity.ActionNotify, err.Error())
	ClientManager.EnqueueMessage(uid, msg)
}
