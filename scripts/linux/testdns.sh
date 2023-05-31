#!/bin/bash
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/dns-example/server-pod1.yaml
sleep 1
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/dns-example/server-pod2.yaml
sleep 1
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/dns-example/server-service1.yaml
sleep 5
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/dns-example/server-service2.yaml
sleep 5
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/dns-example/dns-example.yaml
sleep 35
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/dns-example/user-pod.yaml
sleep 5
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go get dns testDns
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go get service service1
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go get service service2
#docker logs user1-server
sleep 1
curl hanjinbo.com/path1
curl hanjinbo.com/path2