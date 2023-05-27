clean:
	etcdctl del "/api" --prefix

run:
	/usr/local/go/bin/go run ./cmd/kube-apiserver/kube-apiserver.go > log/apiserver.log &
	sleep 5
	/usr/local/go/bin/go run ./cmd/kube-scheduler/kube-scheduler.go > log/scheduler.log &
	/usr/local/go/bin/go run ./cmd/kubelet/kubelet.go --nodename=node1 --nodeip=127.0.0.1 --masterip=http://127.0.0.1:8080 > log/kubelet.log &
	/usr/local/go/bin/go run ./cmd/kube-controller-manager/kube-controller-manager.go > log/controller-manager.log &
	/usr/local/go/bin/go run ./cmd/kube-proxy/kubeproxy.go > log/kubeproxy.log &
	sleep 1
	/usr/local/go/bin/go run ./cmd/kube-service/kubeservice.go > log/kubeservice.log &


m3:
	go run ./cmd/kubelet/kubelet.go --nodename=node3 --nodeip=192.168.1.11 --masterip=http://192.168.1.7:8080 > log/kubelet-m3.log &

stop:
	./scripts/linux/stop.sh

kill:
	sudo docker ps -aq --filter "name=^my-replicaset|^test" | xargs -r docker stop
	sudo docker ps -aq --filter "name=^my-replicaset|^test" | xargs -r docker rm


