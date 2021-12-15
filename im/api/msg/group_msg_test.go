package msg

import (
	"go_im/pkg/db"
	"testing"
)

var groupMsgApi = GroupMsgApi{}

func init() {
	db.Init()
}

func TestGroupMsgApi_GetGroupMessageHistory(t *testing.T) {

	err := groupMsgApi.GetGroupMessageHistory(getContext(1, 1), &GroupMsgHistoryRequest{
		Gid:  4,
		Page: 1,
	})
	if err != nil {
		t.Error(err)
	}
}
