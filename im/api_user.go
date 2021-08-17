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

func (a *userApi) Auth(msg *ApiMessage, request *entity.AuthRequest) error {

	var resp = entity.NewMessage(msg.seq, entity.ActionSuccess)
	uid := dao.UserDao.GetUid(request.Token)
	if uid > 0 {
		ClientManager.ClientSignIn(msg.uid, uid, request.DeviceId)
		ClientManager.EnqueueMessage(uid, resp)
		return nil
	} else {
		return errors.New("login failed")
	}
}

func (a *userApi) Login(msg *ApiMessage, request *entity.LoginRequest) error {

	if len(request.Account) == 0 || len(request.Password) == 0 {
		return errors.New("account or password empty")
	}

	uid, token, err := dao.UserDao.GetUidByLogin(request.Account, request.Password)
	if err != nil {
		return err
	}

	m := entity.NewMessage(msg.seq, entity.ActionSuccess)
	if err = m.SetData(entity.AuthorResponse{Token: token, Uid: uid}); err != nil {
		return err
	}
	ClientManager.ClientSignIn(msg.uid, uid, request.Device)
	ClientManager.EnqueueMessage(uid, m)
	return nil
}

//goland:noinspection GoPreferNilSlice
func (a *userApi) GetAndInitRelationList(msg *ApiMessage) error {

	allContacts, err := dao.UserDao.GetAllContacts(msg.uid)
	if err != nil {
		return err
	}

	friends := []*entity.UserInfoResponse{}
	groups := []*entity.GroupResponse{}

	var uids []int64
	for _, contacts := range allContacts {

		if contacts.Type == dao.ContactsTypeGroup {
			group := GroupManager.GetGroup(contacts.TargetId)
			if group == nil {
				return newApiFatalError("load user group error: nil")
			}
			members, err := GroupManager.GetMembers(group.Gid)
			if err != nil {
				return err
			}
			groups = append(groups, &entity.GroupResponse{
				Group:   *group.group,
				Members: members,
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

func (a *userApi) AddFriend(msg *ApiMessage, request *entity.AddContacts) error {

	hasUser, err := dao.UserDao.HasUser(request.Uid)
	if err != nil {
		return err
	}
	if !hasUser {
		ClientManager.EnqueueMessage(msg.uid, entity.NewErrMessage2(msg.seq, "user not exist"))
		return nil
	}

	hasContacts, err := dao.UserDao.HasContacts(msg.uid, request.Uid, dao.ContactsTypeUser)
	if err != nil {
		return err
	}

	if hasContacts {
		ClientManager.EnqueueMessage(msg.uid, entity.NewErrMessage2(msg.seq, "already added contacts"))
		return nil
	}

	// add to self
	_, err = dao.UserDao.AddContacts(msg.uid, request.Uid, dao.ContactsTypeUser, request.Remark)
	if err != nil {
		return err
	}

	userInfos, err := dao.UserDao.GetUser(msg.uid, request.Uid)
	var me *dao.User
	var friend *dao.User

	if userInfos[0].Uid == msg.uid {
		me = userInfos[0]
		friend = userInfos[1]
	} else {
		me = userInfos[1]
		friend = userInfos[0]
	}
	if err != nil {
		return err
	}

	ccontactResponse := entity.ContactResponse{
		Friends: []*entity.UserInfoResponse{{
			Uid:      friend.Uid,
			Nickname: friend.Nickname,
			Account:  friend.Account,
			Avatar:   friend.Avatar,
		}},
		Groups: []*entity.GroupResponse{},
	}
	ClientManager.EnqueueMessage(msg.uid, entity.NewMessage2(msg.seq, entity.ActionSuccess, ccontactResponse))

	// add to friend
	_, err = dao.UserDao.AddContacts(request.Uid, msg.uid, dao.ContactsTypeUser, "")
	if err != nil {
		return err
	}

	contactRespFriend := entity.ContactResponse{
		Friends: []*entity.UserInfoResponse{{
			Uid:      msg.uid,
			Nickname: me.Nickname,
			Account:  me.Account,
			Avatar:   me.Avatar,
		}},
		Groups: []*entity.GroupResponse{},
	}
	ClientManager.EnqueueMessage(request.Uid, entity.NewMessage2(-1, entity.ActionUserAddFriend, contactRespFriend))

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
	for _, k := range allClient {
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

func (a *userApi) Register(msg *ApiMessage, registerEntity *entity.RegisterRequest) error {

	resp := entity.NewMessage(msg.seq, entity.ActionSuccess)
	err := dao.UserDao.AddUser(registerEntity.Account, registerEntity.Password)

	if err != nil {
		_ = resp.SetData(err.Error())
	} else {
		_ = resp.SetData("register success")
	}
	ClientManager.EnqueueMessage(msg.uid, resp)
	return err
}
