package msgdao

import "go_im/pkg/db"

var Comm CommonDao = commonDao{}

type commonDao struct {
}

func (commonDao) GetMessageID() (int64, error) {
	// TODO 2021-12-17 16:57:04
	result, err := db.Redis.Incr("im:msg:id:incr").Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}
