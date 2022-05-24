package msg

import (
	"github.com/glide-im/glideim/pkg/db"
	"testing"
)

var groupMsgApi = GroupMsgApi{}

func init() {
	db.Init()
}

func TestGroupMsgApi_GetGroupMessageHistory(t *testing.T) {

	err := groupMsgApi.GetGroupMessageHistory(getContext(1, 1), &GroupMsgHistoryRequest{
		Gid: 4,
	})
	if err != nil {
		t.Error(err)
	}
}
