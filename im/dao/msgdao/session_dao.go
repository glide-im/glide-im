package msgdao

import (
	"github.com/go-redis/redis"
	"go_im/pkg/db"
	"go_im/pkg/logger"
	"strconv"
	"strings"
	"time"
)

const (
	keySession      = "im:ses:"
	keyUserSessions = "im:user:ses:"
)

var SessionDaoImpl SessionDao = &sessionDaoImpl{}

type sessionDaoImpl struct{}

func getSessionId(uid1 int64, uid2 int64) (string, int64, int64) {
	lg, sm := uid1, uid2
	if lg < sm {
		lg, sm = sm, lg
	}
	return strconv.FormatInt(lg, 10) + "_" + strconv.FormatInt(sm, 10), lg, sm
}

func sid2Uid(sid string) (int64, int64) {
	split := strings.Split(sid, "_")
	if len(split) != 2 {
		logger.E("split sid to uid error:%s", sid)
		return 0, 0
	}
	lg, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		logger.E("split sid to uid error:%v", err)
		return 0, 0
	}
	sm, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		logger.E("split sid to uid error:%v", err)
		return 0, 0
	}
	return lg, sm
}

func getInt64FromMap(m map[string]string, field string) int64 {
	s, ok := m[field]
	if !ok {
		logger.E("get field from redis session map not exist:%s", field)
		return 0
	}
	parseInt, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logger.E("parse session redis map failed:%v", err)
		return 0
	}
	return parseInt
}

func addToUserSessionList(uid int64, sid string, updateAt int64) error {
	_, err := db.Redis.ZAdd(keyUserSessions+strconv.FormatInt(uid, 10), redis.Z{
		Score:  float64(updateAt),
		Member: sid,
	}).Result()
	return err
}

func (s *sessionDaoImpl) CleanUserSessionUnread(uid1, uid2 int64, uid int64) error {
	id, lg, _ := getSessionId(uid1, uid2)
	key := "lg_unread"
	if lg != uid {
		key = "sm_unread"
	}
	_, err := db.Redis.HSet(keySession+id, key, 0).Result()
	return err
}

func (s *sessionDaoImpl) GetSession(uid int64, uid2 int64) (*Session, error) {
	sid, lg, sm := getSessionId(uid, uid2)
	result, err := db.Redis.HGetAll(keySession + sid).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		logger.E("get session from redis is empty sid=%s", sid)
		return s.CreateSession(uid, uid2, time.Now().Unix())
	}
	se := Session{
		SessionId:   sid,
		Uid:         lg,
		Uid2:        sm,
		LgUidUnread: getInt64FromMap(result, "lg_unread"),
		SmUidUnread: getInt64FromMap(result, "sm_unread"),
		LastMID:     getInt64FromMap(result, "l_mid"),
		UpdateAt:    getInt64FromMap(result, "update"),
		CreateAt:    getInt64FromMap(result, "create"),
	}
	return &se, nil
}

func (s *sessionDaoImpl) CreateSession(uid1 int64, uid2 int64, updateAt int64) (*Session, error) {
	sid, lg, sm := getSessionId(uid1, uid2)
	result, err := db.Redis.Exists(keySession + sid).Result()
	if err != nil {
		return nil, err
	}
	if result == 1 {
		db.Redis.Del(keySession + sid)
	}

	_, err = db.Redis.HMSet(keySession+sid, map[string]interface{}{
		"lg_unread": 0,
		"sm_unread": 0,
		"l_mid":     "0",
		"update":    updateAt,
		"create":    updateAt,
	}).Result()
	db.Redis.ExpireAt(keySession+sid, time.Now().Add(time.Hour*24*30))

	if err != nil {
		return nil, err
	}
	err = addToUserSessionList(uid1, sid, updateAt)
	if err != nil {
		return nil, err
	}
	err = addToUserSessionList(uid2, sid, updateAt)
	if err != nil {
		return nil, err
	}

	session := Session{
		SessionId:   sid,
		Uid:         lg,
		Uid2:        sm,
		LgUidUnread: 0,
		SmUidUnread: 0,
		LastMID:     0,
		UpdateAt:    updateAt,
		CreateAt:    updateAt,
	}
	return &session, nil
}

func (s *sessionDaoImpl) UpdateOrCreateSession(uid1 int64, uid2 int64, sender int64, mid int64, sendAt int64) error {
	sid, lg, _ := getSessionId(uid1, uid2)
	result, err := db.Redis.Exists(keySession + sid).Result()
	if err != nil {
		return err
	}
	if result == 0 {
		err := addToUserSessionList(uid1, sid, sendAt)
		if err != nil {
			return err
		}
		err = addToUserSessionList(uid2, sid, sendAt)
		if err != nil {
			return err
		}
	}

	key := "lg_unread"
	if sender != lg {
		key = "sm_unread"
	}

	m := map[string]interface{}{
		"l_mid": mid,
		// TODO 2021-12-28 NOTE: do business here, clean sender's unread
		key:      0,
		"update": sender,
	}
	_, err = db.Redis.HMSet(keySession+sid, m).Result()
	db.Redis.ExpireAt(keySession+sid, time.Now().Add(time.Hour*24*30))

	if err != nil {
		return err
	}

	if sender == lg {
		_, err := db.Redis.HIncrBy(keySession+sid, "sm_unread", 1).Result()
		if err != nil {
			return err
		}
	} else {
		_, err := db.Redis.HIncrBy(keySession+sid, "lg_unread", 1).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sessionDaoImpl) GetRecentSession(uid int64, before int64, pageSize int64) ([]*Session, error) {
	var se []*Session

	result, err := db.Redis.ZRevRangeByScore(keyUserSessions+strconv.FormatInt(uid, 10), redis.ZRangeBy{
		Min:    "0",
		Max:    strconv.FormatInt(before, 10),
		Offset: 0,
		Count:  pageSize,
	}).Result()
	if err != nil {
		return nil, err
	}
	for _, sid := range result {
		uid1, uid2 := sid2Uid(sid)
		if uid1 == 0 {
			continue
		}
		session, err := s.GetSession(uid1, uid2)
		if err != nil {
			return nil, err
		}
		se = append(se, session)
	}

	return se, nil
}
