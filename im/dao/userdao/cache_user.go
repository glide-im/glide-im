package userdao

import (
	"fmt"
	"go_im/pkg/db"
	"go_im/pkg/logger"
	"strconv"
	"time"
)

var keyTokenPre = "im:auth:token:"

const tokenExpired = time.Hour * 24 * 3

type UserCacheDao struct {
}

func (UserCacheDao) GetUserLoginState(uid int64) ([]*LoginState, error) {
	panic("implement me")
}

func (UserCacheDao) DelUserToken(uid int64, device int64) error {
	panic("implement me")
}

func (UserCacheDao) DelAllToken(uid int64) error {
	k := fmt.Sprintf("%s%d", keyTokenPre, uid)
	_, err := db.Redis.Del(k).Result()
	return err
}

func (UserCacheDao) GetTokenUid(token string) (int64, error) {
	k := fmt.Sprintf("%s%s", keyTokenPre, token)
	result, err := db.Redis.Get(k).Result()
	if err != nil {
		return 0, err
	}
	uid, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		logger.E("redis auth token error", err)
		return 0, err
	}
	return uid, nil
}

func (UserCacheDao) SetUserToken(uid int64, token int64, device int64, expiredAt time.Duration) error {
	k2 := fmt.Sprintf("%s%s", keyTokenPre, token)
	_, err := db.Redis.Set(k2, uid, tokenExpired).Result()

	if err != nil {
		logger.E("redis set auth error", err)
		return err
	}
	return nil
}
