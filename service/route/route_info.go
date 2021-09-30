package route

import "sync"

var (
	mutexDevice = sync.Mutex{}
	mutexGroup  = sync.Mutex{}
	devices     = map[int64]map[int64]string{}
	group       = map[int64]string{}
)

func putDeviceRoute(uid int64, device int64, addr string) {
	mutexDevice.Lock()
	defer mutexDevice.Unlock()
	d, ok := devices[uid]
	if ok {
		d[device] = addr
	} else {
		ds := map[int64]string{
			device: addr,
		}
		devices[uid] = ds
	}
}

func getDeviceRoute(uid int64, device int64) string {
	mutexDevice.Lock()
	defer mutexDevice.Unlock()
	d, ok := devices[uid]
	if ok {
		return d[device]
	}
	return ""
}

func removeDeviceRoute(uid int64, device int64) {
	mutexDevice.Lock()
	defer mutexDevice.Unlock()
	d, ok := devices[uid]
	if ok {
		delete(d, device)
	}
}

func putGroupRoute(gid int64, addr string) {
	mutexGroup.Lock()
	defer mutexGroup.Unlock()
	group[gid] = addr
}

func getGroupRoute(gid int64) string {
	mutexGroup.Lock()
	defer mutexGroup.Unlock()
	return group[gid]
}

func removeGroupRoute(gid int64) {
	mutexGroup.Lock()
	defer mutexGroup.Unlock()
	delete(group, gid)
}
