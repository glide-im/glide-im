package mq_nsq

import (
	"testing"
)

func TestInitConsumer(t *testing.T) {
	c, err := NewConsumer(&NSQConsumerConfig{
		Topic:       "t",
		Channel:     "",
		NsqLookupds: []string{"127.0.0.1:4161"},
	})
	if err != nil {
		t.Error(err)
	}

	err = c.Connect()
	if err != nil {
		t.Error(err)
	}
}
