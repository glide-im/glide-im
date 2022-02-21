package message

import (
	"go_im/im/message/json"
	"go_im/im/message/pb/pb_msg"
)

// CustomerServiceMessage 表示客服消息
type CustomerServiceMessage struct {
	json.CustomerServiceMessage
}

// AckRequest 接收者回复给服务端确认收到消息
type AckRequest struct {
	pb_msg.AckRequest
}

type AckGroupMessage struct {
	pb_msg.AckGroupMessage
}

type Recall struct {
	pb_msg.Recall
}
