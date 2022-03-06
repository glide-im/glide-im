package cache

import "sync"

var gatewayCache = map[int64]string{}
var m = sync.RWMutex{}

func SetGateway(uid int64, device int64, gateway string) error {
	m.Lock()
	defer m.Unlock()
	gatewayCache[uid] = gateway
	return nil
}

func GetGateway(uid int64, device int64) (string, error) {
	m.RLock()
	defer m.RUnlock()
	s := gatewayCache[uid]
	return s, nil
}
