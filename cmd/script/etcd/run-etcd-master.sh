etcd \
--data-dir=data.etcd \
--name master \
--initial-advertise-peer-urls http://127.0.0.1:2380 \
--listen-peer-urls http://127.0.0.1:2380 \
--advertise-client-urls http://127.0.0.1:2379 \
--listen-client-urls http://127.0.0.1:2379 \
--initial-cluster master=http://127.0.0.1:2380,slave1=http://127.0.0.1:2382,slave2=http://127.0.0.1:2384,slave3=http://127.0.0.1:2386 \
--initial-cluster-state new \
--initial-cluster-token my-token
