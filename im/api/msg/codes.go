package msg

import "go_im/im/api/comm"

var (
	errRecentMsgLoadFailed = comm.NewApiBizError(3001, "message load failed")
)
