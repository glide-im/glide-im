package dao

type User struct {
	Uid      int64  `gorm:"primary_key"`
	Account  string `gorm:"unique"`
	Nickname string
	Password string
	Avatar   string

	CreateAt Timestamp `gorm:"type:datetime"`
	UpdateAt Timestamp `gorm:"type:datetime"`
}

type Contacts struct {
	Fid      int64 `gorm:"primary_key"`
	Owner    int64
	TargetId int64
	Remark   string
	Type     int8
	AddTime  Timestamp `gorm:"type:datetime"`
}

type ChatMessageID struct {
	Cid        int64 `gorm:"primary_key"`
	CurrentMid int64
}

type Chat struct {
	Cid          int64 `gorm:"primary_key"`
	ChatType     int8
	TargetId     int64
	CurrentMid   int64
	NewMessageAt Timestamp `gorm:"type:datetime"`
	CreateAt     Timestamp `gorm:"type:datetime"`
}

type UserChat struct {
	UcId int64 `gorm:"primary_key"`
	// IDs 双方ID拼接字符串,小的ID在前
	IDs          string
	Cid          int64
	Owner        int64
	Target       int64
	ChatType     int8
	Unread       int
	NewMessageAt Timestamp `gorm:"type:datetime"`
	ReadAt       Timestamp `gorm:"type:datetime"`
	CreateAt     Timestamp `gorm:"type:datetime"`
}

type ChatMessage struct {
	ID          int64 `gorm:"primary_key"`
	Mid         int64
	Cid         int64
	Sender      int64
	SendAt      Timestamp `gorm:"type:datetime"`
	Message     string
	MessageType int8
	At          string
}

type Group struct {
	Gid      int64 `gorm:"primary_key"`
	Name     string
	Avatar   string
	Owner    int64
	Mute     bool
	Notice   string
	ChatId   int64
	CreateAt Timestamp `gorm:"type:datetime"`
}

type GroupMember struct {
	Id     int64 `gorm:"primary_key"`
	Gid    int64
	Uid    int64
	Mute   int64
	Remark string
	Flag   int32
	JoinAt Timestamp `gorm:"type:datetime"`
}

type GroupMessage struct {
	GmId        int64 `gorm:"primary_key"`
	Cid         int64
	SenderUid   int64
	SendAt      Timestamp `gorm:"type:datetime"`
	Message     string
	MessageType int8
	At          string
}

type GroupNotify struct {
	Id     int64 `gorm:"primary_key"`
	Gid    int64
	Uid    int64
	Remark string
	Type   int8
	State  int8
}
