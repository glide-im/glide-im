package user

import (
	"errors"
	"go_im/im/api/apidep"
	"go_im/im/api/router"
	"go_im/im/dao/groupdao"
	"go_im/im/dao/userdao"
	"go_im/im/message"
	"go_im/pkg/logger"
)

type UserApi struct{}

//goland:noinspection GoPreferNilSlice
func (a *UserApi) GetContactList(msg *route.Context) error {

	allContacts, err := userdao.UserDao.GetAllContacts(msg.Uid)
	if err != nil {
		return err
	}

	friends := []*InfoResponse{}
	groups := []interface{}{}

	var uids []int64
	for _, contacts := range allContacts {

		if contacts.Type == userdao.ContactsTypeGroup {
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
		} else if contacts.Type == userdao.ContactsTypeUser {
			uids = append(uids, contacts.TargetId)
		}
	}
	if len(uids) > 0 {
		user, err := userdao.UserDao.GetUser(uids...)
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
	msg.Response(resp)
	return nil
}

func (a *UserApi) AddContact(msg *route.Context, request *AddContacts) error {

	hasUser, err := userdao.UserDao.HasUser(request.Uid)
	if err != nil {
		return err
	}
	if !hasUser {
		return errors.New("user not exist")
	}

	hasContacts, err := userdao.UserDao.HasContacts(msg.Uid, request.Uid, userdao.ContactsTypeUser)
	if err != nil {
		return err
	}

	if hasContacts {
		return errors.New("already added contacts")
	}

	// add to self
	_, err = userdao.UserDao.AddContacts(msg.Uid, request.Uid, userdao.ContactsTypeUser, request.Remark)
	if err != nil {
		return err
	}

	userInfos, err := userdao.UserDao.GetUser(msg.Uid, request.Uid)
	var me *userdao.User
	var friend *userdao.User

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
		Groups: []interface{}{},
	}
	msg.Response(message.NewMessage(msg.Seq, "api.ActionSuccess", ccontactResponse))

	// add to friend
	_, err = userdao.UserDao.AddContacts(request.Uid, msg.Uid, userdao.ContactsTypeUser, "")
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
		Groups: []interface{}{},
	}
	msg.Response(message.NewMessage(-1, "api.ActionUserAddFriend", contactRespFriend))

	return nil
}

func (a *UserApi) GetUserInfo(msg *route.Context, request *InfoRequest) error {

	users, err := userdao.UserDao.GetUser(request.Uid...)
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

	msg.Response(resp)
	return nil
}

func (a *UserApi) GetOnlineUser(msg *route.Context) error {

	type u struct {
		Uid      int64
		Account  string
		Avatar   string
		Nickname string
	}
	allClient := apidep.ClientManager.AllClient()
	users := make([]u, len(allClient))

	for _, k := range allClient {
		us, err := userdao.UserDao.GetUser(k)
		if err != nil || len(us) == 0 {
			logger.D("get online uid=%d error, error=%v", k, err)
			continue
		}
		user := us[0]
		users = append(users, u{Uid: user.Uid, Account: user.Account, Avatar: user.Avatar, Nickname: user.Nickname})
	}

	m := message.NewMessage(msg.Seq, "api.ActionOnlineUser", users)
	msg.Response(m)
	return nil
}

func (a *UserApi) UserInfo(msg *route.Context) error {

	return nil
}
