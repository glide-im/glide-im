package im

import (
	"errors"
	"go_im/im/entity"
)

var (
	ErrUnknownAction = errors.New("ErrUnknownAction")
	Api              = newApi()
)

var ActionDoNotNeedToken = map[entity.Action]int8{
	entity.ActionUserAuth:     0,
	entity.ActionUserLogin:    0,
	entity.ActionUserRegister: 0,
}

type ApiMessage struct {
	uid int64
	seq int64
}

type api struct {
	*userApi
	*groupApi
}

func newApi() *api {
	ret := new(api)
	ret.userApi = new(userApi)
	ret.groupApi = new(groupApi)
	return ret
}

func (a *api) Handle(client *Client, message *entity.Message) error {

	if err := a.intercept(client, message); err != nil {
		return err
	}

	en := entity.NewRequestFromAction(message.Action)

	if en != nil {
		e := message.DeserializeData(en)
		if e != nil {
			return e
		}
	}

	msg := &ApiMessage{
		uid: client.uid,
		seq: message.Seq,
	}

	switch message.Action {
	case entity.ActionUserLogin:
		req := en.(*entity.LoginRequest)
		m, uid, err := a.Login(msg, req)
		if err != nil {
			return err
		}
		client.SignIn(uid, req.Device)
		client.EnqueueMessage(m)
		return nil
	case entity.ActionUserAuth:
		req := en.(*entity.AuthRequest)
		m, success, err := a.Auth(msg, req)
		if err != nil {
			return err
		}
		if success {
			client.SignIn(req.Uid, req.DeviceId)
		}
		client.EnqueueMessage(m)
		return nil
	case entity.ActionUserRegister:
		req := en.(*entity.RegisterRequest)
		m, err := a.Register(msg, req)
		if err != nil {
			return err
		}
		client.EnqueueMessage(m)
		return nil

	case entity.ActionUserChatList:
		return a.GetUserChatList(msg)
	case entity.ActionUserRelation:
		return a.GetAndInitRelationList(msg)
	case entity.ActionOnlineUser:
		return a.GetOnlineUser(msg)
	case entity.ActionUserNewChat:
		return a.NewChat(msg, en.(*entity.UserNewChatRequest))
	case entity.ActionUserLogout:
	case entity.ActionUserEditInfo:
	case entity.ActionUserGetInfo:
		return a.GetUserInfo(msg, en.(*entity.UserInfoRequest))
	case entity.ActionUserInfo:
		return a.UserInfo(msg)
	default:
		return ErrUnknownAction
	}

	return ErrUnknownAction
}

func (a *api) intercept(client *Client, message *entity.Message) error {

	_, ok := ActionDoNotNeedToken[message.Action]
	if client.uid <= 0 && !ok {
		return errors.New("unauthorized")
	}

	// verify fields

	// something else
	return nil
}
