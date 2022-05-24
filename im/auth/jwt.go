package auth

import (
	"github.com/golang-jwt/jwt"
	"go_im/im/api/comm"
	"time"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte("glide-im-jwt-secret")
}

type AuthInfo struct {
	jwt.StandardClaims
	Uid    int64 `json:"uid"`
	Device int64 `json:"device"`
	Ver    int64 `json:"ver"`
}

func genJwt(payload AuthInfo) (string, error) {

	expireAt := time.Now().Add(time.Hour * 24)
	return genJwtExp(payload, expireAt)
}

func genJwtExp(payload AuthInfo, expiredAt time.Time) (string, error) {
	payload.ExpiresAt = expiredAt.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return t, nil
}

func parseJwt(token string) (*AuthInfo, error) {
	j := AuthInfo{}
	t, err := jwt.ParseWithClaims(token, &j, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	jwtToken, ok := t.Claims.(*AuthInfo)
	if !ok {
		return nil, comm.NewApiBizError(1, "invalid token")
	}
	return jwtToken, nil
}

func genJwtVersion() int64 {
	return time.Now().Unix()
}
