#!/bin/bash
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete dns testDns
sleep 10
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
docker volume ls -q --filter "name=^volume0" | xargs -r docker volume rm || true
docker volume ls -q --filter "name=^volume-usr" | xargs -r docker volume rm || true