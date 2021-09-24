package uid

import (
	"go_im/pkg/db"
	"math/rand"
)

const KeyTempIdIncr = "im:uid:temp:incr"
const KeyUidIncr = "im:uid:incr"
const KeySystemIdIncr = "im:uid:sys:incr"

const (
	SystemIdStart = 1_000
	SystemIdEnd   = SystemIdStart + 100_000

	UserIdStart = SystemIdEnd + 100_000
	UserIdEnd   = UserIdStart + 1_000_000_000_000

	TempIdStart = UserIdEnd + 100_000
	TempIdEnd   = TempIdStart + 1_000_000_000
)

func IsUserId(uid int64) bool {
	return UserIdStart < uid && uid < UserIdEnd
}

func IsSystemId(uid int64) bool {
	return SystemIdStart < uid && uid < SystemIdEnd
}

func IsTempId(uid int64) bool {
	return TempIdStart < uid && uid < TempIdEnd
}

func GenSysUid() int64 {
	r, err := db.Redis.Get(KeySystemIdIncr).Int64()
	if err != nil {
		return 0
	}
	next := r + 1
	db.Redis.Set(KeySystemIdIncr, next, 0)
	return next
}

func GenUid() int64 {
	r, err := db.Redis.Get(KeyUidIncr).Int64()
	if err != nil {
		return 0
	}
	next := r + rand.Int63n(4)
	db.Redis.Set(KeyUidIncr, next, 0)
	return next
}

func GenTemp() int64 {
	r, err := db.Redis.Get(KeyTempIdIncr).Int64()
	if err != nil {
		return 0
	}
	next := r + rand.Int63n(4)
	if next >= TempIdEnd {
		next = TempIdStart
	}
	db.Redis.Set(KeyTempIdIncr, next, 0)
	return next
}
