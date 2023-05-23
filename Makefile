clean:
	etcdctl del "" --prefix

run:
	go run ./cmd/kube-apiserver/kube-apiserver.go > log/apiserver.log &
	sleep 2
	go run ./cmd/kube-controller/kube-controller.go > log/controller.log &
	go run ./cmd/kube-scheduler/kube-scheduler.go > log/scheduler.log &
	go run ./cmd/kubelet/kubelet.go > log/kubelet.log &

stop:
	./scripts/linux/stop.sh