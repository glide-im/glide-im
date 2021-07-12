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
	actionEntityMap map[entity.Action]interface{}
}

func newApi() *api {
	ret := new(api)
	ret.actionEntityMap = map[entity.Action]interface{}{
		entity.ActionUserLogin:    &entity.LoginEntity{},
		entity.ActionUserRegister: &entity.RegisterEntity{},
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
		return a.login(client, message.Seq, en.(*entity.LoginEntity))
	case entity.ActionUserRegister:
		return a.register(client, message.Seq, en.(*entity.RegisterEntity))
	default:
		return ErrUnknownAction
	}
}

func (a *api) login(client *Client, seq int64, loginEntity *entity.LoginEntity) error {
	if len(loginEntity.Password) != 0 && len(loginEntity.Username) != 0 {
		client.EnqueueMessage(&entity.Message{
			Seq:  seq,
			Data: []byte("login success"),
		})
	} else {
		client.EnqueueMessage(entity.NewMessage(seq, entity.ActionUserUnauthorized, "unauthorized"))
	}
	return nil
}

func (a *api) register(client *Client, seq int64, registerEntity *entity.RegisterEntity) error {

	return nil
}
