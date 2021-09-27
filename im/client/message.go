package client

import (
	"go_im/im/dao"
)

type GroupMessage struct {
	TargetId    int64
	Cid         int64
	UcId        int64
	MessageType int8
	Message     string
	SendAt      dao.Timestamp
}

// SenderChatMessage simple chat room message
type SenderChatMessage struct {
	Cid         int64
	UcId        int64
	TargetId    int64
	MessageType int8
	Message     string
	SendAt      dao.Timestamp
}

type ReceiverChatMessage struct {
	Mid         int64
	Cid         int64
	UcId        int64
	Sender      int64
	MessageType int8
	Message     string
	SendAt      dao.Timestamp
}

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
