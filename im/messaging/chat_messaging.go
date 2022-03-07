package messaging

import (
	"go_im/im/client"
	"go_im/im/dao/msgdao"
	"go_im/im/dao/uid"
	"go_im/im/message"
	"go_im/pkg/logger"
	"strconv"
)

// handleChatMessage 分发用户单聊消息
func handleChatMessage(from int64, device int64, m *message.Message) {
	if uid.IsTempId(from) {
		logger.D("not sign in")
		client.EnqueueMessage(from, message.NewMessage(0, message.ActionNotifyNeedAuth, ""))
		return
	}
	msg := new(message.ChatMessage)
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
		if m.Action == message.ActionChatMessageRecall {
			r := &message.Recall{}
			err := message.UnmarshallJson(msg.Content, r)
			if err != nil || r.RecallBy != from {
				return
			}
			err = msgdao.ChatMsgDaoImpl.UpdateChatMessageStatus(r.Mid, r.RecallBy, msg.To, msgdao.ChatMessageStatusRecalled)
			if err != nil {
				logger.E("update message status error %v", err)
				return
			}
		} else {
			dbMsg := msgdao.ChatMessage{
				MID:       msg.Mid,
				From:      from,
				To:        msg.To,
				Type:      msg.Type,
				SendAt:    msg.SendAt,
				Content:   msg.Content,
				CliSeq:    msg.Seq,
				SessionID: sessionId,
			}
			// 保存消息
			_, err := msgdao.AddChatMessage(&dbMsg)
			if err != nil {
				logger.E("save chat message error %v", err)
				return
			}
		}
	}

	// 告诉客户端服务端已收到
	ackChatMessage(from, device, msg.Mid)

	// 对方不在线, 下发确认包
	// TODO 2022-1-17 处理假在线, 假链接
	if !client.IsOnline(msg.To) {
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

func handleChatRecallMessage(from int64, device int64, msg *message.Message) {
	handleChatMessage(from, device, msg)
}

func ackNotifyMessage(from int64, mid int64) {
	ackNotify := message.NewAckNotify(mid)
	msg := message.NewMessage(0, message.ActionAckNotify, &ackNotify)
	client.EnqueueMessage(from, msg)
}

func ackChatMessage(from int64, device int64, mid int64) {
	ackMsg := message.NewAckMessage(mid, 0)
	ack := message.NewMessage(0, message.ActionAckMessage, &ackMsg)
	client.EnqueueMessageToDevice(from, device, ack)
}

// dispatchOffline 接收者不在线, 离线推送
func dispatchOffline(from int64, message *message.Message) {

}

// dispatchOnline 接收者在线, 直接投递消息
func dispatchOnline(from int64, msg *message.ChatMessage) {

	receiverMsg := msg
	msg.From = from
	dispatchMsg := message.NewMessage(-1, message.ActionChatMessage, receiverMsg)
	client.EnqueueMessage(msg.To, dispatchMsg)
}
