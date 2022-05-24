package dispatch

import (
	"github.com/glide-im/glideim/pkg/mq_nsq"
	"google.golang.org/protobuf/proto"
)

type nsqMsgProducer struct {
	producer *mq_nsq.NSQProducer
}

func newNsqMsgProducer(addr string) (*nsqMsgProducer, error) {
	c := &mq_nsq.NSQProducerConfig{
		Addr: addr,
	}
	producer, err := mq_nsq.NewProducer(c)
	if err != nil {
		return nil, err
	}
	p := &nsqMsgProducer{producer: producer}
	return p, nil
}

func (m *nsqMsgProducer) publish(node string, msg proto.Message) error {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	err = m.producer.Publish(node, bytes)
	return err
}
