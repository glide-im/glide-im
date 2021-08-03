package entity

import (
	"go_im/im/dao"
)

type GroupMessage struct {
	Gid         int64
	Uid         int64
	MessageType uint
	Content     string
}

// SenderChatMessage simple chat room message
type SenderChatMessage struct {
	Cid         int64
	UcId        int64
	Receiver    int64
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
