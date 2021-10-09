package api

import (
	"errors"
	"go_im/im/client"
	"go_im/im/dao"
	"go_im/im/message"
	"go_im/pkg/logger"
)

type UserApi struct{}

func (a *UserApi) Auth(msg *RequestInfo, request *AuthRequest) error {

	var resp = message.NewMessage(msg.Seq, ActionSuccess, "success")
	uid := dao.UserDao.GetUid(request.Token)
	if uid > 0 {
		client.Manager.ClientSignIn(msg.Uid, uid, request.DeviceId)
		respondMessage(uid, resp)
		return nil
	} else {
		return errors.New("login failed")
	}
}

func (a *UserApi) Login(msg *RequestInfo, request *LoginRequest) error {

	if len(request.Account) == 0 || len(request.Password) == 0 {
		return errors.New("account or password empty")
	}

	uid, token, err := dao.UserDao.GetUidByLogin(request.Account, request.Password)
	if err != nil {
		return err
	}

	m := message.NewMessage(msg.Seq, ActionSuccess, "success")
	if err = m.SetData(AuthorResponse{Token: token, Uid: uid}); err != nil {
		return err
	}
	client.Manager.ClientSignIn(msg.Uid, uid, request.Device)
	respondMessage(uid, m)
	return nil
}

//goland:noinspection GoPreferNilSlice
func (a *UserApi) GetAndInitRelationList(msg *RequestInfo) error {

	allContacts, err := dao.UserDao.GetAllContacts(msg.Uid)
	if err != nil {
		return err
	}

	friends := []*UserInfoResponse{}
	groups := []*GroupResponse{}

	var uids []int64
	for _, contacts := range allContacts {

		if contacts.Type == dao.ContactsTypeGroup {
			g, er := dao.GroupDao.GetGroup(contacts.TargetId)
			if er != nil {
				return er
			}
			if g == nil {
				return errors.New("group not exist")
			}
			members, err := dao.GroupDao.GetMembers(g.Gid)
			if err != nil {
				return err
			}
			groups = append(groups, &GroupResponse{
				Group:   *g,
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
			friends = append(friends, &UserInfoResponse{
				Uid:      u.Uid,
				Account:  u.Account,
				Nickname: u.Nickname,
				Avatar:   u.Avatar,
			})
		}
	}

	body := ContactResponse{
		Friends: friends,
		Groups:  groups,
	}

	resp := message.NewMessage(msg.Seq, ActionSuccess, body)
	respondMessage(msg.Uid, resp)
	return nil
}

func (a *UserApi) AddFriend(msg *RequestInfo, request *AddContacts) error {

	hasUser, err := dao.UserDao.HasUser(request.Uid)
	if err != nil {
		return err
	}
	if !hasUser {
		return errors.New("user not exist")
	}

	hasContacts, err := dao.UserDao.HasContacts(msg.Uid, request.Uid, dao.ContactsTypeUser)
	if err != nil {
		return err
	}

	if hasContacts {
		return errors.New("already added contacts")
	}

	// add to self
	_, err = dao.UserDao.AddContacts(msg.Uid, request.Uid, dao.ContactsTypeUser, request.Remark)
	if err != nil {
		return err
	}

	userInfos, err := dao.UserDao.GetUser(msg.Uid, request.Uid)
	var me *dao.User
	var friend *dao.User

	if userInfos[0].Uid == msg.Uid {
		me = userInfos[0]
		friend = userInfos[1]
	} else {
		me = userInfos[1]
		friend = userInfos[0]
	}
	if err != nil {
		return err
	}

	ccontactResponse := ContactResponse{
		Friends: []*UserInfoResponse{{
			Uid:      friend.Uid,
			Nickname: friend.Nickname,
			Account:  friend.Account,
			Avatar:   friend.Avatar,
		}},
		Groups: []*GroupResponse{},
	}
	respondMessage(msg.Uid, message.NewMessage(msg.Seq, ActionSuccess, ccontactResponse))

	// add to friend
	_, err = dao.UserDao.AddContacts(request.Uid, msg.Uid, dao.ContactsTypeUser, "")
	if err != nil {
		return err
	}

	contactRespFriend := ContactResponse{
		Friends: []*UserInfoResponse{{
			Uid:      msg.Uid,
			Nickname: me.Nickname,
			Account:  me.Account,
			Avatar:   me.Avatar,
		}},
		Groups: []*GroupResponse{},
	}
	respondMessage(request.Uid, message.NewMessage(-1, ActionUserAddFriend, contactRespFriend))

	return nil
}

func (a *UserApi) GetUserInfo(msg *RequestInfo, request *UserInfoRequest) error {

	users, err := dao.UserDao.GetUser(request.Uid...)
	if err != nil {
		return err
	}
	resp := message.NewMessage(msg.Seq, ActionOnlineUser, "success")
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

	respondMessage(msg.Uid, resp)
	return nil
}

func (a *UserApi) GetChatInfo(msg *RequestInfo, request *ChatInfoRequest) error {

	uc, err := dao.ChatDao.GetUserChatFromChat(request.Cid, msg.Uid)
	if err != nil {
		return err
	}
	resp := message.NewMessage(msg.Seq, ActionUserChatInfo, uc)
	respondMessage(msg.Uid, resp)
	return nil
}

func (a *UserApi) GetChatHistory(msg *RequestInfo, request *ChatHistoryRequest) error {

	chatMessages, err := dao.ChatDao.GetChatHistory(request.Cid, 20)
	if err != nil {
		return err
	}

	resp := message.NewMessage(msg.Seq, ActionUserChatHistory, chatMessages)

	respondMessage(msg.Uid, resp)
	return nil
}

func (a *UserApi) GetOnlineUser(msg *RequestInfo) error {

	type u struct {
		Uid      int64
		Account  string
		Avatar   string
		Nickname string
	}
	allClient := client.Manager.AllClient()
	users := make([]u, len(allClient))

	for _, k := range allClient {
		us, err := dao.UserDao.GetUser(k)
		if err != nil || len(us) == 0 {
			logger.D("get online uid=%d error, error=%v", k, err)
			continue
		}
		user := us[0]
		users = append(users, u{Uid: user.Uid, Account: user.Account, Avatar: user.Avatar, Nickname: user.Nickname})
	}

	m := message.NewMessage(msg.Seq, ActionOnlineUser, users)
	respondMessage(msg.Uid, m)
	return nil
}

func (a *UserApi) NewChat(msg *RequestInfo, request *UserNewChatRequest) error {

	uid := msg.Uid
	target := request.Id

	// todo remove
	chat, err := dao.ChatDao.GetChatByTarget(target, request.Type)

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
		resp := message.NewMessage(msg.Seq, ActionSuccess, m2)
		respondMessage(msg.Uid, resp)
	} else if request.Type == dao.ChatTypeGroup {
		m, e := dao.ChatDao.NewUserChat(chat.Cid, uid, target, dao.ChatTypeGroup)
		if e != nil {
			return e
		}
		resp := message.NewMessage(msg.Seq, ActionUserChatInfo, m)
		respondMessage(msg.Uid, resp)
	} else {
		return errors.New("unknown chat type")
	}
	return nil
}

func (a *UserApi) GetUserChatList(msg *RequestInfo) error {

	list, err := dao.ChatDao.GetUserChatList(msg.Uid)
	if err != nil {
		return err
	}
	resp := message.NewMessage(msg.Seq, ActionSuccess, list)
	respondMessage(msg.Uid, resp)
	return nil
}

func (a *UserApi) UserInfo(msg *RequestInfo) error {

	return nil
}

func (a *UserApi) Register(msg *RequestInfo, registerEntity *RegisterRequest) error {

	resp := message.NewMessage(msg.Seq, ActionSuccess, "success")
	err := dao.UserDao.AddUser(registerEntity.Account, registerEntity.Password)

	if err != nil {
		resp = message.NewMessage(msg.Seq, ActionFailed, err)
	}
	respondMessage(msg.Uid, resp)
	return err
}
