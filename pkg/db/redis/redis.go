package redis

import (
	"go_im/pkg/db"
	"time"
)

func Get(key string) string {
	return db.Redis.Get(key).Val()
}

func Set(key string, value string, expiration time.Duration) {
	db.Redis.Set(key, value, expiration)
}
