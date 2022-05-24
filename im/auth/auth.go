package auth

import (
	"fmt"
	"go_im/im/api/comm"
	"go_im/im/dao/userdao"
	"go_im/pkg/logger"
	"time"
)

type Token struct {
	Token string
}

type Result struct {
	Uid     int64
	Token   string
	Servers []string
}

func ParseToken(token string) (*AuthInfo, error) {
	return parseJwt(token)
}

func Auth(from int64, device int64, t *Token) (*Result, error) {

	token, err := parseJwt(t.Token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}
	version, err := userdao.Dao.GetTokenVersion(token.Uid, token.Device)
	if err != nil || version == 0 || version > token.Ver {
		return nil, fmt.Errorf("invalid token")
	}
	if from == token.Uid && device == token.Device {
		// logged in
		logger.D("auth token for a connection is logged in")
	}

	return &Result{
		Uid:     token.Uid,
		Token:   t.Token,
		Servers: nil,
	}, nil
}

func GenerateTokenExpire(uid int64, device int64, expire int64) (string, error) {
	jt := AuthInfo{
		Uid:    uid,
		Device: device,
		Ver:    genJwtVersion(),
	}
	expir := time.Now().Add(time.Hour * time.Duration(expire))
	token, err := genJwtExp(jt, expir)
	if err != nil {
		return "", comm.NewUnexpectedErr("login failed", err)
	}

	err = userdao.Dao.SetTokenVersion(jt.Uid, jt.Device, jt.Ver, time.Duration(jt.ExpiresAt))
	if err != nil {
		return "", fmt.Errorf("generate token failed")
	}

	return token, nil
}

func GenerateToken(uid int64, device int64) (string, error) {
	return GenerateTokenExpire(uid, device, 24*3)
}
