package user

import (
	"errors"
	"go_im/im/api/groups"
	"go_im/im/api/router"
	"go_im/im/client"
	"go_im/im/dao"
	"go_im/im/dao/groupdao"
	"go_im/im/message"
	"go_im/pkg/logger"
)

var ResponseHandleFunc func(uid int64, device int64, message *message.Message)

func respond(uid int64, seq int64, action message.Action, data interface{}) {
	resp := message.NewMessage(seq, action, data)
	respondMessage(uid, resp)
}

func respondMessage(uid int64, msg *message.Message) {
	ResponseHandleFunc(uid, 0, msg)
}

type UserApi struct{}

//goland:noinspection GoPreferNilSlice
func (a *UserApi) GetAndInitRelationList(msg *route.RequestInfo) error {

	allContacts, err := dao.UserDao.GetAllContacts(msg.Uid)
	if err != nil {
		return err
	}

	friends := []*InfoResponse{}
	groups := []*groups.GroupResponse{}

	var uids []int64
	for _, contacts := range allContacts {

		if contacts.Type == dao.ContactsTypeGroup {
			g, er := groupdao.GroupDao.GetGroup(contacts.TargetId)
			if er != nil {
				return er
			}
			if g == nil {
				return errors.New("group not exist")
			}
			//members, err := groupdao.GroupDao.GetMembers(g.Gid)
			if err != nil {
				return err
			}
			//groups = append(groups, &groups.GroupResponse{
			//	Group:   *g,
			//	Members: members,
			//})
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
			friends = append(friends, &InfoResponse{
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

	resp := message.NewMessage(msg.Seq, "api.ActionSuccess", body)
	respondMessage(msg.Uid, resp)
	return nil
}

func (a *UserApi) AddFriend(msg *route.RequestInfo, request *AddContacts) error {

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
		Friends: []*InfoResponse{{
			Uid:      friend.Uid,
			Nickname: friend.Nickname,
			Account:  friend.Account,
			Avatar:   friend.Avatar,
		}},
		Groups: []*groups.GroupResponse{},
	}
	respondMessage(msg.Uid, message.NewMessage(msg.Seq, "api.ActionSuccess", ccontactResponse))

	// add to friend
	_, err = dao.UserDao.AddContacts(request.Uid, msg.Uid, dao.ContactsTypeUser, "")
	if err != nil {
		return err
	}

	contactRespFriend := ContactResponse{
		Friends: []*InfoResponse{{
			Uid:      msg.Uid,
			Nickname: me.Nickname,
			Account:  me.Account,
			Avatar:   me.Avatar,
		}},
		Groups: []*groups.GroupResponse{},
	}
	respondMessage(request.Uid, message.NewMessage(-1, "api.ActionUserAddFriend", contactRespFriend))

	return nil
}

func (a *UserApi) GetUserInfo(msg *route.RequestInfo, request *InfoRequest) error {

	users, err := dao.UserDao.GetUser(request.Uid...)
	if err != nil {
		return err
	}
	resp := message.NewMessage(msg.Seq, "api.ActionOnlineUser", "success")
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

func (a *UserApi) GetOnlineUser(msg *route.RequestInfo) error {

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

	m := message.NewMessage(msg.Seq, "api.ActionOnlineUser", users)
	respondMessage(msg.Uid, m)
	return nil
}

func (a *UserApi) NewChat(msg *route.RequestInfo, request *NewChatRequest) error {

	uid := msg.Uid
	target := request.Id

	// todo remove
	chat, err := dao.ChatDao.GetChatByTarget(target, request.Type)

	if err != nil {
		chat, err = dao.ChatDao.CreateChat(request.Type, msg.Uid, target)
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
		resp := message.NewMessage(msg.Seq, "api.ActionSuccess", m2)
		respondMessage(msg.Uid, resp)
	} else if request.Type == dao.ChatTypeGroup {
		m, e := dao.ChatDao.NewUserChat(chat.Cid, uid, target, dao.ChatTypeGroup)
		if e != nil {
			return e
		}
		resp := message.NewMessage(msg.Seq, "api.ActionUserChatInfo", m)
		respondMessage(msg.Uid, resp)
	} else {
		return errors.New("unknown chat type")
	}
	return nil
}

func (a *UserApi) GetUserChatList(msg *route.RequestInfo) error {

	list, err := dao.ChatDao.GetUserChatList(msg.Uid)
	if err != nil {
		return err
	}
	resp := message.NewMessage(msg.Seq, "api.ActionSuccess", list)
	respondMessage(msg.Uid, resp)
	return nil
}

func (a *UserApi) UserInfo(msg *route.RequestInfo) error {

	return nil
}
