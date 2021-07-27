package im

import (
	"go_im/im/dao"
	"go_im/im/entity"
)

type ApiFatalError struct {
	msg string
}

func (f *ApiFatalError) Error() string {
	return f.msg
}

func newApiFatalError(msg string) *ApiFatalError {
	return &ApiFatalError{msg: msg}
}

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

	if len(request.Account) == 0 || len(request.Password) == 0 {
		return entity.NewSimpleMessage(msg.seq, entity.RespActionUserUnauthorized, "account or password empty"), -1, nil
	}

	uid, token, err := dao.UserDao.GetUidByLogin(request.Account, request.Password)
	if err != nil {
		return nil, uid, err
	}

	m := entity.NewMessage(msg.seq, entity.RespActionSuccess)
	if err = m.SetData(entity.AuthorResponse{Token: token}); err != nil {
		return nil, uid, err
	}
	return m, uid, nil
}

func (a *userApi) GetAndInitRelationList(msg *ApiMessage) error {

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

func (a *userApi) GetUserInfo(msg *ApiMessage, request *entity.UserInfoRequest) error {

	return nil
}

func (a *userApi) GetOnlineUser(msg *ApiMessage) error {

	m := entity.NewMessage(msg.seq, entity.ActionOnlineUser)
	type u struct {
		Uid      int64
		Account  string
		Avatar   string
		Nickname string
	}
	var users []u

	for k := range ClientManager.AllClient() {
		user, err := dao.UserDao.GetUser(k)
		if err != nil {
			logger.D("get online uid=%d error, error=%v", k, err)
			continue
		}
		users = append(users, u{Uid: user.Uid, Account: user.Account, Avatar: user.Avatar, Nickname: user.Nickname})
	}

	_ = m.SetData(users)
	ClientManager.EnqueueMessage(msg.uid, m)
	return nil
}

func (a *userApi) NewChat(msg *ApiMessage, request *entity.UserNewChatRequest) error {

	uid := msg.uid
	target := request.Id

	c, err := dao.MessageDao.NewChat(uid, target, request.Type)
	if err != nil {
		return err
	}

	// chat
	if request.Type == 1 {
		m2, err2 := dao.MessageDao.NewUserChat(c.Cid, uid, target, 1)
		if err2 != nil {
			return err2
		}
		_, err = dao.MessageDao.NewUserChat(c.Cid, int64(target), uint64(uid), 1)
		if err != nil {
			return err
		}
		resp := entity.NewMessage(msg.seq, entity.RespActionSuccess)
		if err = resp.SetData(m2); err != nil {
			return err
		}
		ClientManager.EnqueueMessage(msg.uid, resp)
	} else {
		ClientManager.EnqueueMessage(msg.uid, entity.NewAckMessage(msg.seq))
	}
	return nil
}

func (a *userApi) GetUserChatList(msg *ApiMessage) error {

	resp := entity.NewMessage(msg.seq, entity.RespActionSuccess)
	list, err := dao.MessageDao.GetUserChatList(msg.uid)
	if err != nil {
		return err
	}
	if err = resp.SetData(list); err != nil {
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, resp)
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
