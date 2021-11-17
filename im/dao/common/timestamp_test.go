package common

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnixTime_MarshalJSON(t *testing.T) {
	var ut Timestamp = Timestamp(time.Now())
	s, err := json.Marshal(&ut)
	t.Log(string(s), err)
}

func TestUnixTime_UnmarshalJSON(t *testing.T) {
	ut := Timestamp(time.Now())
	s, _ := json.Marshal(&ut)

	o := new(Timestamp)
	err := json.Unmarshal(s, o)
	t.Log(o.String(), err)
}
