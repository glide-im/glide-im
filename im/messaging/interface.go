package messaging

import (
	"go_im/im/client"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/pkg/logger"
)

func HandleMessage(from int64, device int64, msg *message.Message) {
	handleMessage(from, device, msg)
}

func dispatchGroupMessage(gid int64, msg *message.ChatMessage) error {
	return group.DispatchMessage(gid, msg)
}

func dispatchRecallMessage(gid int64, msg *message.ChatMessage) error {
	return group.DispatchRecallMessage(gid, msg)
}

func enqueueMessage(uid int64, message *message.Message) {
	err := client.EnqueueMessage(uid, message)
	if err != nil {
		logger.E("%v", err)
	}
}
