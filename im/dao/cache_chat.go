package dao

import (
	"fmt"
	"github.com/go-redis/redis"
	"go_im/pkg/db"
	"go_im/pkg/logger"
	"strconv"
	"time"
)

const (
	// keyIncrMessageId 每个聊天会话的消息 ID 自增值
	keyIncrMessageId = "im:msg:incr:mid:"

	// keyChatUpdateAt 会话ID, 按更新时间排序的有序集合, 每次生成消息 ID 时, 将更新指定会话 ID 的 score 为当前时间
	keyChatUpdateAt = "im:msg:chat:update"

	keyChatIdIncr = "im:msg:incr:cid"

	keyUserChatIdIncr = "im:msg:incr:ucid"
)

func GetUserChatId(uid int64, chatID int64) (int64, error) {
	result, err := db.Redis.Incr(keyUserChatIdIncr).Result()
	if err == nil {
		return 0, err
	}
	return result, nil
}

func GetNextChatId(chatType int8) (int64, error) {
	result, err := db.Redis.Incr(keyChatIdIncr).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// GetMessageId 获取会话的下一个消息 ID, 这个消息 ID 是按 Chat 自增的
func GetMessageId(chatId int64) int64 {
	k := fmt.Sprintf("%s%d", keyIncrMessageId, chatId)
	result, err := db.Redis.Incr(k).Result()
	if err != nil || result == 0 {
		r := getMidFromDb(chatId)
		if r == 0 {
			logger.E("gen message id for chat error", chatId)
			return 0
		}
		db.Redis.Set(k, r, 0)
	}
	updateChat(chatId)
	return result
}

// getMidFromDb 从数据库中获取 mid, 一般情况是会话过期删除了
func getMidFromDb(cid int64) int64 {
	chat := ChatDao.GetChat(cid)
	if chat == nil {
		return 0
	}
	return chat.NextMid
}

// removeExpiredChat 移除过期的 chat, 通过有序集合 keyChatUpdateAt 中的时间为准, 从 0 到 now-secAgo 的会话信息
func removeExpiredChat(secAgo int64) {
	now := time.Now().Unix()
	expiredCid, err := db.Redis.ZRangeByScore(keyChatUpdateAt, redis.ZRangeBy{
		Min:    "0",
		Max:    strconv.FormatInt(now-secAgo, 10),
		Offset: 0,
	}).Result()
	if err != nil {
		logger.E("redis query chat update error", err)
		return
	}

	var expiredChatMidIncr []string
	var expiredCidUpdate []interface{}

	for _, cid := range expiredCid {
		keyChatIncr := fmt.Sprintf("%s%s", keyIncrMessageId, cid)
		expiredChatMidIncr = append(expiredChatMidIncr, keyChatIncr)
		expiredCidUpdate = append(expiredCidUpdate, cid)
	}

	// remove update
	_, er := db.Redis.ZRem(keyChatUpdateAt, expiredCidUpdate...).Result()
	if er != nil {
		logger.E("redis rm chat update error", er)
	}

	_, err = db.Redis.Del(expiredChatMidIncr...).Result()

	if err != nil {
		logger.E("redis remove chat update error", err)
	}
}

// updateChat 更新会话访问时间
func updateChat(chatId int64) {
	_, err := db.Redis.ZAdd(keyChatUpdateAt, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: chatId,
	}).Result()

	if err != nil {
		logger.E("redis update chat visit", err)
	}
}
