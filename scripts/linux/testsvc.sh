#!/bin/bash

/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/service-example/server-pod.yaml
sleep 1
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/service-example/server-pod2.yaml
sleep 1
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/service-example/server-service.yaml
sleep 10
curl 11.1.1.1
curl 11.1.1.1
curl 11.1.1.1
curl 11.1.1.1
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go get service webservice
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/service-example/server-pod3.yaml
sleep 5
curl 11.1.1.1
curl 11.1.1.1
curl 11.1.1.1
curl 11.1.1.1
curl 11.1.1.1
curl 11.1.1.1
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go get service webservice