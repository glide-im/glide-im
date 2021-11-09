package messaging

import (
	"github.com/panjf2000/ants/v2"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/dao"
	"go_im/im/group"
	"go_im/im/message"
	"go_im/im/statistics"
	"go_im/pkg/logger"
)

// execPool 100 capacity goroutine pool, 假设每个消息处理需要10ms, 一个协程则每秒能处理100条消息
var execPool *ants.Pool

func init() {
	client.MessageHandleFunc = messageHandler

	var err error
	execPool, err = ants.NewPool(200000,
		ants.WithNonblocking(true),
		ants.WithPanicHandler(onHandleMessagePanic),
		//ants.WithPreAlloc(true),
	)
	if err != nil {
		panic(err)
	}
}

// messageHandler 处理接收到的所有类型消息, 所有消息处理的入口
func messageHandler(from int64, device int64, msg *message.Message) {
	err := execPool.Submit(func() {
		statistics.SMsgInput()
		switch msg.Action {
		case message.ActionChatMessage:
			dispatchChatMessage(from, msg)
		case message.ActionGroupMessage:
			dispatchGroupMsg(from, msg)
		case message.ActionCSMessage:
			dispatchCustomerServiceMsg(from, msg)
		case message.ActionMessageAck:
			handleAckMsg(from, msg)
		default:
			if msg.Action.Contains(message.ActionApi) {
				api.Handle(from, msg)
			} else {
				client.EnqueueMessage(from, message.NewMessage(-1, message.ActionNotify, "unknown action"))
				logger.W("receive a unknown action message")
			}
		}
	})
	if err != nil {
		client.EnqueueMessage(from, message.NewMessage(-1, message.ActionNotify, "internal server error"))
		logger.E("async handle message error %v", err)
	}
}

// dispatchChatMessage 分发用户单聊消息
func dispatchChatMessage(from int64, msg *message.Message) {
	senderMsg := new(message.SenderChatMessage)
	if !unwrap(from, msg, senderMsg) {
		return
	}

	// 保存到历史记录
	chatMsg, err := dao.ChatDao.NewChatMessage(senderMsg.Cid, from, senderMsg.Message, senderMsg.MessageType)
	if err != nil {
		return
	}

	// 对方不在线, 下发确认包
	if !client.Manager.IsOnline(from) {
		ackMsg := message.AckReceived{
			Mid:    chatMsg.Mid,
			CMid:   0,
			Sender: from,
		}
		ackNotify := message.NewMessage(0, message.ActionMessageAck, ackMsg)
		client.EnqueueMessage(ackMsg.Sender, ackNotify)
		dispatchOffline(from, msg)
	} else {
		dispatchOnline(from, chatMsg, senderMsg)
	}
}

// dispatchOffline 接收者不在线, 离线推送
func dispatchOffline(from int64, message *message.Message) {

}

// dispatchOnline 接收者在线, 直接投递消息
func dispatchOnline(from int64, chatMsg *dao.ChatMessage, senderMsg *message.SenderChatMessage) {

	receiverMsg := message.ReceiverChatMessage{
		Mid:         chatMsg.Mid,
		Cid:         senderMsg.Cid,
		Sender:      from,
		MessageType: senderMsg.MessageType,
		Message:     senderMsg.Message,
		SendAt:      chatMsg.SendAt.Unix(),
	}

	dispatchMsg := message.NewMessage(-1, message.ActionChatMessage, receiverMsg)
	client.EnqueueMessage(senderMsg.TargetId, dispatchMsg)
}

// handleAckMsg 处理接收者收到消息发回来的确认消息
func handleAckMsg(from int64, msg *message.Message) {
	ackMsg := new(message.AckReceived)
	if !unwrap(from, msg, ackMsg) {
		return
	}
	ackNotify := message.NewMessage(0, message.ActionMessageAck, ackMsg)
	// 通知发送者, 对方已收到消息
	client.EnqueueMessage(ackMsg.Sender, ackNotify)
}

// dispatchGroupMsg 分发群消息
func dispatchGroupMsg(from int64, msg *message.Message) {
	groupMsg := new(message.GroupMessage)
	if !unwrap(from, msg, groupMsg) {
		return
	}
	groupMsg.Sender = from
	group.Manager.DispatchMessage(groupMsg.TargetId, groupMsg)
}

// dispatchCustomerServiceMsg 分发客服消息
func dispatchCustomerServiceMsg(from int64, msg *message.Message) {
	csMsg := new(message.CustomerServiceMessage)
	if !unwrap(from, msg, csMsg) {
		return
	}
	// 发送消息给客服
	client.EnqueueMessage(csMsg.CsId, msg)
}

func onHandleMessagePanic(i interface{}) {
	logger.E("handler message panic, %v", i)
}

// unwrap 解包, 反序列化消息包中数据到对象
func unwrap(from int64, msg *message.Message, to interface{}) bool {
	err := msg.DeserializeData(to)
	if err != nil {
		client.EnqueueMessage(from, message.NewMessage(msg.Seq, message.ActionNotify, "send message failed"))
		logger.E("sender chat senderMsg %v", err)
		return false
	}
	return true
}
