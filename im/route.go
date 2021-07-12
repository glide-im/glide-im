package im

import "go_im/im/entity"

var MsgHandler = NewMsgHandler()

type MessageHandler struct {
	apiInterceptor *interceptorChain
	msgInterceptor *interceptorChain
}

func NewMsgHandler() *MessageHandler {
	mh := new(MessageHandler)
	mh.apiInterceptor = NewInterceptorChain()
	return mh
}

func (r *MessageHandler) Handle(client *Client, message *entity.Message) {

}
