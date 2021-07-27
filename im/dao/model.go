package dao

import "time"

type User struct {
	Uid      int64 `gorm:"primary_key"`
	Account  string
	Nickname string
	Password string
	Avatar   string

	CreateAt time.Time
	UpdateAt time.Time
}

type Chat struct {
	Cid          int64 `gorm:"primary_key"`
	ChatType     int8
	NewMessageAt time.Time
	CreateAt     time.Time
}

type UserChat struct {
	UcId         int64 `gorm:"primary_key"`
	Cid          int64
	Owner        int64
	Target       uint64
	ChatType     int8
	Unread       int
	NewMessageAt time.Time
	ReadAt       time.Time
	CreateAt     time.Time
}

type ChatMessage struct {
	Mid         int64 `gorm:"primary_key"`
	Cid         uint64
	SenderUid   int64
	SendAt      time.Time
	Message     string
	MessageType int8
	At          string
}
