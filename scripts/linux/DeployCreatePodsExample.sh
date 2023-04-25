# 直接在云主机上编译，由于依赖一些github的包，可能会下不下来

go build -o ../../build/createPodsExample ../../cmd/kubelet/dockerClient/createPodsExample.go