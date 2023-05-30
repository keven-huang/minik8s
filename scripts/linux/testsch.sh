#!/bin/bash
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/sche-example/sche-pod.yaml
sleep 1
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/sche-example/sche-pod2.yaml
sleep 1
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/sche-example/sche-pod3.yaml
sleep 1
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/sche-example/sche-pod4.yaml
sleep 1
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/sche-example/sche-pod5.yaml
sleep 1
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go create -f ./cmd/kubectl/sche-example/sche-pod6.yaml
sleep 5
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go get pod sche