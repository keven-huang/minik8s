clean:
	etcdctl del "/api" --prefix

run:
	go run ./cmd/kube-apiserver/kube-apiserver.go > log/apiserver.log &
	sleep 10
	go run ./cmd/kube-scheduler/kube-scheduler.go > log/scheduler.log &
	sleep 1
	go run ./cmd/kubelet/kublet.go > log/kubelet.log &
	sleep 1
	go run ./cmd/kube-proxy/kubeproxy.go > log/kubeproxy.log &
	sleep 1
	go run ./cmd/kube-service/kubeservice.go > log/kubeservice.log &


stop:
	./scripts/linux/stop.sh