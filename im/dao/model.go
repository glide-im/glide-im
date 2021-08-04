package dao

type User struct {
	Uid      int64 `gorm:"primary_key"`
	Account  string
	Nickname string
	Password string
	Avatar   string

	CreateAt Timestamp `gorm:"type:datetime"`
	UpdateAt Timestamp `gorm:"type:datetime"`
}

type Friend struct {
	Fid     int64 `gorm:"primary_key"`
	Owner   int64
	Uid     int64
	Remark  string
	AddTime Timestamp `gorm:"type:datetime"`
}

type Chat struct {
	Cid          int64 `gorm:"primary_key"`
	ChatType     int8
	NewMessageAt Timestamp `gorm:"type:datetime"`
	CreateAt     Timestamp `gorm:"type:datetime"`
}

type UserChat struct {
	UcId         int64 `gorm:"primary_key"`
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
	Mid         int64 `gorm:"primary_key"`
	Cid         int64
	SenderUid   int64
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
	CreateAt Timestamp `gorm:"type:datetime"`
}

type GroupMember struct {
	Id     int64 `gorm:"primary_key"`
	Gid    int64
	Uid    int64
	Mute   int64
	Remark string
	Type   int8
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
