package im

import (
	"errors"
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

	var resp = entity.NewMessage(msg.seq, entity.ActionSuccess)
	uid := dao.UserDao.GetUid(request.Token)
	if uid == request.Uid {
		return resp, true, nil
	} else {
		return resp, false, nil
	}
}

func (a *userApi) Login(msg *ApiMessage, request *entity.LoginRequest) (*entity.Message, int64, error) {

	if len(request.Account) == 0 || len(request.Password) == 0 {
		return entity.NewSimpleMessage(msg.seq, entity.ActionUserUnauthorized, "account or password empty"), -1, nil
	}

	uid, token, err := dao.UserDao.GetUidByLogin(request.Account, request.Password)
	if err != nil {
		return nil, uid, err
	}

	m := entity.NewMessage(msg.seq, entity.ActionSuccess)
	if err = m.SetData(entity.AuthorResponse{Token: token, Uid: uid}); err != nil {
		return nil, uid, err
	}
	return m, uid, nil
}

func (a *userApi) GetAndInitRelationList(msg *ApiMessage) error {

	allContacts, err := dao.UserDao.GetAllContacts(msg.uid)
	if err != nil {
		return err
	}

	var friends []*entity.UserInfoResponse
	var groups []*entity.GroupResponse

	var uids []int64
	for _, contacts := range allContacts {

		if contacts.Type == dao.ContactsTypeGroup {
			group := GroupManager.GetGroup(contacts.TargetId)
			if group == nil {
				return newApiFatalError("load user group error: nil")
			}
			ClientManager.GetClient(msg.uid).AddGroup(group)
			groups = append(groups, &entity.GroupResponse{
				Group:   *group.group,
				Members: group.GetMembers(),
			})
		} else if contacts.Type == dao.ContactsTypeUser {
			uids = append(uids, contacts.TargetId)
		}
	}
	if len(uids) > 0 {
		user, err := dao.UserDao.GetUser(uids...)
		if err != nil {
			return err
		}
		for _, u := range user {
			friends = append(friends, &entity.UserInfoResponse{
				Uid:      u.Uid,
				Account:  u.Account,
				Nickname: u.Nickname,
				Avatar:   u.Avatar,
			})
		}
	}

	body := entity.ContactResponse{
		Friends: friends,
		Groups:  groups,
	}

	resp := entity.NewMessage2(msg.seq, entity.ActionSuccess, body)
	ClientManager.EnqueueMessage(msg.uid, resp)
	return nil
}

func (a *userApi) AddContacts(msg *ApiMessage, request *entity.AddContacts) error {

	users, err := dao.UserDao.GetUser(request.Uid)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		resp := entity.NewErrMessage2(msg.seq, "user not exist")
		ClientManager.EnqueueMessage(msg.uid, resp)
		return nil
	}

	hasContacts, err := dao.UserDao.HasContacts(msg.uid, request.Uid, dao.ContactsTypeUser)
	if err != nil {
		return err
	}

	if hasContacts {
		resp := entity.NewErrMessage2(msg.seq, "already added contacts")
		ClientManager.EnqueueMessage(msg.uid, resp)
		return nil
	}

	// add to self
	friend, err := dao.UserDao.AddContacts(msg.uid, request.Uid, dao.ContactsTypeUser, request.Remark)
	if err != nil {
		return err
	}
	resp := entity.NewMessage2(msg.seq, entity.ActionSuccess, friend)
	ClientManager.EnqueueMessage(msg.uid, resp)

	// add to friend
	f, err := dao.UserDao.AddContacts(request.Uid, msg.uid, dao.ContactsTypeUser, "")
	if err != nil {
		return err
	}
	resp1 := entity.NewMessage2(-1, entity.ActionUserAddFriend, f)
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
	ret := make([]u, 0, len(users))
	for _, user := range users {
		retU := u{
			Uid:      user.Uid,
			Account:  user.Account,
			Avatar:   user.Avatar,
			Nickname: user.Nickname,
		}
		ret = append(ret, retU)
	}
	if err = resp.SetData(ret); err != nil {
		return err
	}

	ClientManager.EnqueueMessage(msg.uid, resp)
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
	allClient := ClientManager.AllClient()
	users := make([]u, len(allClient))

	ClientManager.Update()
	for k := range allClient {
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

	// todo remove
	chat, err := dao.ChatDao.GetChat(target, request.Type)

	if err != nil {
		chat, err = dao.ChatDao.CreateChat(request.Type, target)
		if err != nil {
			return err
		}
	}

	if request.Type == dao.ChatTypeUser {
		m2, err2 := dao.ChatDao.NewUserChat(chat.Cid, uid, target, dao.ChatTypeUser)
		if err2 != nil {
			return err2
		}
		_, err = dao.ChatDao.NewUserChat(chat.Cid, target, uid, dao.ChatTypeUser)
		if err != nil {
			return err
		}
		resp := entity.NewMessage(msg.seq, entity.ActionSuccess)
		if err = resp.SetData(m2); err != nil {
			return err
		}
		ClientManager.EnqueueMessage(msg.uid, resp)
	} else if request.Type == dao.ChatTypeGroup {
		m, e := dao.ChatDao.NewUserChat(chat.Cid, uid, target, dao.ChatTypeGroup)
		if e != nil {
			return e
		}
		resp := entity.NewMessage2(msg.seq, entity.ActionUserChatInfo, m)
		ClientManager.EnqueueMessage(msg.uid, resp)
	} else {
		return errors.New("unknown chat type")
	}
	return nil
}

func (a *userApi) GetUserChatList(msg *ApiMessage) error {

	resp := entity.NewMessage(msg.seq, entity.ActionSuccess)
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

	resp := entity.NewMessage(msg.seq, entity.ActionSuccess)
	err := dao.UserDao.AddUser(registerEntity.Account, registerEntity.Password)

	if err != nil {
		_ = resp.SetData(err.Error())
	} else {
		_ = resp.SetData("register success")
	}

	return resp, err
}
