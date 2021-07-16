package im

import (
	"errors"
	"go_im/im/entity"
)

var (
	ErrUnknownAction = errors.New("ErrUnknownAction")

	Api = newApi()
)

type api struct {
	*userApi
	*groupApi

	actionEntityMap map[entity.Action]interface{}
}

func newApi() *api {
	ret := new(api)
	ret.userApi = new(userApi)
	ret.groupApi = new(groupApi)
	ret.actionEntityMap = map[entity.Action]interface{}{
		entity.ActionUserLogin:    &entity.LoginRequest{},
		entity.ActionUserRegister: &entity.RegisterRequest{},
	}
	return ret
}

func (a *api) Handle(client *Client, message *entity.Message) error {

	en := a.actionEntityMap[message.Action]
	e := message.DeserializeData(en)
	if e != nil {
		return e
	}

	switch message.Action {
	case entity.ActionUserLogin:
		return a.login(client, message.Seq, en.(*entity.LoginRequest))
	case entity.ActionUserRegister:
		return a.register(client, message.Seq, en.(*entity.RegisterRequest))
	default:
		return ErrUnknownAction
	}
}
