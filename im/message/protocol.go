package message

import "encoding/json"

func UnmarshallJson(json_ string, i interface{}) error {
	return json.Unmarshal([]byte(json_), i)
}
