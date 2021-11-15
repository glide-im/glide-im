package groupdao

import "go_im/pkg/db"

const keyGroupIDIncr = "im:group:incr:gid"

func GetNextGid() (int64, error) {
	result, err := db.Redis.Incr(keyGroupIDIncr).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}
