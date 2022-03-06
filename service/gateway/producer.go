package gateway

import (
	"go_im/im/message"
	"go_im/service/cache"
	"go_im/service/mq_nsq"
)

var producer *mq_nsq.NSQProducer
var consumer *mq_nsq.NSQConsumer
var topicPrefix = "im_gateway_"

// InitMessageProducer  init service as a gateway message producer, nsqdAddr is the address of local nsqd
func InitMessageProducer(nsqdAddr string) error {
	var err error
	c := &mq_nsq.NSQProducerConfig{
		Addr: nsqdAddr,
	}
	producer, err = mq_nsq.NewProducer(c)
	return err
}

type gateway struct {
}

func (g gateway) ClientSignIn(oldUid int64, uid int64, device int64) error {
	topic, err := cache.GetGateway(oldUid, 0)
	if err != nil {
		return err
	}
	return producer.Publish("im_signin_"+topic, nil)
}

func (g gateway) ClientLogout(uid int64, device int64) error {
	topic, err := cache.GetGateway(uid, device)
	if err != nil {
		return err
	}
	return producer.Publish("im_logout_"+topic, nil)
}

func (g gateway) EnqueueMessage(uid int64, device int64, message *message.Message) error {
	topic, err := cache.GetGateway(uid, device)
	if err != nil {
		return err
	}
	return producer.Publish("im_msg_"+topic, nil)
}
