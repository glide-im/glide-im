package gateway

import (
	"go_im/im/client"
	"go_im/im/message"
	"go_im/pkg/logger"
	"go_im/pkg/mq_nsq"
	"go_im/protobuff/gen/pb_rpc"
	"google.golang.org/protobuf/proto"
	"sync"
	"testing"
	"time"
)

func TestProto(t *testing.T) {
	m := pb_rpc.NSQUserMessage{
		Uid:     1,
		Message: message.NewMessage(1, message.ActionChatMessage, nil).CommMessage,
	}
	bts, err := proto.Marshal(&m)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(bts))
	m2 := pb_rpc.NSQUserMessage{}
	err = proto.Unmarshal(bts, &m2)
	if err != nil {
		t.Error(err)
	}
}

func TestInitMQ(t *testing.T) {

	err := InitMessageProducer("127.0.0.1:4159")
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	for j := 0; j < 1; j++ {
		wg.Add(1)
		go func() {
			logger.D("start")
			for i := 0; i < 1; i++ {
				time.Sleep(time.Millisecond * 200)
				m := message.NewMessage(1, message.ActionChatMessage, nil)
				_ = client.EnqueueMessage(0, m)
				if err != nil {
					t.Error(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	logger.D("DONE")
	_ = producer.Stop()
}

func TestRegisterGateway(t *testing.T) {

	//config,_ := service.GetConfig()
	//newServer := NewServer(&rpc.ServerOptions{
	//	Name:        config.Gateway.Server.Name,
	//	Network:     config.Gateway.Server.Network,
	//	Addr:        config.Gateway.Server.Addr,
	//	Port:        config.Gateway.Server.Port,
	//	EtcdServers: []string{},
	//})
	err := RegisterGateway(nil, &mq_nsq.NSQConsumerConfig{
		NsqLookupds: []string{"127.0.0.1:4161"},
	})
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Minute * 1)
	err = consumer.Disconnect()
	if err != nil {
		t.Error(err)
		return
	}
}
