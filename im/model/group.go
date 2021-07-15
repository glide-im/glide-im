package model

/**
Member

Type 1: 群员 2: 管理 3: 群主
State 状态位 0000 : 0-0-通知开关-被禁言
*/
type Member struct {
	Uid      int64
	Nickname string
	Avatar   string
	Type     uint8
	State    uint8
}
