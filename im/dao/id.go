package dao

import (
	"fmt"
	"go_im/pkg/db"
	"go_im/pkg/logger"
)

var keyMessageId = "im:mid:incr"

func GenMessageId(chatId int64) int64 {
	k := fmt.Sprintf("%s:%d", keyMessageId, chatId)
	result, err := db.Redis.Incr(k).Result()
	if err != nil || result == 0 {
		r := getMidFromDb(chatId)
		if r == 0 {
			logger.E("gen message id for chat error", chatId)
			return 0
		}
		db.Redis.Set(k, r, 0)
	}
	return result
}

func getMidFromDb(cid int64) int64 {
	chat := ChatDao.GetChat(cid)
	if chat == nil {
		return 0
	}
	return chat.NextMid
}
