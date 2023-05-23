clean:
	etcdctl del "" --prefix

run:
	# 检查目录是否存在
	if [ ! -d "log/" ]; then
	  echo "log/ 目录不存在，开始创建..."
	  mkdir log/
	  echo "log/ 目录创建成功！"
	else
	  echo "log/ 目录已存在，无需创建。"
	fi

	/usr/local/go/bin/go run ./cmd/kube-apiserver/kube-apiserver.go > log/apiserver.log &
	sleep 5
	
	/usr/local/go/bin/go run ./cmd/kube-scheduler/kube-scheduler.go > log/scheduler.log &
	/usr/local/go/bin/go run ./cmd/kubelet/kubelet.go > log/kubelet.log &
	/usr/local/go/bin/go run ./cmd/kube-controller-manager/kube-controller-manager.go > log/kube-controller-manager.log &

stop:
	./scripts/linux/stop.sh