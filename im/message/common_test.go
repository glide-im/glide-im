package message

import (
	"testing"
)

func TestNewMessage(t *testing.T) {
	c := NewChatMessage(1, 1, 1, 1, 1, "", 1)
	message := NewMessage(1, "", &c)
	encode, err := DefaultCodec.Encode(message)
	if err != nil {
		t.Error(err)
	}
	message = NewEmptyMessage()
	err = DefaultCodec.Decode(encode, message)
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
}

func TestJsonCodec_Decode(t *testing.T) {
	c := NewChatMessage(1, 1, 1, 1, 1, "", 1)
	m := NewMessage(1, ActionGroupMessageRecall, &c)
	j, err := JsonCodec.Encode(&m)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(j))
}
