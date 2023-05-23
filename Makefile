clean:
	etcdctl del "" --prefix

run:
	/usr/local/go/bin/go run ./cmd/kube-apiserver/kube-apiserver.go > log/apiserver.log &
	sleep 8
	/usr/local/go/bin/go run ./cmd/kube-controller/kube-controller.go > log/controller.log &
	/usr/local/go/bin/go run ./cmd/kube-scheduler/kube-scheduler.go > log/scheduler.log &
	/usr/local/go/bin/go run ./cmd/kubelet/kubelet.go > log/kubelet.log &
	/usr/local/go/bin/go run ./cmd/kubelet/kube-controller-manager.go > log/kube-controller-manager.log &

stop:
	./scripts/linux/stop.sh