package uid

import (
	"go_im/pkg/db"
	"math/rand"
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

type gen struct{}

var instance Gen = &gen{}

func (g *gen) GenSysUid() int64 {
	r, err := db.Redis.Get(keySystemIdIncr).Int64()
	if err != nil {
		return 0
	}
	next := r + 1
	db.Redis.Set(keySystemIdIncr, next, 0)
	return next
}

func (g *gen) GenUid() int64 {
	r, err := db.Redis.Get(keyUidIncr).Int64()
	if err != nil {
		return 0
	}
	next := r + rand.Int63n(4)
	db.Redis.Set(keyUidIncr, next, 0)
	return next
}

func (g *gen) GenTempUid() int64 {
	r, err := db.Redis.Get(keyTempIdIncr).Int64()
	if err != nil {
		return 0
	}
	next := r + rand.Int63n(4)
	if next >= tempIdEnd {
		next = tempIdStart
	}
	db.Redis.Set(keyTempIdIncr, next, 0)
	return next
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
