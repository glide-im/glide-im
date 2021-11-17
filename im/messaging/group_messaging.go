package messaging

import (
	"go_im/im/group"
	"go_im/im/message"
)

// dispatchGroupMsg 分发群消息
func dispatchGroupMsg(from int64, msg *message.Message) {
	groupMsg := new(message.UpChatMessage)
	if !unwrap(from, msg, groupMsg) {
		return
	}
	groupMsg.From_ = from
	group.Manager.DispatchMessage(groupMsg.To, groupMsg)
}
