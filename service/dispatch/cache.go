package dispatch

import (
	"sync"
)

type routeCache struct {
	c  map[int64]string
	mu sync.RWMutex
}

func newRouteCache() *routeCache {
	r := new(routeCache)
	r.mu = sync.RWMutex{}
	r.c = map[int64]string{}
	return r
}

func (r *routeCache) getRoute(srvName string, id int64) string {
	r.mu.RLock()
	node, ok := r.c[id]
	r.mu.RUnlock()
	if !ok {
		// TODO 2022-3-11 17:27:53 Query redis
	}
	return node
}

func (r *routeCache) updateRoute(srvName string, id int64, node string) {
	r.mu.Lock()
	r.c[id] = node
	r.mu.Unlock()
}
