package etcd

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
)

type EtcdStore struct {
	client *clientv3.Client
}

// constructort of etcdstore
func InitEtcdStore() (*EtcdStore, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	fmt.Printf("in\n")
	if err != nil {
		fmt.Printf("connect to etcd failed, err: %v\n", err)
	}

	fmt.Println("connect to etcd success...")
	return &EtcdStore{client: cli}, err
}

// deconstructor of EtcdStore

// operations : put
func (store *EtcdStore) Put(key string, val string) error {
	ctx, cal := context.WithTimeout(context.Background(), requestTimeout)
	kv := clientv3.NewKV(store.client)
	_, err := kv.Put(ctx, key, val)
	cal()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	fmt.Printf("out\n")
	return err
}

// operations : get
func (store *EtcdStore) Get(key string) ([]string, error) {
	fmt.Printf("in\n")
	kv := clientv3.NewKV(store.client)
	resp, err := kv.Get(context.TODO(), key)
	if err != nil {
		fmt.Printf("get to etcd failed, err:%v\n", err)
		log.Fatal(err)
	}
	if len(resp.Kvs) == 0 {
		return []string{}, err
	} else {
		res := []string{}
		for _, ev := range resp.Kvs {
			res = append(res, string(ev.Value))
		}
		return res, err
	}
}

// operation : del
func (store *EtcdStore) Del(key string) error {
	_, err := store.client.Delete(context.TODO(), key)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// watch
func (store *EtcdStore) Listen(key string) {
	wat := store.client.Watch(context.Background(), key)
	for w := range wat {
		for _, v := range w.Events {
			fmt.Printf("type: %s key: %s value: %s", v.Type, v.Kv.Key, v.Kv.Value)
		}
	}
}

// perfix watch
