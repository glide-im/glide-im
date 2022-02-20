package json

import (
	"encoding/json"
	"go_im/pkg/logger"
)

type Data struct {
	des interface{}
}

func NewData(d interface{}) Data {
	return Data{
		des: d,
	}
}

func (d *Data) UnmarshalJSON(bytes []byte) error {
	d.des = bytes
	return nil
}

func (d *Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.des)
}

func (d *Data) Bytes() []byte {
	bytes, ok := d.des.([]byte)
	if ok {
		return bytes
	}
	marshalJSON, err := d.MarshalJSON()
	if err != nil {
		logger.E("message data marshal json error %v", err)
		return nil
	}
	return marshalJSON
}

func (d *Data) Deserialize(i interface{}) error {
	s, ok := d.des.([]byte)
	if ok {
		return json.Unmarshal(s, i)
	}
	return nil
}

type CommMessage struct {
	Ver    int64
	Seq    int64
	Action string
	Data   Data
}
