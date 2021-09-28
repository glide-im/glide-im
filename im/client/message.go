package client

import (
	"go_im/im/dao"
)

// GroupMessage 代表一个群消息
type GroupMessage struct {
	TargetId    int64
	Cid         int64
	UcId        int64
	MessageType int8
	Message     string
	SendAt      dao.Timestamp
}

// SenderChatMessage 表示服务端收到发送者的消息
type SenderChatMessage struct {
	Cid         int64
	UcId        int64
	TargetId    int64
	MessageType int8
	Message     string
	SendAt      dao.Timestamp
}

// ReceiverChatMessage 表示服务端分发给接受者的聊天消息
type ReceiverChatMessage struct {
	Mid         int64
	Cid         int64
	UcId        int64
	Sender      int64
	MessageType int8
	Message     string
	SendAt      dao.Timestamp
}

// CustomerServiceMessage 表示客服消息
type CustomerServiceMessage struct {
	// sender's id
	From int64
	// receiver's id
	To int64
	// customer service id
	CsId int64

	ChatId      int64
	UserChatId  int64
	MessageType int8
	Message     string
	SendAt      dao.Timestamp
}
