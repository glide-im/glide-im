package comm

import (
	"github.com/golang-jwt/jwt"
	"time"
)

var jwtSecret []byte

func init() {
	// TODO 2021-12-17 11:11:34 update
	jwtSecret = []byte("glide-im-jwt-secret")
}

type AuthInfo struct {
	jwt.StandardClaims
	Uid    int64 `json:"uid"`
	Device int64 `json:"device"`
	Ver    int64 `json:"ver"`
}

func GenJwt(payload AuthInfo) (string, error) {

	expireAt := time.Now().Add(time.Hour * 24)
	payload.ExpiresAt = expireAt.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return t, nil
}

func ParseJwt(token string) (*AuthInfo, error) {
	j := AuthInfo{}
	t, err := jwt.ParseWithClaims(token, &j, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	jwtToken, ok := t.Claims.(*AuthInfo)
	if !ok {
		return nil, NewApiBizError(1, "invalid token")
	}
	return jwtToken, nil
}

func GenJwtVersion() int64 {
	return time.Now().Unix()
}
