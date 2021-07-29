package dao

type User struct {
	Uid      int64 `gorm:"primary_key"`
	Account  string
	Nickname string
	Password string
	Avatar   string

	CreateAt Timestamp
	UpdateAt Timestamp
}

type Chat struct {
	Cid          int64 `gorm:"primary_key"`
	ChatType     int8
	NewMessageAt Timestamp
	CreateAt     Timestamp
}

type UserChat struct {
	UcId         int64 `gorm:"primary_key"`
	Cid          int64
	Owner        int64
	Target       uint64
	ChatType     int8
	Unread       int
	NewMessageAt Timestamp
	ReadAt       Timestamp
	CreateAt     Timestamp
}

type ChatMessage struct {
	Mid         int64 `gorm:"primary_key"`
	Cid         uint64
	SenderUid   int64
	SendAt      Timestamp
	Message     string
	MessageType int8
	At          string
}
