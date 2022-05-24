package msg

import "github.com/glide-im/glideim/im/api/comm"

var (
	errRecentMsgLoadFailed = comm.NewApiBizError(3001, "message load failed")
)
