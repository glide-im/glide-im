package msgdao

import (
	"go_im/pkg/db"
	"strconv"
)

var Cache = cacheDao{}

const (
	keyUserMsgSeq = "im:msg:seq:"
)

type cacheDao struct{}

func (cacheDao) GetUserMsgSeq(uid int64) (int64, error) {
	k := keyUserMsgSeq + strconv.FormatInt(uid, 10)
	seq, err := db.Redis.Get(k).Result()
	if err != nil {
		return 0, err
	}
	seqI, err := strconv.ParseInt(seq, 10, 64)
	if err != nil {
		return 0, err
	}

	return seqI, nil
}

func (cacheDao) GetIncrUserMsgSeq(uid int64) (int64, error) {
	k := keyUserMsgSeq + strconv.FormatInt(uid, 10)
	seq, err := db.Redis.Incr(k).Result()
	if err != nil {
		return 0, err
	}
	return seq, nil
}
