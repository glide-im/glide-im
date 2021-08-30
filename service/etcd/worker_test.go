package etcd

import (
	"testing"
)

func TestNewWorker(t *testing.T) {

	master := NewMaster([]string{
		"http://127.0.0.1:2379",
		"http://127.0.0.1:2381",
		"http://127.0.0.1:2383",
		"http://127.0.0.1:2385"})

	master.WatchWorkers()
}
