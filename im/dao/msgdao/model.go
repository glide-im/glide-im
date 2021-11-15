package msgdao

// ChatMessageModel 一对一聊天全量消息
type ChatMessageModel struct {
	MID        int64 `gorm:"primary_key"`
	cliMsgID   int64
	ReceiveSeq int64
	From       int64
	To         int64
	Type       int64
	SendAt     int64
	Content    string
}

// OfflineMessageModel 用户不在线, 离线消息
type OfflineMessageModel struct {
	ID  int64 `gorm:"primary_key"`
	MID int64
	UID int64
}

// GroupMessageModel 全量群消息
type GroupMessageModel struct {
	MID      int64 `gorm:"primary_key"`
	cliMsgID int64
	Seq      int64
	To       int64
	From     int64
	Type     int64
	SendAt   int64
	Content  string
}

// GroupMemberMsgStateModel 群成员确认收到消息记录, 用于计算离线消息的同步量
type GroupMemberMsgStateModel struct {
	MbID       string `gorm:"primary_key"`
	GID        int64
	UID        int64
	LastAckMID int64
	LastAckSeq int64
}

// GroupMessageStateModel 群消息最新状态 ID 及 seq
type GroupMessageStateModel struct {
	GID       int64 `gorm:"primary_key"`
	LastMID   int64
	LastSeq   int64
	LastMsgAt int64
}

// GroupMsgSeqModel 群消息 seq 状态
type GroupMsgSeqModel struct {
	GID  int64 `gorm:"primary_key"`
	Seq  int64
	Step int64
}
