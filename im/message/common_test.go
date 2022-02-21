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
