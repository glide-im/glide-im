package userdao

import (
	"errors"
	"fmt"
	"go_im/pkg/db"
	"go_im/pkg/logger"
	"strconv"
	"strings"
	"time"
)

var keyToken2Uid = "im:auth:token:"
var keyUid2Token = "im:auth:login:"

type UserCacheDao struct {
}

func (UserCacheDao) IsUserSignIn(uid int64, device int64) (bool, error) {
	s := fmt.Sprintf("%d_%d", device, uid)
	result, err := db.Redis.Exists(keyUid2Token + s).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (UserCacheDao) GetUserSignState(uid int64) ([]*LoginState, error) {
	// TODO 2021-11-25
	panic("implement me")
}

func (UserCacheDao) DelAuthToken(uid int64, device int64) error {
	s := fmt.Sprintf("%d_%d", device, uid)

	db.Redis.Pipeline()
	token, err := db.Redis.Get(keyUid2Token + s).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			logger.D("not signed in")
			return nil
		}
		return err
	}
	if token == "" {
		logger.W("token not exist")
		return nil
	}

	r, err := db.Redis.Del(keyUid2Token + s).Result()
	if r == 0 {
		return nil
	}
	if err != nil {
		return err
	}
	_, err = db.Redis.Del(keyToken2Uid + token).Result()
	if err != nil {
		return err
	}
	return nil
}

func (UserCacheDao) DelToken(token string) error {
	_, err := db.Redis.Del(keyToken2Uid + token).Result()
	if err != nil {
		return err
	}
	return nil
}

func (UserCacheDao) DelAllToken(uid int64) error {
	k := fmt.Sprintf("%s%d", keyToken2Uid, uid)
	_, err := db.Redis.Del(k).Result()
	return err
}

func (UserCacheDao) GetTokenInfo(token string) (int64, int64, error) {
	k := fmt.Sprintf("%s%s", keyToken2Uid, token)
	exist, err := db.Redis.Exists(k).Result()
	if err != nil {
		return 0, 0, err
	}
	if exist != 1 {
		return 0, 0, nil
	}

	result, err := db.Redis.Get(k).Result()
	if err != nil || len(result) == 0 {
		return 0, 0, err
	}
	sp := strings.Split(result, "_")
	if len(sp) != 2 {
		return 0, 0, errors.New("wrong token info")
	}
	deviceId, err := strconv.ParseInt(sp[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	uid, err := strconv.ParseInt(sp[1], 10, 64)
	if err != nil {
		logger.E("redis auth token error", err)
		return 0, 0, err
	}
	return uid, deviceId, nil
}

func (UserCacheDao) SetSignInToken(uid int64, device int64, token string, expiredAt time.Duration) error {
	s := fmt.Sprintf("%d_%d", device, uid)
	_, err := db.Redis.Set(keyToken2Uid+token, s, expiredAt).Result()
	if err != nil {
		logger.E("redis set auth error", err)
		return err
	}
	_, err = db.Redis.Set(keyUid2Token+s, token, expiredAt).Result()
	if err != nil {
		logger.E("redis set auth error", err)
		return err
	}

	return nil
}
