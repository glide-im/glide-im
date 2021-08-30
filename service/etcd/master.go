package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go_im/pkg/logger"
	"time"
)

type Master struct {
	members map[string]*clientv3.Member
	cli     *clientv3.Client
}

func NewMaster(endpoints []string) *Master {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		Username:    "",
		Password:    "",
		Logger:      nil,
	})
	if err != nil {
		panic(err)
	}
	return &Master{
		members: make(map[string]*clientv3.Member),
		cli:     c,
	}
}

func (m *Master) WatchWorkers() {

	watchChan := m.cli.Watch(context.TODO(), "workers", clientv3.WithPrefix())

	for response := range watchChan {

		for _, event := range response.Events {

			if event.Type.String() == "PUT" {
				logger.D("worker online:%s", event.Kv.String())

			} else if event.Type.String() == "DELETE" {
				logger.D("worker offline:%s", event.Kv.String())
			}
		}
	}
}

func (m *Master) Run() {
	m.WatchWorkers()
}
