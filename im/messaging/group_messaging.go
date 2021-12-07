package messaging

import (
	"go_im/im/client"
	"go_im/im/dao/msgdao"
	"go_im/im/group"
	"go_im/im/message"
)

// dispatchGroupMsg 分发群消息
func dispatchGroupMsg(from int64, msg *message.Message) {
	groupMsg := new(message.UpChatMessage)
	if !unwrap(from, msg, groupMsg) {
		return
	}
	groupMsg.From = from
	err := group.Manager.DispatchMessage(groupMsg.To, groupMsg)
	if err != nil {
		notify := message.NewMessage(0, message.ActionMessageFailed, message.AckNotify{Mid: groupMsg.Mid})
		client.EnqueueMessage(from, notify)
	}
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
