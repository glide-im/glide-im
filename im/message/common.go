package message

import (
	"errors"
	"github.com/glide-im/glideim/im/message/json"
	"github.com/glide-im/glideim/im/message/pb"
	"github.com/glide-im/glideim/pkg/logger"
	"github.com/glide-im/glideim/protobuf/gen/pb_im"
	"github.com/glide-im/glideim/protobuf/gen/pb_rpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Message struct {
	pb   *pb_im.CommMessage
	json *json.ComMessage
	data interface{}
}

func (m *Message) SetSeq(seq int64) {
	if m.json != nil {
		m.json.Seq = seq
	}
	if m.pb != nil {
		m.pb.Seq = seq
	}
}

func (m *Message) GetSeq() int64 {
	if m.json != nil {
		return m.json.Seq
	}
	if m.pb != nil {
		return m.pb.Seq
	}
	return 0
}

func (m *Message) GetAction() string {
	if m.json != nil {
		return m.json.Action
	}
	if m.pb != nil {
		return m.pb.Action
	}
	return ""
}

func (m *Message) GetData() interface{} {
	if m.data != nil {
		return m.data
	}
	return nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	if m.json == nil {
		if m.pb == nil {
			m.json = json.NewMessage(0, "", nil)
		} else {
			//if m.data == nil {
			//	return nil, errors.New("cannot marshal protobuf msg to json")
			//}
			if m.data == nil {
				m.data = m.pb.Data
			}
			m.json = json.NewMessage(m.pb.Seq, m.pb.Action, m.data)
		}
	}
	return JsonCodec.Encode(m.json)
}

func (m *Message) UnmarshalJSON(bytes []byte) error {
	if m.json == nil {
		m.json = &json.ComMessage{}
	}
	return JsonCodec.Decode(bytes, m.json)
}

func (m *Message) ProtoReflect() protoreflect.Message {
	return m.GetProtobuf().ProtoReflect()
}

func (m *Message) GetProtobuf() *pb_im.CommMessage {
	if m.pb == nil {
		if m.json == nil {
			m.pb = &pb_im.CommMessage{}
		} else {
			data := m.json.Data
			_, ok := data.Data().(proto.Message)
			if !ok {
				jb, err := data.MarshalJSON()
				if err != nil {
					logger.E("%v", err)
					return nil
				}
				m.data = &pb_rpc.JsonString{Json: string(jb)}
			} else {
				m.data = data
			}
			m.pb = pb.NewMessage(m.json.Seq, m.json.Action, m.data)
		}
	}
	return m.pb
}

func FromProtobuf(message *pb_im.CommMessage) *Message {
	return &Message{
		pb:   message,
		json: nil,
		data: nil,
	}
}

func NewMessage(seq int64, action Action, data interface{}) *Message {

	message := Message{
		pb:   nil,
		json: nil,
		data: data,
	}
	_, ok := data.(proto.Message)
	if ok {
		message.pb = pb.NewMessage(seq, string(action), data)
	} else {
		message.json = json.NewMessage(seq, string(action), data)
	}

	return &message
}

func NewEmptyMessage() *Message {
	return &Message{
		pb:   nil,
		json: nil,
		data: nil,
	}
}

func (m *Message) DeserializeData(v interface{}) error {
	if m.pb != nil {
		return ProtoBuffCodec.Decode(m.pb.Data.Value, v)
	}
	if m.json != nil {
		return m.json.Data.Deserialize(v)
	}
	return errors.New("the data is nil")
}

func (m *Message) String() string {
	b, err := JsonCodec.Encode(m)
	if err != nil {
		return "-"
	}
	return string(b)
}
