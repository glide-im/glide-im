package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go_im/pkg/logger"
	"time"
)

type WorkerInfo struct {
	Name string
	IP   string
	CPU  int
	Mem  int
}

func (receiver *WorkerInfo) String() string {
	return fmt.Sprintf("{\"Name\":\"%s\", \"IP\":\"%s\", \"CPU\":\"%d\", \"Mem\":\"%d\"}",
		receiver.Name, receiver.IP, receiver.CPU, receiver.Mem)
}

type Worker struct {
	Name string
	IP   string
	Cli  *clientv3.Client
}

func NewWorker(name string, IP string, endpoints []string) *Worker {

	config := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}

	return &Worker{Name: name, IP: IP, Cli: client}
}

func (w *Worker) Watch() {
	watchChan := w.Cli.Watch(context.Background(), "key")
	for response := range watchChan {
		for _, event := range response.Events {
			logger.D("key on event: %s, %v", "key", event.Kv.String())
		}
	}
}

func (w *Worker) Heartbeat() {
	for {
		leaseResp, err := w.Cli.Lease.Grant(context.TODO(), 10)
		info := WorkerInfo{
			Name: w.Name,
			IP:   w.IP,
			CPU:  0,
			Mem:  0,
		}
		if err != nil {

		}
		_, err = w.Cli.Put(context.TODO(), "workers."+w.Name, info.String(), clientv3.WithLease(leaseResp.ID))
		if err != nil {

		}
		time.Sleep(3 * time.Second)
	}
}

func (w *Worker) Run() {
	go w.Watch()
	go w.Heartbeat()
}
