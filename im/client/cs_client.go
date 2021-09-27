package client

import (
	"go_im/im/comm"
	"go_im/im/message"
	"go_im/pkg/logger"
)

type CustomerServiceClient struct {
	// uid id of customer service
	uid int64

	messages chan *message.Message
	seq      *comm.AtomicInt64

	// waiter the online customer service staff
	waiter map[int64]interface{}
	// waiter to customer map
	w2c map[int64]int64
}

func NewCustomerServiceClient(uid int64) *CustomerServiceClient {
	client := new(CustomerServiceClient)
	client.messages = make(chan *message.Message, 200)
	client.uid = uid
	client.seq = new(comm.AtomicInt64)
	client.waiter = map[int64]interface{}{}
	return client
}

func (o *CustomerServiceClient) SignIn(_ int64, _ int64) {}

func (o *CustomerServiceClient) Closed() bool {
	// no connection, always online
	return false
}

func (o *CustomerServiceClient) EnqueueMessage(message *message.Message) {
	csMsg := new(CustomerServiceMessage)
	err := message.DeserializeData(csMsg)
	if err != nil {
		logger.E("customer service handle msg error", err)
		return
	}

	_, ok := o.waiter[csMsg.From]
	if ok {
		// customer service staffs' message, dispatch to customer
		EnqueueMessage(csMsg.To, message)
	} else {

	}
}

func (o *CustomerServiceClient) Exit() {

}

func (o *CustomerServiceClient) Run() {

}
