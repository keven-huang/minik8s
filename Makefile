clean:
	etcdctl del "/api" --prefix
	rm -rf ./bin
	rm -rf /root/nginx
	rm ./iptable_ori

construct:
	mkdir -p /root/nginx
	mkdir -p bin
	iptables-save > ./iptable_ori
	/usr/local/go/bin/go build -o ./bin/kube-apiserver ./cmd/kube-apiserver/kube-apiserver.go
	/usr/local/go/bin/go build -o ./bin/kube-scheduler ./cmd/kube-scheduler/kube-scheduler.go
	/usr/local/go/bin/go build -o ./bin/kubelet ./cmd/kubelet/kubelet.go
	/usr/local/go/bin/go build -o ./bin/kube-controller-manager ./cmd/kube-controller-manager/kube-controller-manager.go
	/usr/local/go/bin/go build -o ./bin/kubeproxy ./cmd/kube-proxy/kubeproxy.go
	/usr/local/go/bin/go build -o ./bin/kubeservice ./cmd/kube-service/kubeservice.go

build:
	/usr/local/go/bin/go build -o ./bin/kube-apiserver ./cmd/kube-apiserver/kube-apiserver.go
	/usr/local/go/bin/go build -o ./bin/kube-scheduler ./cmd/kube-scheduler/kube-scheduler.go
	/usr/local/go/bin/go build -o ./bin/kubelet ./cmd/kubelet/kubelet.go
	/usr/local/go/bin/go build -o ./bin/kube-controller-manager ./cmd/kube-controller-manager/kube-controller-manager.go
	/usr/local/go/bin/go build -o ./bin/kubeproxy ./cmd/kube-proxy/kubeproxy.go
	/usr/local/go/bin/go build -o ./bin/kubeservice ./cmd/kube-service/kubeservice.go

run:
	./bin/kube-apiserver > log/apiserver.log &
	sleep 5
	./bin/kube-scheduler > log/scheduler.log &
	./bin/kubelet --nodename=node1 --nodeip=127.0.0.1 --masterip=http://127.0.0.1:8080 > log/kubelet.log &
	./bin/kube-controller-manager > log/controller-manager.log &
	./bin/kubeproxy > log/kubeproxy.log &
	sleep 1
	./bin/kubeservice > log/kubeservice.log &


m3:
	go run ./cmd/kubelet/kubelet.go --nodename=node3 --nodeip=192.168.1.11 --masterip=http://192.168.1.7:8080 > log/kubelet-m3.log &

stop:
	./scripts/linux/stop.sh
	sleep 2
	iptables-restore < ./iptable_ori
	sudo docker ps -aq --filter "name=^coreDNS" | xargs -r docker stop
	sudo docker ps -aq --filter "name=^coreDNS" | xargs -r docker rm

kill:
	sudo docker ps -aq --filter "name=^my-replicaset|^test" | xargs -r docker stop
	sudo docker ps -aq --filter "name=^my-replicaset|^test" | xargs -r docker rm


