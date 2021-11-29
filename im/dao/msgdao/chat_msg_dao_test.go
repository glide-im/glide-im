package msgdao

import (
	"testing"
	"time"
)

func TestChatMsgDao_GetRecentChatMessages(t *testing.T) {
	messages, err := instance.GetRecentChatMessages(1, 1637650186)
	if err != nil {
		t.Error(err)
	}
	for _, message := range messages {
		t.Log(message)
	}
}

func TestChatMsgDao_GetOfflineMessage(t *testing.T) {
	m, err := GetOfflineMessage(1)
	if err != nil {
		t.Error(err)
	}
	t.Log(m)
}

func TestChatMsgDao_AddOfflineMessage(t *testing.T) {
	err := AddOfflineMessage(1, 4)
	if err != nil {
		t.Error(err)
	}
}

func TestChatMsgDao_DelOfflineMessage(t *testing.T) {
	err := DelOfflineMessage(1, []int64{1, 2, 3, 4})
	if err != nil {
		t.Error(err)
	}
}

func TestChatMsgDao_AddOrUpdateChatMessage(t *testing.T) {

	message, err := AddChatMessage(&ChatMessage{
		MID:        14,
		SessionTag: "2_1",
		CliSeq:     2,
		From:       2,
		To:         1,
		Type:       1,
		SendAt:     time.Now().Unix(),
		CreateAt:   time.Now().Unix(),
		Content:    "hello",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func TestChatMsgDao_GetChatMessageSeqAfter(t *testing.T) {
	after, err := GetChatMessageMidAfter(1, 2, 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(after)
}
