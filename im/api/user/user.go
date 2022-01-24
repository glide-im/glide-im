package user

import (
	"errors"
	"go_im/im/api/apidep"
	"go_im/im/api/comm"
	"go_im/im/api/router"
	"go_im/im/dao/userdao"
	"go_im/im/message"
	"go_im/pkg/logger"
)

type UserApi struct{}

func (a *UserApi) GetUserProfile(msg *route.Context) error {
	// TODO 2021-11-29 我的详细信息
	return nil
}

func (a *UserApi) UpdateUserProfile(msg *route.Context, request *UpdateProfileRequest) error {
	// TODO 2021-11-29 更新我的信息
	return nil
}

func (a *UserApi) GetUserInfo(ctx *route.Context, request *InfoRequest) error {
	info, err := userdao.UserInfoDao.GetUserSimpleInfo(request.Uid...)
	if err != nil {
		return comm.NewDbErr(err)
	}
	//goland:noinspection ALL
	resp := []InfoResponse{}
	for _, i := range info {
		resp = append(resp, InfoResponse{
			Uid:      i.Uid,
			Nickname: i.Nickname,
			Account:  i.Account,
			Avatar:   i.Avatar,
		})
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, resp))
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
		us, err := userdao.UserDao2.GetUser(k)
		if err != nil || len(us) == 0 {
			logger.D("get online uid=%d error, error=%v", k, err)
			continue
		}
		user := us[0]
		users = append(users, u{Uid: user.Uid, Account: user.Account, Avatar: user.Avatar, Nickname: user.Nickname})
	}

	m := message.NewMessage(msg.Seq, comm.ActionSuccess, users)
	msg.Response(m)
	return nil
}

func (a *UserApi) UserProfile(ctx *route.Context) error {

	info, err := userdao.UserInfoDao.GetUserSimpleInfo(ctx.Uid)
	if err != nil {
		return comm.NewDbErr(err)
	}
	//goland:noinspection ALL
	resp := []InfoResponse{}
	for _, i := range info {
		resp = append(resp, InfoResponse{
			Uid:      i.Uid,
			Nickname: i.Nickname,
			Account:  i.Account,
			Avatar:   i.Avatar,
		})
	}
	if len(resp) != 1 {
		return comm.NewUnexpectedErr("no such user", errors.New("user info is empty"))
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, resp[0]))
	return nil
}
