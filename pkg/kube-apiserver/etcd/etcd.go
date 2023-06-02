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
	if err != nil {
		fmt.Printf("connect to etcd failed, err: %v\n", err)
	}

	fmt.Println("connect to etcd success...")
	return &EtcdStore{client: cli}, err
}

// deconstructor of EtcdStore

// operations : put
func (store *EtcdStore) Put(key string, val string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	kv := clientv3.NewKV(store.client)
	_, err := kv.Put(ctx, key, val)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return err
}

// operations : get
func (store *EtcdStore) Get(key string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	kv := clientv3.NewKV(store.client)
	resp, err := kv.Get(ctx, key)
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

// operations : get
func (store *EtcdStore) GetWithPrefix(prefix string) ([]ListRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	kv := clientv3.NewKV(store.client)
	resp, err := kv.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		fmt.Printf("get to etcd failed, err:%v\n", err)
		log.Fatal(err)
	}
	if len(resp.Kvs) == 0 {
		return []ListRes{}, err
	} else {
		res := []ListRes{}
		for _, ev := range resp.Kvs {
			res = append(res, ListRes{ResourceVersion: ev.Version, Key: string(ev.Key), Value: string(ev.Value)})
		}
		return res, err
	}
}

func (store *EtcdStore) GetExact(prefix string) ([]ListRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	kv := clientv3.NewKV(store.client)
	resp, err := kv.Get(ctx, prefix)
	if err != nil {
		fmt.Printf("get to etcd failed, err:%v\n", err)
		log.Fatal(err)
	}
	if len(resp.Kvs) == 0 {
		return []ListRes{}, err
	} else {
		res := []ListRes{}
		for _, ev := range resp.Kvs {
			res = append(res, ListRes{ResourceVersion: ev.Version, Key: string(ev.Key), Value: string(ev.Value)})
		}
		return res, err
	}
}

// operations : get
func (store *EtcdStore) GetAll(prefix string) ([]ListRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	kv := clientv3.NewKV(store.client)
	resp, err := kv.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		fmt.Printf("get to etcd failed, err:%v\n", err)
		log.Fatal(err)
	}
	if len(resp.Kvs) == 0 {
		return []ListRes{}, err
	} else {
		res := []ListRes{}
		for _, ev := range resp.Kvs {
			res = append(res, ListRes{ResourceVersion: ev.Version, Key: string(ev.Key), Value: string(ev.Value)})
		}
		return res, err
	}
}

// operation : del
func (store *EtcdStore) Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := store.client.Delete(ctx, key)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// DelAll delete all keys with prefix
func (store *EtcdStore) DelAll(prefix string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	deleteRes, err := store.client.Delete(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	return deleteRes.Deleted, nil
}

// TO DO: close channel and client maintain
// watch
func (store *EtcdStore) Watch(key string, isPrefix bool) <-chan Event {
	watchChan := make(chan Event)
	watcher := func(c chan<- Event) {
		var wat clientv3.WatchChan
		if isPrefix {
			wat = store.client.Watch(context.Background(), key, clientv3.WithPrefix())
		} else {
			wat = store.client.Watch(context.Background(), key)
		}
		for w := range wat {
			for _, event := range w.Events {
				fmt.Println("[etcd] etcd have watched ", string(event.Kv.Key), " ", event.Type)
				//fmt.Println("[etcd] etcd have watched")
				//fmt.Print(string(event.Kv.Key), " ", event.Type, "\n")
				fmt.Print(event.Kv.Version, " ", event.Kv.CreateRevision, " ", event.Kv.ModRevision, "\n")
				fmt.Println("----------")
				var watchedEvent Event
				watchedEvent.Type = getType(event)
				watchedEvent.Key = string(event.Kv.Key)
				watchedEvent.Val = string(event.Kv.Value)
				c <- watchedEvent
			}
		}
		close(c)
		log.Println("watcher closed")
	}
	go watcher(watchChan)
	return watchChan
}
