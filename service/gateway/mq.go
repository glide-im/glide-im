package gateway

import (
	"github.com/gogo/protobuf/proto"
	"github.com/nsqio/go-nsq"
	"go_im/im/client"
	"go_im/protobuff/pb_im"
	"go_im/protobuff/pb_rpc"
	"go_im/service/mq_nsq"
)

var producer *mq_nsq.NSQProducer

func InitMQ(addr string) error {
	var err error
	c := &mq_nsq.NSQProducerConfig{
		Addr: addr,
	}
	producer, err = mq_nsq.NewProducer(c)
	return err
}

func PublishMsg(uid int64, message *pb_im.CommMessage) error {
	m := pb_rpc.NSQUserMessage{
		Uid:     uid,
		Message: message,
	}
	b, err := proto.Marshal(&m)
	if err != nil {
		return err
	}
	err = producer.Publish("topic", b)
	return err
}

type msgHandler struct {
}

func (msgHandler) HandleMessage(msg *nsq.Message) error {
	m := pb_rpc.NSQUserMessage{}
	err := proto.Unmarshal(msg.Body, &m)
	if err != nil {
		return err
	}
	client.EnqueueMessage(m.Uid, m.Message)
	return nil
}

func RegisterGateway(s Server, config *mq_nsq.NSQConsumerConfig) error {
	c, err := mq_nsq.NewConsumer(config)
	if err != nil {
		return err
	}
	c.AddHandler(&msgHandler{})
	return err
}
