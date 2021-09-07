package rpc

import (
	"go_im/config"
	"go_im/pkg/db"
	"strconv"
	"testing"
)

func TestRedisHSet(t *testing.T) {

	config.Init()
	db.Init()

	//db.Redis.Expire("im:user:host", time.Second*3)
	for i := 0; i < 10_0000; i++ {
		db.Redis.HSet("im:user:h", strconv.FormatInt(int64(i), 10), "127.0.0.1")
	}
}
