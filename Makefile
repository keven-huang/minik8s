clean:
	etcdctl del "/api" --prefix

run:
	go run ./cmd/kube-apiserver/kube-apiserver.go > log/apiserver.log &
	sleep 10
	go run ./cmd/kube-scheduler/kube-scheduler.go > log/scheduler.log &
	go run ./cmd/kubelet/kublet.go > log/kubelet.log &
	go run ./cmd/kube-service/kubeservice.go > log/kubeservice.log &
	go run ./cmd/kube-proxy/kubeproxy.go > log/kubeproxy.log &

stop:
	./scripts/linux/stop.sh