package message

// UpChatMessage 上行消息, 表示服务端收到发送者的消息
type UpChatMessage struct {
	// Mid 消息ID
	Mid int64
	// CSeq 发送者消息 seq
	CSeq int64
	// From internal
	From int64
	// To 接收者 ID
	To int64
	// Type 消息类型
	Type int
	// Content 消息内容
	Content string
	// CTime 发送时间
	CTime int64
}

// DownChatMessage 表示服务端分发给接受者的聊天消息
type DownChatMessage struct {
	Mid     int64
	CSeq    int64
	From    int64
	To      int64
	Type    int
	Content string
	CTime   int64
}

// DownGroupMessage 下行群消息
type DownGroupMessage struct {
	Mid int64
	// Seq 群消息 Seq
	Seq     int64
	Gid     int64
	Type    int
	From    int64
	Content string
	SendAt  int64
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
	Gid  int64
	Type int64
}
