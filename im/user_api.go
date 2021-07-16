package im

import "go_im/im/entity"

type userApi struct{}

func (a *userApi) login(client *Client, seq int64, loginEntity *entity.LoginRequest) error {
	if len(loginEntity.Password) != 0 && len(loginEntity.Username) != 0 {
		m := entity.NewMessage(seq, entity.ActionSuccess)
		if err := m.SetData(entity.AuthorResponse{Token: "this is token"}); err != nil {
			return err
		}
		client.SignIn(1234, loginEntity.Device)
		client.EnqueueMessage(m)
	} else {
		client.EnqueueMessage(entity.NewSimpleMessage(seq, entity.ActionUserUnauthorized, "unauthorized"))
	}
	return nil
}

func (a *userApi) getInfo() {

}

func (a *userApi) register(client *Client, seq int64, registerEntity *entity.RegisterRequest) error {

	return nil
}
