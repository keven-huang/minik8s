clean:
	etcdctl del "/api" --prefix
	#rm -rf ./bin 2>/dev/null || true
	pwd
	rm -rf /root/nginx 2>/dev/null || true
	rm ./iptable_ori 2>/dev/null || true

# deprecated
construct:
	mkdir -p bin
	/usr/local/go/bin/go build -o ./bin/kube-apiserver ./cmd/kube-apiserver/kube-apiserver.go
	/usr/local/go/bin/go build -o ./bin/kube-scheduler ./cmd/kube-scheduler/kube-scheduler.go
	/usr/local/go/bin/go build -o ./bin/kubelet ./cmd/kubelet/kubelet.go
	/usr/local/go/bin/go build -o ./bin/kube-controller-manager ./cmd/kube-controller-manager/kube-controller-manager.go
	/usr/local/go/bin/go build -o ./bin/kubeproxy ./cmd/kube-proxy/kubeproxy.go
	/usr/local/go/bin/go build -o ./bin/kubeservice ./cmd/kube-service/kubeservice.go
	/usr/local/go/bin/go build -o ./bin/kubectl ./cmd/kubectl/kubectl.go

# deprecated
build:
	/usr/local/go/bin/go build -o ./binx/kube-apiserver ./cmd/kube-apiserver/kube-apiserver.go
	/usr/local/go/bin/go build -o ./binx/kube-scheduler ./cmd/kube-scheduler/kube-scheduler.go
	/usr/local/go/bin/go build -o ./binx/kubelet ./cmd/kubelet/kubelet.go
	/usr/local/go/bin/go build -o ./binx/kube-controller-manager ./cmd/kube-controller-manager/kube-controller-manager.go
	/usr/local/go/bin/go build -o ./binx/kubeproxy ./cmd/kube-proxy/kubeproxy.go
	/usr/local/go/bin/go build -o ./binx/kubeservice ./cmd/kube-service/kubeservice.go

run:
	mkdir -p /root/nginx
	/usr/local/go/bin/go run ./cmd/kube-apiserver/kube-apiserver.go > log/apiserver.log &
	sleep 5
	/usr/local/go/bin/go run ./cmd/kube-scheduler/kube-scheduler.go > log/scheduler.log &
	/usr/local/go/bin/go run ./cmd/kubelet/kubelet.go --nodename=node1 --nodeip=127.0.0.1 --masterip=http://127.0.0.1:8080 > log/kubelet.log &
	/usr/local/go/bin/go run ./cmd/kube-controller-manager/kube-controller-manager.go > log/controller-manager.log &
	/usr/local/go/bin/go run ./cmd/kube-proxy/kubeproxy.go > log/kubeproxy.log &
	sleep 1
	/usr/local/go/bin/go run ./cmd/kube-service/kubeservice.go > log/kubeservice.log &
	sleep 15

m3:
	go run ./cmd/kubelet/kubelet.go --nodename=node3 --nodeip=192.168.1.11 --masterip=http://192.168.1.7:8080 > log/kubelet-m3.log &

stop:
	./scripts/linux/stop.sh
	sleep 2
	iptables-restore < /root/iptable_ori
	docker ps -aq --filter "name=^coreDNS" | xargs -r docker stop
	docker ps -aq --filter "name=^coreDNS" | xargs -r docker rm
	docker volume rm volume-coredns 2>/dev/null || true

kill:
	sudo docker ps -aq --filter "name=^my-replicaset|^test" | xargs -r docker stop
	sudo docker ps -aq --filter "name=^my-replicaset|^test" | xargs -r docker rm

testsvc:
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/service-example/server-pod.yaml
	sleep 5
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/service-example/server-pod2.yaml
	sleep 5
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/service-example/server-service.yaml
	sleep 10
	curl 11.1.1.1
	curl 11.1.1.1
	curl 11.1.1.1
	curl 11.1.1.1
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go get service webservice
	sleep 3
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/service-example/server-pod3.yaml
	sleep 10
	curl 11.1.1.1
	curl 11.1.1.1
	curl 11.1.1.1
	curl 11.1.1.1
	curl 11.1.1.1
	curl 11.1.1.1
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go get service webservice
	sleep 3
delsvc:
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete service webservice
	sleep 5
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete pod tinyserver1
	sleep 2
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete pod tinyserver2
	sleep 2
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete pod tinyserver3

# need project in /home/minik8s, webs in /home/webs
testdns:
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/dns-example/server-pod1.yaml
	sleep 5
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/dns-example/server-pod2.yaml
	sleep 5
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/dns-example/server-service1.yaml
	sleep 10
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/dns-example/server-service2.yaml
	sleep 10
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/dns-example/dns-example.yaml
	sleep 30
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/dns-example/user-pod.yaml
	sleep 10
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go get dns testDns
	sleep 5
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go get service service1
	sleep 3
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go get service service2
	sleep 3
	docker logs user1-server
	sleep 1
	curl hanjinbo.com/path1
	curl hanjinbo.com/path2

deltestdns:
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete dns testDns
	sleep 20
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete service service1
	sleep 5
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete service service2
	sleep 5
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete pod dns-example-pod1
	sleep 2
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete pod dns-example-pod2
	sleep 2
	/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete pod user1
	sleep 3
	docker volume ls -q --filter "name=^gatewayvolume" | xargs -r docker volume rm || true
	docker volume rm volume04 volume05 volume-usr || true