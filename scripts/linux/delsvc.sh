#!/bin/bash

/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete service webservice
sleep 5
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete pod tinyserver1
sleep 2
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete pod tinyserver2
sleep 2
/usr/local/go/bin/go run ./cmd/kubectl/kubectl.go delete pod tinyserver3