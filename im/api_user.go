package im

import (
	"go_im/im/dao"
	"go_im/im/entity"
)

type userApi struct{}

func (a *userApi) Auth(msg *ApiMessage, request *entity.AuthRequest) (*entity.Message, bool, error) {

	var resp = entity.NewMessage(msg.seq, entity.RespActionSuccess)
	uid := dao.UserDao.GetUid(request.Token)
	if uid == request.Uid {
		return resp, true, nil
	} else {
		return resp, false, nil
	}
}

func (a *userApi) Login(msg *ApiMessage, request *entity.LoginRequest) (*entity.Message, int64, error) {

	uid, token, err := dao.UserDao.GetUidByLogin(request.Account, request.Password)
	if err != nil {
		return nil, uid, err
	}

	if len(request.Password) != 0 && len(request.Account) != 0 {
		m := entity.NewMessage(msg.seq, entity.RespActionSuccess)
		if err := m.SetData(entity.AuthorResponse{Token: token}); err != nil {
			return nil, uid, err
		}
		return m, uid, nil
	} else {
		return entity.NewSimpleMessage(msg.seq, entity.RespActionUserUnauthorized, "unauthorized"), uid, nil
	}
}

func (a *userApi) GetRelationList(msg *ApiMessage) error {

	groups := dao.GroupDao.GetUserGroup(msg.uid)
	for _, gid := range groups {
		group := GroupManager.GetGroup(gid)
		mc := ClientManager.GetClient(msg.uid).messages
		group.Subscribe(msg.uid, mc)
	}

	friends := dao.UserDao.GetFriends(msg.uid)
	relation := entity.RelationResponse{
		Groups:  groups,
		Friends: friends,
	}

	resp := entity.NewMessage(msg.seq, entity.RespActionSuccess)
	if err := resp.SetData(relation); err != nil {
		return err
	}

	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (a *userApi) SyncMessageList(msg *ApiMessage) error {

	chats := dao.UserDao.GetMessageList(msg.uid)
	for _, chat := range chats {
		dao.MessageDao.GetChatInfo(chat)
	}

	return nil
}

func (a *userApi) GetUserInfo(msg *ApiMessage, request *entity.UserInfoRequest) error {

	return nil
}

func (a *userApi) UserInfo(msg *ApiMessage) error {

	return nil
}

func (a *userApi) Register(msg *ApiMessage, registerEntity *entity.RegisterRequest) (*entity.Message, error) {

	resp := entity.NewMessage(msg.seq, entity.RespActionSuccess)
	err := dao.UserDao.AddUser(registerEntity.Account, registerEntity.Password)

	if err != nil {
		resp.Data = "seq err=" + err.Error()
	} else {
		resp.Data = "register success"
	}

	return resp, err
}
