package route

import "sync"

var gatewayCache = map[int64]string{}
var groupCache = map[int64]GroupRouteInfo{}

var muGroup = sync.RWMutex{}
var muGate = sync.RWMutex{}

func GetGateway(uid int64, device int64) (string, error) {
	muGate.RLock()
	defer muGate.RUnlock()
	s := gatewayCache[uid]
	return s, nil
}

func GetGroup(gid int64) (GroupRouteInfo, error) {
	muGroup.RLock()
	defer muGroup.RUnlock()
	s := groupCache[gid]
	return s, nil
}

func SetGateway(uid int64, device int64, rt string) error {
	muGate.Lock()
	defer muGate.Unlock()
	gatewayCache[uid] = rt
	return nil
}

func setGroup(gid int64, rt GroupRouteInfo) error {
	muGroup.Lock()
	defer muGroup.Unlock()
	groupCache[gid] = rt
	return nil
}
