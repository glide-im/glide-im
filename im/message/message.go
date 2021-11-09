package message

// GroupMessage 代表一个群消息
type GroupMessage struct {
	TargetId int64
	// internal
	Sender      int64 `json:"-"`
	Cid         int64
	UcId        int64
	MessageType int8
	Message     string
	SendAt      int64
}

// SenderChatMessage 表示服务端收到发送者的消息
type SenderChatMessage struct {
	Cid         int64
	UcId        int64
	Seq         int64
	TargetId    int64
	MessageType int8
	Message     string
	SendAt      int64
}

type ChatMessageAck struct {
	Seq int64
	Mid int64
}

type AckReceived struct {
	Mid    int64
	CMid   int64
	Sender int64
}

type SyncChatMessage struct {
}

// ReceiverChatMessage 表示服务端分发给接受者的聊天消息
type ReceiverChatMessage struct {
	Mid         int64
	Seq         int64
	AlignTag    string
	Cid         int64
	Sender      int64
	MessageType int8
	Message     string
	SendAt      int64
}

// CustomerServiceMessage 表示客服消息
type CustomerServiceMessage struct {
	// sender's id
	Sender int64
	// receiver's id
	Receiver int64
	// customer service id
	CsId int64

	ChatId      int64
	UserChatId  int64
	MessageType int8
	Message     string
	SendAt      int64
}
