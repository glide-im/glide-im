package group

import (
	"go_im/pkg/db"
	"testing"
)

func TestLoadAllGroup(t *testing.T) {

	db.Init()
	groups := LoadAllGroup()
	for _, group := range groups {
		t.Log(group.gid, group.msgSequence, len(group.members))
	}
}
