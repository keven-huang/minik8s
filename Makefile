clean:
	etcdctl del "/api" --prefix
	#rm -rf ./bin 2>/dev/null || true
	pwd
	rm -rf /root/nginx 2>/dev/null || true
	rm ./iptable_ori 2>/dev/null || true

construct:
	mkdir -p /root/nginx
	mkdir -p bin
	iptables-save > ./iptable_ori
	/usr/local/go/bin/go build -o ./binx/kube-apiserver ./cmd/kube-apiserver/kube-apiserver.go
	/usr/local/go/bin/go build -o ./binx/kube-scheduler ./cmd/kube-scheduler/kube-scheduler.go
	/usr/local/go/bin/go build -o ./binx/kubelet ./cmd/kubelet/kubelet.go
	/usr/local/go/bin/go build -o ./binx/kube-controller-manager ./cmd/kube-controller-manager/kube-controller-manager.go
	/usr/local/go/bin/go build -o ./binx/kubeproxy ./cmd/kube-proxy/kubeproxy.go
	/usr/local/go/bin/go build -o ./binx/kubeservice ./cmd/kube-service/kubeservice.go

build:
	/usr/local/go/bin/go build -o ./binx/kube-apiserver ./cmd/kube-apiserver/kube-apiserver.go
	/usr/local/go/bin/go build -o ./binx/kube-scheduler ./cmd/kube-scheduler/kube-scheduler.go
	/usr/local/go/bin/go build -o ./binx/kubelet ./cmd/kubelet/kubelet.go
	/usr/local/go/bin/go build -o ./binx/kube-controller-manager ./cmd/kube-controller-manager/kube-controller-manager.go
	/usr/local/go/bin/go build -o ./binx/kubeproxy ./cmd/kube-proxy/kubeproxy.go
	/usr/local/go/bin/go build -o ./binx/kubeservice ./cmd/kube-service/kubeservice.go

run:
	./binx/kube-apiserver > log/apiserver.log &
	sleep 5
	./binx/kube-scheduler > log/scheduler.log &
	./binx/kubelet --nodename=node1 --nodeip=127.0.0.1 --masterip=http://127.0.0.1:8080 > log/kubelet.log &
	./binx/kube-controller-manager > log/controller-manager.log &
	./binx/kubeproxy > log/kubeproxy.log &
	sleep 1
	./binx/kubeservice > log/kubeservice.log &


m3:
	go run ./cmd/kubelet/kubelet.go --nodename=node3 --nodeip=192.168.1.11 --masterip=http://192.168.1.7:8080 > log/kubelet-m3.log &

stop:
	./scripts/linux/stop.sh
	sleep 2
ifeq (iptable_ori, $(wildcard iptable_ori))
	iptables-restore < iptable_ori
endif
	docker ps -aq --filter "name=^coreDNS" | xargs -r docker stop
	docker ps -aq --filter "name=^coreDNS" | xargs -r docker rm

kill:
	sudo docker ps -aq --filter "name=^my-replicaset|^test" | xargs -r docker stop
	sudo docker ps -aq --filter "name=^my-replicaset|^test" | xargs -r docker rm


