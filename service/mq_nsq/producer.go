package mq_nsq

import "github.com/nsqio/go-nsq"

type NSQProducerConfig struct {
	Addr string
}

type NSQProducer struct {
	producer *nsq.Producer
}

func NewProducer(c *NSQProducerConfig) (*NSQProducer, error) {

	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(c.Addr, config)
	if err != nil {
		return nil, err
	}
	p := NSQProducer{
		producer: producer,
	}
	return &p, nil
}

func (p *NSQProducer) Publish(topic string, m []byte) error {
	return p.producer.Publish(topic, m)
}
