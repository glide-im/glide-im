package uid

import (
	"go_im/pkg/db"
	"go_im/pkg/logger"
	"math/rand"
	"strconv"
	"sync"
)

const keyTempIdIncr = "im:uid:temp:incr"
const keyUidIncr = "im:uid:incr"
const keySystemIdIncr = "im:uid:sys:incr"

const (
	systemIdStart = 1_000
	systemIdEnd   = systemIdStart + 100_000

	userIdStart = systemIdEnd + 100_000
	userIdEnd   = userIdStart + 1_000_000_000_000

	tempIdStart = userIdEnd + 100_000
	tempIdEnd   = tempIdStart + 1_000_000_000
)

func Init() {
	checkIncrKey(keyTempIdIncr, tempIdStart+1)
	checkIncrKey(keySystemIdIncr, systemIdStart+1)
	checkIncrKey(keyUidIncr, userIdStart+1)
}

func checkIncrKey(key string, initialValue int64) {
	result, err := db.Redis.Exists(key).Result()
	if err != nil {
		panic(err)
	}
	if result == 0 {
		db.Redis.Set(key, initialValue, 0)
	}
}

func IsUserId(uid int64) bool {
	return userIdStart < uid && uid < userIdEnd
}

func IsSystemId(uid int64) bool {
	return systemIdStart < uid && uid < systemIdEnd
}

func IsTempId(uid int64) bool {
	return tempIdStart < uid && uid < tempIdEnd
}

type Gen interface {
	GenSysUid() int64
	GenUid() int64
	GenTempUid() int64
}

type gen struct {
	muTempUid sync.Mutex
}

var instance Gen = &gen{
	muTempUid: sync.Mutex{},
}

func (g *gen) GenSysUid() int64 {
	rs, err := g.getInt64(keyTempIdIncr)
	if err != nil {
		return 0
	}
	next := rs + 1
	db.Redis.Set(keySystemIdIncr, next, 0)
	return next
}

func (g *gen) GenUid() int64 {
	rs, err := g.getInt64(keyUidIncr)
	if err != nil {
		return 0
	}
	next := rs + rand.Int63n(4)
	db.Redis.Set(keyUidIncr, next, 0)
	return next
}

func (g *gen) GenTempUid() int64 {
	result, err := db.Redis.Incr(keyTempIdIncr).Result()
	if err != nil {
		return 0
	}
	if result >= tempIdEnd {
		result = tempIdStart
		db.Redis.Set(keyTempIdIncr, result, 0)
	}
	return result
}

func (g *gen) getInt64(key string) (int64, error) {
	r, err := db.Redis.Get(key).Result()
	if err != nil {
		logger.E("get redis error", err)
		return 0, err
	}
	rs, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		logger.E("conv redis value to int64 err", err)
		return 0, err
	}
	return rs, nil
}

func GenSysUid() int64 {
	return instance.GenSysUid()
}

func GenUid() int64 {
	return instance.GenUid()
}

func GenTemp() int64 {
	return instance.GenTempUid()
}
