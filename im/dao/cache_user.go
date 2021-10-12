package dao

import (
	"fmt"
	"go_im/pkg/db"
	"go_im/pkg/logger"
	"strconv"
	"time"
)

var keyTokenPre = "im:auth:token:"

const tokenExpired = time.Hour * 24 * 3

func setAuthToken(uid int64, token string, device int64) error {

	k2 := fmt.Sprintf("%s%s", keyTokenPre, token)
	_, err := db.Redis.Set(k2, uid, tokenExpired).Result()

	if err != nil {
		logger.E("redis set auth error", err)
		return err
	}
	return nil
}

func authToken(token string) int64 {
	k := fmt.Sprintf("%s%s", keyTokenPre, token)
	result, err := db.Redis.Get(k).Result()
	if err != nil {
		return 0
	}
	uid, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		logger.E("redis auth token error", err)
		return 0
	}
	return uid
}

func delAuthToken(token string) error {
	k := fmt.Sprintf("%s%s", keyTokenPre, token)
	_, err := db.Redis.Del(k).Result()
	return err
}
