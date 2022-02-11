package messaging

import (
	"go_im/im/client"
	"go_im/im/dao/msgdao"
	"go_im/im/dao/uid"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/pkg/logger"
)

// dispatchGroupMsg 分发群消息
func dispatchGroupMsg(from int64, msg *message.Message) {
	if uid.IsTempId(from) {
		logger.D("not sign in, uid=%d", from)
		return
	}
	groupMsg := new(message.UpChatMessage)
	if !unwrap(from, msg, groupMsg) {
		return
	}
	groupMsg.From = from

	var err error
	if msg.Action == message.ActionGroupMessageRecall {
		err = group.DispatchMessage(groupMsg.To, groupMsg)
	} else {
		err = group.DispatchRecallMessage(groupMsg.To, groupMsg)
	}
	if err != nil {
		logger.E("dispatch group message error: %v", err)
		notify := message.NewMessage(0, message.ActionMessageFailed, message.AckNotify{Mid: groupMsg.Mid})
		client.EnqueueMessage(from, notify)
	}
}

func dispatchGroupRecallMsg(from int64, msg *message.Message) {
	dispatchGroupMsg(from, msg)
}

func handleAckGroupMsgRequest(from int64, msg *message.Message) {
	ack := new(message.AckGroupMessage)
	if !unwrap(from, msg, ack) {
		return
	}
	err := msgdao.UpdateGroupMemberMsgState(ack.Gid, from, ack.Mid, ack.Seq)
	if err != nil {

	}
}
