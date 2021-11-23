package auth

import (
	"github.com/pkg/errors"
	"go_im/im/api/apidep"
	"go_im/im/api/comm"
	"go_im/im/api/router"
	"go_im/im/dao/userdao"
	"go_im/im/message"
)

type Interface interface {
	AuthToken(info *route.Context, req *AuthTokenReq) error
	SignIn(info *route.Context, req *SignInRequest) error
	Logout(info *route.Context) error
	Register(info *route.Context, req *RegisterRequest) error
}

type AuthApi struct {
}

func (*AuthApi) AuthToken(info *route.Context, req *AuthTokenReq) error {
	panic("implement me")
}

func (*AuthApi) SignIn(ctx *route.Context, request *SignInRequest) error {
	if len(request.Account) == 0 || len(request.Password) == 0 {
		return errors.New("account or password empty")
	}

	uid, token, err := userdao.UserDao2.GetUidByLogin(request.Account, request.Password, request.Device)
	if err != nil {
		return err
	}

	m := message.NewMessage(ctx.Seq, "", "success")
	if err = m.SetData(AuthorResponse{Token: token, Uid: uid}); err != nil {
		return err
	}
	apidep.ClientManager.ClientSignIn(ctx.Uid, uid, request.Device)
	ctx.Response(m)
	return nil
}

func (*AuthApi) Register(ctx *route.Context, registerEntity *RegisterRequest) error {

	err := userdao.UserDao2.AddUser(registerEntity.Account, registerEntity.Password)
	if err != nil {
		return err
	}
	ctx.Response(message.NewMessage(ctx.Seq, comm.ActionSuccess, "success"))
	return err
}

func (a *AuthApi) Logout(ctx *route.Context, r *LogoutRequest) error {
	err := userdao.UserDao2.Logout(ctx.Uid, r.Device, r.Token)
	if err != nil {
		return err
	}
	ctx.Response(message.NewMessage(ctx.Seq, "", ""))
	apidep.ClientManager.ClientLogout(ctx.Uid, r.Device)
	return nil
}

func (a *AuthApi) Auth(ctx *route.Context, request *AuthRequest) error {

	var resp = message.NewMessage(ctx.Seq, "", "success")
	uid := userdao.UserDao2.GetUid(request.Token)
	if uid > 0 {
		apidep.ClientManager.ClientSignIn(ctx.Uid, uid, request.DeviceId)
		ctx.Response(resp)
		return nil
	} else {
		return errors.New("token expired")
	}
}
