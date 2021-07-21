package im

import (
	"errors"
	"go_im/im/entity"
)

var (
	ErrUnknownAction = errors.New("ErrUnknownAction")

	Api = newApi()
)

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
	case entity.ActionUserRegister:
		return a.Register(msg, en.(*entity.RegisterRequest))
	case entity.ActionUserSyncMsg:
		return a.SyncMessageList(msg)
	case entity.ActionUserRelation:
		return a.GetRelationList(msg)
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
