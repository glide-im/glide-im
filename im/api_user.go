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
	if err = m.SetData(entity.AuthorResponse{Token: token, Uid: uid}); err != nil {
		return nil, uid, err
	}
	return m, uid, nil
}

func (a *userApi) GetAndInitRelationList(msg *ApiMessage) error {

	groups, err := dao.GroupDao.GetUserGroup(msg.uid)
	if err != nil {
		return err
	}
	friends, err := dao.UserDao.GetFriends(msg.uid)
	if err != nil {
		return err
	}

	contacts := make([]entity.ContactResponse, len(groups))

	for _, g := range groups {
		contacts = append(contacts, entity.ContactResponse{
			Id:     g.Gid,
			Avatar: g.Avatar,
			Name:   g.Name,
			Type:   2,
		})
		group := GroupManager.GetGroup(g.Gid)
		client := ClientManager.GetClient(msg.uid)
		client.AddGroup(group)
	}

	for _, friend := range friends {
		contacts = append(contacts, entity.ContactResponse{
			Id:     friend.Uid,
			Avatar: "",
			Name:   friend.Remark,
			Type:   1,
		})
	}

	resp := entity.NewMessage(msg.seq, entity.RespActionSuccess)
	if err = resp.SetData(contacts); err != nil {
		return err
	}

	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (a *userApi) AddFriend(msg *ApiMessage, request *entity.AddFriendRequest) error {

	users, err := dao.UserDao.GetUser(request.Uid)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		resp := entity.NewSimpleMessage(msg.seq, entity.RespActionFailed, "user not exist.")
		ClientManager.EnqueueMessage(msg.uid, resp)
		return nil
	}

	friend, err := dao.UserDao.AddFriend(msg.uid, request.Uid, request.Remark)
	if err != nil {
		return err
	}
	resp := entity.NewMessage(msg.seq, entity.RespActionSuccess)
	if err = resp.SetData(friend); err != nil {
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, resp)

	// friend
	f, err := dao.UserDao.AddFriend(request.Uid, msg.uid, "")
	if err != nil {
		return err
	}
	resp1 := entity.NewMessage(msg.seq, entity.RespActionFriendApproval)
	if err = resp.SetData(f); err != nil {
		return err
	}
	ClientManager.EnqueueMessage(request.Uid, resp1)

	return nil
}

func (a *userApi) GetUserInfo(msg *ApiMessage, request *entity.UserInfoRequest) error {

	users, err := dao.UserDao.GetUser(request.Uid...)
	if err != nil {
		return err
	}
	resp := entity.NewMessage(msg.seq, entity.ActionOnlineUser)
	type u struct {
		Uid      int64
		Account  string
		Avatar   string
		Nickname string
	}
	var ret []u
	for _, user := range users {
		retU := u{
			Uid:      user.Uid,
			Account:  user.Account,
			Avatar:   user.Account,
			Nickname: user.Nickname,
		}
		ret = append(ret, retU)
	}
	if err = resp.SetData(ret); err != nil {
		return err
	}
	return nil
}

func (a *userApi) GetChatInfo(msg *ApiMessage, request *entity.ChatInfoRequest) error {

	uc, err := dao.ChatDao.GetUserChatFromChat(request.Cid, msg.uid)
	if err != nil {
		return err
	}
	resp := entity.NewMessage(msg.seq, entity.ActionUserChatInfo)
	if err = resp.SetData(uc); err != nil {
		return err
	}
	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (a *userApi) GetChatHistory(msg *ApiMessage, request *entity.ChatHistoryRequest) error {

	chatMessages, err := dao.ChatDao.GetChatHistory(request.Cid, 20)
	if err != nil {
		return err
	}

	resp := entity.NewMessage(msg.seq, entity.ActionUserChatHistory)
	if err = resp.SetData(chatMessages); err != nil {
		return err
	}

	ClientManager.EnqueueMessage(msg.uid, resp)
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

	ClientManager.Update()
	for k := range ClientManager.AllClient() {
		us, err := dao.UserDao.GetUser(k)
		if err != nil || len(us) == 0 {
			logger.D("get online uid=%d error, error=%v", k, err)
			continue
		}
		user := us[0]
		users = append(users, u{Uid: user.Uid, Account: user.Account, Avatar: user.Avatar, Nickname: user.Nickname})
	}

	_ = m.SetData(users)
	ClientManager.EnqueueMessage(msg.uid, m)
	return nil
}

func (a *userApi) NewChat(msg *ApiMessage, request *entity.UserNewChatRequest) error {

	uid := msg.uid
	target := request.Id

	c, err := dao.ChatDao.NewChat(uid, target, request.Type)
	if err != nil {
		return err
	}

	// chat
	if request.Type == 1 {
		m2, err2 := dao.ChatDao.NewUserChat(c.Cid, uid, target, 1)
		if err2 != nil {
			return err2
		}
		_, err = dao.ChatDao.NewUserChat(c.Cid, int64(target), uid, 1)
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
	list, err := dao.ChatDao.GetUserChatList(msg.uid)
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
