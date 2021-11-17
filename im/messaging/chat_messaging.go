package messaging

import (
	"go_im/im/client"
	"go_im/im/dao/msgdao"
	"go_im/im/message"
)

// dispatchChatMessage 分发用户单聊消息
func dispatchChatMessage(from int64, m *message.Message) {
	msg := new(message.UpChatMessage)
	if !unwrap(from, m, msg) {
		return
	}

	dbMsg := msgdao.ChatMessage{
		MID:        msg.Mid,
		ReceiveSeq: 0,
		CliSeq:     msg.CSeq,
		From:       from,
		To:         msg.To,
		Type:       msg.Type,
		SendAt:     msg.CTime,
		Content:    msg.Content,
	}

	// 保存到历史记录
	err := msgdao.AddChatMessage(&dbMsg)

	if err != nil {
		return
	}

	// 对方不在线, 下发确认包
	if !client.Manager.IsOnline(from) {
		ackMsg := message.AckNotify{Mid: msg.Mid}
		ackNotify := message.NewMessage(0, message.ActionMessageAck, ackMsg)
		client.EnqueueMessage(from, ackNotify)
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
		GSeq:    0,
		From:    from,
		To:      msg.To,
		Content: msg.Content,
		CTime:   msg.CTime,
	}

	dispatchMsg := message.NewMessage(-1, message.ActionChatMessage, receiverMsg)
	client.EnqueueMessage(msg.To, dispatchMsg)
}

func saveMsg2Db() {

}
