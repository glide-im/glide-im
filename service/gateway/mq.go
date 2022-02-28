package gateway

import (
	"github.com/nsqio/go-nsq"
	"go_im/protobuff/gen/pb_im"
	"go_im/protobuff/gen/pb_rpc"
	"go_im/service/cache"
	"go_im/service/mq_nsq"
	"google.golang.org/protobuf/proto"
)

var producer *mq_nsq.NSQProducer
var consumer *mq_nsq.NSQConsumer
var topicPrefix = "im_gateway_"

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
	bts, err := proto.Marshal(&m)
	if err != nil {
		return err
	}
	gateway, err := cache.GetGateway(uid)
	if err != nil {
		return err
	}
	err = producer.Publish(topicPrefix+gateway, bts)
	return err
}

type msgHandler struct {
}

func (msgHandler) HandleMessage(msg *nsq.Message) error {
	if len(msg.Body) == 0 {
		return nil
	}
	m := pb_rpc.NSQUserMessage{}
	err := proto.Unmarshal(msg.Body, &m)
	if err != nil {
		return err
	}
	msg.Finish()
	//client.EnqueueMessage(m.Uid, &message.Message{CommMessage: m.Message})
	return nil
}

func RegisterGateway(s *Server, config *mq_nsq.NSQConsumerConfig) error {
	var err error
	config.Channel = "ch"
	config.Topic = topicPrefix + s.Options.Addr
	consumer, err = mq_nsq.NewConsumer(config)
	if err != nil {
		return err
	}
	consumer.AddHandler(&msgHandler{})
	return consumer.Connect()
}
