package main

import (
	"errors"
	"minik8s/pkg/kube-apiserver/etcd"
	"testing"

	"github.com/go-playground/assert/v2"
	"gotest.tools/v3/assert"
)

func TestEtcd(t *testing.T) {
	// initialize etcd
	etcdstore, err := etcd.InitEtcdStore()
	if err != nil {
		assert.Error(t, errors.New("etcd establish wrong"), "")
	}
	etcdstore.Put("hello", "1023")
	res, err := etcdstore.Get("hello")
	if err != nil {
		assert.Error(t, errors.New("etcd get wrong"), "")
	}
}
