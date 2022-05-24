package message

import (
	"go_im/im/message/json"
	"go_im/protobuf/gen/pb_im"
)

// CustomerServiceMessage 表示客服消息
type CustomerServiceMessage struct {
	json.CustomerServiceMessage
}

// AckRequest 接收者回复给服务端确认收到消息
type AckRequest struct {
	pb_im.AckRequest
}

type AckGroupMessage struct {
	pb_im.AckGroupMessage
}

type Recall struct {
	pb_im.Recall
}
