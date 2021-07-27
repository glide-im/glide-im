package entity

import "time"

type GroupMessage struct {
	Gid         uint64
	Uid         uint64
	MessageType uint
	Content     string
}

// SenderChatMessage simple chat room message
type SenderChatMessage struct {
	Cid         uint64
	UcId        int64
	Receiver    int64
	MessageType int8
	Message     string
	SendAt      time.Time
}

type ReceiverChatMessage struct {
	Mid         int64
	Cid         uint64
	UcId        int64
	Sender      int64
	MessageType int8
	Message     string
	SendAt      time.Time
}
