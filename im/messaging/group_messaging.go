package messaging

import (
	"go_im/im/dao/msgdao"
	"go_im/im/dao/uid"
	"go_im/im/message"
	"go_im/pkg/logger"
)

// handleGroupMsg 分发群消息
func handleGroupMsg(from int64, device int64, msg *message.Message) {
	if uid.IsTempId(from) {
		logger.D("not sign in, uid=%d", from)
		return
	}
	groupMsg := new(message.ChatMessage)
	if !unwrap(from, msg, groupMsg) {
		return
	}
	groupMsg.From = from

	var err error
	if msg.Action == message.ActionGroupMessageRecall {
		err = dispatchRecallMessage(groupMsg.To, groupMsg)
	} else {
		err = dispatchGroupMessage(groupMsg.To, groupMsg)
	}
	if err != nil {
		logger.E("dispatch group message error: %v", err)
		notify := message.NewMessage(0, message.ActionMessageFailed, message.NewAckNotify(groupMsg.Mid))
		enqueueMessage(from, notify)
	}
}

func handleGroupRecallMsg(from int64, device int64, msg *message.Message) {
	handleGroupMsg(from, device, msg)
}

func handleAckGroupMsgRequest(from int64, device int64, msg *message.Message) {
	ack := new(message.AckGroupMessage)
	if !unwrap(from, msg, ack) {
		return
	}
	err := msgdao.UpdateGroupMemberMsgState(ack.Gid, from, ack.Mid, ack.Seq)
	if err != nil {

	}
}
