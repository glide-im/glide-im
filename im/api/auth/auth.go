package auth

import (
	"go_im/im/api/apidep"
	"go_im/im/api/comm"
	"go_im/im/api/router"
	"go_im/im/dao/common"
	"go_im/im/dao/userdao"
	"go_im/im/message"
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

var (
	ErrInvalidToken      = comm.NewApiBizError(1001, "token is invalid, plz sign in")
	ErrSignInAccountInfo = comm.NewApiBizError(1002, "check your account and password")
	ErrReplicatedLogin   = comm.NewApiBizError(1003, "replicated login")
)

type AuthApi struct {
}

func (*AuthApi) AuthToken(ctx *route.Context, req *AuthTokenRequest) error {
	token, err := comm.ParseJwt(req.Token)
	if err != nil {
		return ErrInvalidToken
	}
	version, err := userdao.Dao.GetTokenVersion(token.Uid, token.Device)
	if err != nil || version == 0 || version > token.Ver {
		return ErrInvalidToken
	}
	if ctx.Uid == token.Uid && ctx.Device == token.Device {
		return ErrReplicatedLogin
	}
	if ctx.Uid != 0 {
		apidep.ClientManager.ClientSignIn(ctx.Uid, token.Uid, token.Device)
		ctx.Uid = token.Uid
		ctx.Device = token.Device
	}
	resp := AuthResponse{
		Uid: token.Uid,
		Servers: []string{
			"ws://192.168.1.123:8080/ws",
		},
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, resp))
	return nil
}

func (*AuthApi) SignIn(ctx *route.Context, request *SignInRequest) error {
	if len(request.Account) == 0 || len(request.Password) == 0 {
		return ErrSignInAccountInfo
	}
	uid, err := userdao.Dao.GetUidInfoByLogin(request.Account, request.Password)
	if err != nil || uid == 0 {
		if err == common.ErrNoRecordFound || uid == 0 {
			return ErrSignInAccountInfo
		}
		return comm.NewDbErr(err)
	}
	jt := comm.AuthInfo{
		Uid:    uid,
		Device: request.Device,
		Ver:    comm.GenJwtVersion(),
	}
	token, err := comm.GenJwt(jt)
	if err != nil {
		return comm.NewUnexpectedErr("login failed", err)
	}
	err = userdao.Dao.SetTokenVersion(jt.Uid, jt.Device, jt.Ver, time.Duration(jt.ExpiresAt))
	if err != nil {
		return comm.NewDbErr(err)
	}

	tk := AuthResponse{
		Uid:   uid,
		Token: token,
		Servers: []string{
			"ws://192.168.1.123:8080/ws",
		},
	}
	resp := message.NewMessage(ctx.Seq, comm.ActionSuccess, tk)
	apidep.ClientManager.ClientSignIn(ctx.Uid, uid, request.Device)

	ctx.Uid = uid
	ctx.Device = request.Device
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
		return comm.NewDbErr(err)
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, ""))
	return err
}

func (a *AuthApi) Logout(ctx *route.Context) error {
	err := userdao.Dao.DelAuthToken(ctx.Uid, ctx.Device)
	if err != nil {
		return comm.NewDbErr(err)
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, ""))
	apidep.ClientManager.ClientLogout(ctx.Uid, ctx.Device)
	return nil
}

func (a *AuthApi) offline(uid int64, device int64) {
	if apidep.ClientManager.IsDeviceOnline(uid, device) {
		//err := userdao.Dao.DelAuthToken(uid, device)
		//if err != nil {
		//	logger.E("del user token failed %v", err)
		//}
		notify := message.NewMessage(0, message.ActionNotify, "your account has sign in on another device")
		apidep.SendMessage(uid, device, notify)
		apidep.ClientManager.ClientLogout(uid, device)
	}
}
