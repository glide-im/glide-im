package entity

type GroupMessage struct {
	Gid         uint64
	Uid         uint64
	MessageType uint
	Content     string
}

type ChatMessage struct {
	ChatId      uint64
	Target      int64
	MessageType int
	Message     string
}
