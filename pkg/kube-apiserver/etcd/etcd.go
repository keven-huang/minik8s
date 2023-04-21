package etcd

import (
	"context"
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

type Result struct {
	key string
	val string
}

// constructort of etcdstore
func InitEtcdStore() (*EtcdStore, error) {
	cli, err := clientv3.New(
		clientv3.Config{
			Endpoints:   []string{"http://127.0.0.1:30000"},
			DialTimeout: dialTimeout,
		})
	if err != nil {
		log.Fatal(err)
	}
	return &EtcdStore{client: cli}, err
}

// deconstructor of EtcdStore

// operations : put
func (store *EtcdStore) PutKey(key string, val string) error {
	_, err := store.client.Put(context.TODO(), key, val)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// operations : get
func (store *EtcdStore) Get(key string) ([]Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := store.client.Get(ctx, key)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Kvs) == 0 {
		return []Result{}, err
	} else {
		res := []Result{}
		for _, ev := range resp.Kvs {
			res = append(res, Result{
				key: string(ev.Key),
				val: string(ev.Value),
			})
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
