package messaging

import (
	"github.com/panjf2000/ants/v2"
	"go_im/im/api"
	"go_im/im/client"
	"go_im/im/message"
	"go_im/im/statistics"
	"go_im/pkg/logger"
)

// execPool 100 capacity goroutine pool, 假设每个消息处理需要10ms, 一个协程则每秒能处理100条消息
var execPool *ants.Pool

var messageHandlerFunMap = map[message.Action]func(from int64, msg *message.Message){
	message.ActionGroupMessageRecall: dispatchGroupRecallMsg,
	message.ActionChatMessageRecall:  dispatchChatRecallMessage,
	message.ActionChatMessage:        dispatchChatMessage,
	message.ActionChatMessageRetry:   dispatchChatMessage,
	message.ActionChatMessageResend:  dispatchChatMessage,
	message.ActionGroupMessage:       dispatchGroupMsg,
	message.ActionCSMessage:          dispatchCustomerServiceMsg,
	message.ActionAckRequest:         handleAckRequest,
	message.ActionAckGroupMsg:        handleAckGroupMsgRequest,
}

func Init() {
	client.MessageHandleFunc = messageHandler

	var err error
	execPool, err = ants.NewPool(100_0000,
		ants.WithNonblocking(true),
		ants.WithPanicHandler(onHandleMessagePanic),
		ants.WithPreAlloc(true),
	)
	if err != nil {
		panic(err)
	}
}

// messageHandler 处理接收到的所有类型消息, 所有消息处理的入口
func messageHandler(from int64, device int64, msg *message.Message) {
	logger.D("new message: uid=%d, %v", from, msg)
	err := execPool.Submit(func() {
		statistics.SMsgInput()
		h, ok := messageHandlerFunMap[msg.Action]
		if ok {
			h(from, msg)
			return
		}
		switch msg.Action {
		case message.ActionHeartbeat:
			handleHeartbeat(from, device, msg)
		default:
			if msg.Action.Contains(message.ActionApi) {
				api.Handle(from, device, msg)
			} else {
				client.EnqueueMessage(from, message.NewMessage(-1, message.ActionNotifyError, "unknown action"))
				logger.W("receive a unknown action message: " + string(msg.Action))
			}
		}
	})
	if err != nil {
		if err == ants.ErrPoolOverload {
			logger.E("Messaging.MessageHandler goroutine pool is overload")
			return
		}
		if err == ants.ErrPoolClosed {
			logger.E("Messaging.MessageHandler goroutine pool is closed")
			return
		}
		client.EnqueueMessage(from, message.NewMessage(-1, message.ActionNotifyError, "internal server error"))
		logger.E("async handle message error %v", err)
	}
}

func handleHeartbeat(from int64, device int64, msg *message.Message) {
	// TODO 2021-11-15 处理心跳消息
}

// handleAckRequest 处理接收者收到消息发回来的确认消息
func handleAckRequest(from int64, msg *message.Message) {
	ackMsg := new(message.AckRequest)
	if !unwrap(from, msg, ackMsg) {
		return
	}
	ackNotify := message.NewMessage(0, message.ActionAckNotify, ackMsg)
	// 通知发送者, 对方已收到消息
	client.EnqueueMessage(ackMsg.From, ackNotify)
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
	statistics.SError(i.(error))
	logger.E("handler message panic, %v", i)
}

// unwrap 解包, 反序列化消息包中数据到对象
func unwrap(from int64, msg *message.Message, to interface{}) bool {
	err := msg.DeserializeData(to)
	if err != nil {
		client.EnqueueMessage(from, message.NewMessage(msg.Seq, message.ActionNotifyError, "send message failed"))
		logger.E("sender chat senderMsg %v", err)
		return false
	}
	return true
}
