package json

// ChatMessage 上行消息, 表示服务端收到发送者的消息
type ChatMessage struct {
	// Mid 消息ID
	Mid int64
	// Seq 发送者消息 seq
	Seq int64
	// From internal
	From int64
	// To 接收者 ID
	To int64
	// Type 消息类型
	Type int32
	// Content 消息内容
	Content string
	// SendAt 发送时间
	SendAt int64
}

// DownGroupMessage 下行群消息
type DownGroupMessage struct {
	ChatMessage
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

// AckRequest 接收者回复给服务端确认收到消息
type AckRequest struct {
	Seq  int64
	Mid  int64
	From int64
}

type AckGroupMessage struct {
	Gid int64
	Mid int64
	Seq int64
}

// AckMessage 服务端通知发送者的服务端收到消息
type AckMessage struct {
	Mid int64
	Seq int64
}

// AckNotify 服务端下发给发送者的消息送达通知
type AckNotify struct {
	Mid int64
}

type GroupNotify struct {
	Mid       int64
	Gid       int64
	Type      int64
	Seq       int64
	Timestamp int64
	Data      interface{}
}

type Recall struct {
	RecallBy int64
	Mid      int64
}
