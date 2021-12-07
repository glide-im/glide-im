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
	msg.From = from

	if m.Action != message.ActionChatMessageResend {
		lg := from
		sm := msg.To
		if lg < sm {
			lg, sm = sm, lg
		}
		sessionId := strconv.FormatInt(lg, 10) + "_" + strconv.FormatInt(sm, 10)
		dbMsg := msgdao.ChatMessage{
			MID:       msg.Mid,
			From:      from,
			To:        msg.To,
			Type:      msg.Type,
			SendAt:    msg.CTime,
			Content:   msg.Content,
			CliSeq:    msg.CSeq,
			SessionID: sessionId,
		}
		// 保存消息
		_, err := msgdao.AddChatMessage(&dbMsg)
		if err != nil {
			logger.E("save chat message error %v", err)
			return
		}
	}

	// 告诉客户端服务端已收到
	ackChatMessage(from, msg.Mid)

	// 对方不在线, 下发确认包
	if !client.Manager.IsOnline(msg.To) {
		ackNotifyMessage(from, msg.Mid)
		err := msgdao.AddOfflineMessage(msg.To, msg.Mid)
		if err != nil {
			logger.E("save offline message error %v", err)
		}
		dispatchOffline(from, m)
	} else {
		dispatchOnline(from, msg)
	}
}

func ackNotifyMessage(from int64, mid int64) {
	ackNotify := message.AckNotify{Mid: mid}
	msg := message.NewMessage(0, message.ActionAckNotify, ackNotify)
	client.EnqueueMessage(from, msg)
}

func ackChatMessage(from int64, mid int64) {
	ackMsg := message.AckMessage{Mid: mid}
	ack := message.NewMessage(0, message.ActionAckMessage, ackMsg)
	client.EnqueueMessage(from, ack)
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
		Type:    msg.Type,
		Content: msg.Content,
		CTime:   msg.CTime,
	}

	dispatchMsg := message.NewMessage(-1, message.ActionChatMessage, receiverMsg)
	client.EnqueueMessage(msg.To, dispatchMsg)
}
