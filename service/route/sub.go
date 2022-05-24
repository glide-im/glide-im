package route

import (
	"encoding/json"
	"github.com/glide-im/glideim/pkg/mq_nsq"
	"github.com/glide-im/glideim/service"
	"github.com/nsqio/go-nsq"
)

var producer *mq_nsq.NSQProducer
var consumer *mq_nsq.NSQConsumer

const (
	groupRouteTopic = "group_rt_update"
)

type GroupRouteInfo struct {
	Gid  int64
	Addr string
	Port int
}

type groupRouteHandler struct {
}

type ProducerConf struct {
	Nsqd string
}

type ConsumerConf struct {
	Channel     string
	NsqLookupds []string
}

func (g *groupRouteHandler) HandleMessage(message *nsq.Message) error {
	info := GroupRouteInfo{}
	err := json.Unmarshal(message.Body, &info)
	if err != nil {
		return err
	}
	err = setGroup(info.Gid, info)
	return err
}

// Init the nsq consumer, producer, pass nil will not init.
// group service publish group route info.
func Init(pConf *ProducerConf, cConf *ConsumerConf) error {
	if cConf != nil {
		config := &mq_nsq.NSQConsumerConfig{
			Topic:       groupRouteTopic,
			Channel:     cConf.Channel,
			NsqLookupds: cConf.NsqLookupds,
		}
		var err error
		consumer, err = mq_nsq.NewConsumer(config)
		if err != nil {
			return err
		}
		consumer.AddHandler(&groupRouteHandler{})
		err = consumer.Connect()
		if err != nil {
			return err
		}
	}
	if pConf != nil {
		var err error = nil
		config := &mq_nsq.NSQProducerConfig{Addr: pConf.Nsqd}
		producer, err = mq_nsq.NewProducer(config)
		if err != nil {
			return err
		}
	}
	return nil
}

// PublishGroupRoute notify the service update group route info who dependence group service,
//calling when group service load or init group
func PublishGroupRoute(gid int64) error {
	conf, err := service.GetConfig()
	if err != nil {
		return err
	}
	info := GroupRouteInfo{
		Gid:  gid,
		Addr: conf.GroupMessaging.Server.Addr,
		Port: conf.GroupMessaging.Server.Port,
	}
	j, err := json.Marshal(&info)
	if err != nil {
		return err
	}
	return producer.Publish(groupRouteTopic, j)
}
