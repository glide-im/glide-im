package auth

import (
	"github.com/pkg/errors"
	"go_im/im/api/apidep"
	"go_im/im/api/comm"
	"go_im/im/api/router"
	"go_im/im/dao/userdao"
	"go_im/im/message"
	"go_im/pkg/logger"
	"math/rand"
	"time"
)

var avatars = []string{
	"https://dengzii.com/static/a.webp",
	"https://dengzii.com/static/b.webp",
	"https://dengzii.com/static/c.webp",
	"https://dengzii.com/static/d.webp",
	"https://dengzii.com/static/e.webp",
	"https://dengzii.com/static/f.webp",
	"https://dengzii.com/static/g.webp",
	"https://dengzii.com/static/h.webp",
	"https://dengzii.com/static/i.webp",
	"https://dengzii.com/static/j.webp",
	"https://dengzii.com/static/k.webp",
	"https://dengzii.com/static/l.webp",
	"https://dengzii.com/static/m.webp",
	"https://dengzii.com/static/n.webp",
	"https://dengzii.com/static/o.webp",
	"https://dengzii.com/static/p.webp",
	"https://dengzii.com/static/q.webp",
	"https://dengzii.com/static/r.webp",
}

var nicknames = []string{"佐菲", "赛文", "杰克", "艾斯", "泰罗", "雷欧", "阿斯特拉", "艾迪", "迪迦", "杰斯", "奈克斯", "梦比优斯", "盖亚", "戴拿"}

type Interface interface {
	AuthToken(info *route.Context, req *AuthTokenRequest) error
	SignIn(info *route.Context, req *SignInRequest) error
	Logout(info *route.Context) error
	Register(info *route.Context, req *RegisterRequest) error
}

type AuthApi struct {
}

func (*AuthApi) AuthToken(info *route.Context, req *AuthTokenRequest) error {
	uid, device, err := userdao.Dao.GetTokenInfo(req.Token)
	if err != nil {
		return err
	}
	if uid == 0 {
		info.Response(message.NewMessage(info.Seq, comm.ActionFailed, "token is invalid, plz sign in"))
		return nil
	}
	if req.Device != device {

	}
	apidep.ClientManager.ClientSignIn(info.Uid, uid, device)
	info.Response(message.NewMessage(info.Seq, comm.ActionSuccess, AuthResponse{Uid: uid}))
	return nil
}

func (*AuthApi) SignIn(ctx *route.Context, request *SignInRequest) error {
	if len(request.Account) == 0 || len(request.Password) == 0 {
		return errors.New("account or password empty")
	}
	uid, err := userdao.Dao.GetUidInfoByLogin(request.Account, request.Password)
	if err != nil {
		return err
	}
	if uid == 0 {
		return comm.NewApiBizError(1001, "check your account and password")
	}

	if apidep.ClientManager.IsDeviceOnline(uid, request.Device) {
		err = userdao.Dao.DelAuthToken(uid, request.Device)
		if err != nil {
			logger.E("del user token failed %v", err)
		}
		notify := message.NewMessage(0, message.ActionNotify, "your account has sign in on another device")
		apidep.ClientManager.EnqueueMessage(uid, request.Device, notify)
		apidep.ClientManager.ClientLogout(uid, request.Device)
	}

	token := genToken(32)
	err = userdao.Dao.SetSignInToken(uid, request.Device, token, time.Hour*24*7)
	if err != nil {
		return comm.NewUnexpectedErr("login failed cause internal server error", err)
	}
	resp := message.NewMessage(ctx.Seq, comm.ActionSuccess, AuthResponse{Token: token, Uid: uid})
	apidep.ClientManager.ClientSignIn(ctx.Uid, uid, request.Device)
	ctx.Response(resp)
	return nil
}

func (*AuthApi) Register(ctx *route.Context, req *RegisterRequest) error {
	u := &userdao.User{
		Account:  req.Account,
		Password: req.Password,
		Nickname: nicknames[rand.Intn(len(nicknames))],
		Avatar:   avatars[rand.Intn(len(avatars))],
	}
	err := userdao.UserInfoDao.AddUser(u)
	if err != nil {
		return comm.NewUnexpectedErr("register failed cause internal server error", err)
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, ""))
	return err
}

func (a *AuthApi) Logout(ctx *route.Context, r *LogoutRequest) error {
	err := userdao.Dao.DelAuthToken(ctx.Uid, ctx.Device)
	if err != nil {
		return comm.NewUnexpectedErr("logout failed due to internal server error", err)
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, ""))
	apidep.ClientManager.ClientLogout(ctx.Uid, ctx.Device)
	return nil
}
