package messaging

import (
	"go_im/im/client"
	"go_im/im/dao/msgdao"
	"go_im/im/message"
	"go_im/pkg/logger"
	"strconv"
)

// dispatchChatMessage 分发用户单聊消息
func dispatchChatMessage(from int64, m *message.Message) {
	msg := new(message.UpChatMessage)
	if !unwrap(from, m, msg) {
		return
	}

	lg := from
	sm := msg.To
	if lg < sm {
		lg, sm = sm, lg
	}
	dbMsg := msgdao.ChatMessage{
		MID:        msg.Mid,
		From:       from,
		To:         msg.To,
		Type:       msg.Type,
		SendAt:     msg.CTime,
		Content:    msg.Content,
		CliSeq:     msg.CSeq,
		SessionTag: strconv.FormatInt(lg, 10) + "_" + strconv.FormatInt(sm, 10),
	}
	// 保存消息
	_, err := msgdao.AddChatMessage(&dbMsg)
	if err != nil {
		logger.E("save chat message error %v", err)
		return
	}

	// 对方不在线, 下发确认包
	if !client.Manager.IsOnline(from) {
		ackMsg := message.AckNotify{Mid: msg.Mid}
		ackNotify := message.NewMessage(0, message.ActionMessageAck, ackMsg)
		client.EnqueueMessage(from, ackNotify)
		err = msgdao.AddOfflineMessage(msg.To, msg.Mid)
		if err != nil {
			logger.E("save offline message error %v", err)
		}
		dispatchOffline(from, m)
	} else {
		dispatchOnline(from, msg)
	}
}

// dispatchOffline 接收者不在线, 离线推送
func dispatchOffline(from int64, message *message.Message) {

}

// dispatchOnline 接收者在线, 直接投递消息
func dispatchOnline(from int64, msg *message.UpChatMessage) {

	receiverMsg := message.DownChatMessage{
		Mid:     msg.Mid,
		CSeq:    msg.CSeq,
		From:    from,
		To:      msg.To,
		Content: msg.Content,
		CTime:   msg.CTime,
	}

	dispatchMsg := message.NewMessage(-1, message.ActionChatMessage, receiverMsg)
	client.EnqueueMessage(msg.To, dispatchMsg)
}
