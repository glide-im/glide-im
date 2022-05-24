package groups

import "go_im/im/api/comm"

var (
	ErrGroupNotExit       = comm.NewApiBizError(3001, "ErrGroupNotExit")
	ErrMemberAlreadyExist = comm.NewApiBizError(3002, "ErrMemberAlreadyExist")
)
