// Copyright 2016 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/client/v3"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
	// endpoints      = []string{"localhost:2379", "localhost:22379", "localhost:32379"}
	endpoints = []string{"localhost:2379"}
)

func main() {
	// Log + client connect + Put
	Example()
	// client connect + Put
	ExampleKV_put()
	// recommend error handling, example with empty key.
	ExampleKV_putErrorHandling()
	// client connect + Put + Get
	ExampleKV_get()
	// put, get返回的第一个参数是一个指针, get的时候传进去这个指针能查到过去的版本
	ExampleKV_getWithRev()
	// 前缀查询并排序后返回
	ExampleKV_getSortedPrefix()
	// delete (prefix)
	ExampleKV_delete()
	// 压缩compact
	ExampleKV_compact()
	// important! 事务操作
	ExampleKV_txn()
	//执行多个不同的操作(非事务)
	ExampleKV_do()
	// watch 某个key
	// 注意这个会spin,需要在命令行同时执行
	ExampleWatcher_watch()
	//ExampleWatcher_watchWithPrefix()
	//ExampleWatcher_watchWithRange()
	//ExampleWatcher_watchWithProgressNotify()
}

func Example() {
	// 3.5.8的log有很多冗余的，跟3.3.27不太一样
	//clientv3.SetLogger(grpclog.NewLoggerV2(os.Stderr, os.Stderr, os.Stderr))

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // make sure to close the client

	_, err = cli.Put(context.TODO(), "foo", "bar")
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleKV_put() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.Put(ctx, "sample_key", "sample_value")
	cancel()
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleKV_putErrorHandling() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.Put(ctx, "", "sample_value")
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			fmt.Printf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			fmt.Printf("ctx is attached with a deadline is exceeded: %v\n", err)
		//case rpctypes.ErrEmptyKey:
		//	fmt.Printf("client-side error: %v\n", err)
		default:
			fmt.Printf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
	}
	// Output: client-side error: etcdserver: key is not provided
}

func ExampleKV_get() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	_, err = cli.Put(context.TODO(), "foo", "bar")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "foo")
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	// Output: foo : bar
}

func ExampleKV_getWithRev() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	presp, err := cli.Put(context.TODO(), "foo", "bar1")
	if err != nil {
		log.Fatal(err)
	}
	_, err = cli.Put(context.TODO(), "foo", "bar2")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "foo", clientv3.WithRev(presp.Header.Revision))
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	// Output: foo : bar1
}

func ExampleKV_getSortedPrefix() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	for i := range make([]int, 3) {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		_, err = cli.Put(ctx, fmt.Sprintf("key_%d", i), "value")
		cancel()
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "key", clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	// Output:
	// key_2 : value
	// key_1 : value
	// key_0 : value
}

func ExampleKV_delete() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// count keys about to be deleted
	gresp, err := cli.Get(ctx, "key", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}

	// delete the keys
	dresp, err := cli.Delete(ctx, "key", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted all keys:", int64(len(gresp.Kvs)) == dresp.Deleted)
	// Output:
	// Deleted all keys: true
}

func ExampleKV_compact() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, "foo")
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	compRev := resp.Header.Revision // specify compact revision of your choice

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.Compact(ctx, compRev)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleKV_txn() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	kvc := clientv3.NewKV(cli)

	_, err = kvc.Put(context.TODO(), "key", "xyz")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err = kvc.Txn(ctx).
		// txn value comparisons are lexical
		If(clientv3.Compare(clientv3.Value("key"), ">", "abc")).
		// the "Then" runs, since "xyz" > "abc"
		Then(clientv3.OpPut("key", "XYZ")).
		// the "Else" does not run
		Else(clientv3.OpPut("key", "ABC")).
		Commit()
	cancel()
	if err != nil {
		log.Fatal(err)
	}

	gresp, err := kvc.Get(context.TODO(), "key")
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range gresp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	// Output: key : XYZ
}

func ExampleKV_do() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ops := []clientv3.Op{
		clientv3.OpPut("put-key", "123"),
		clientv3.OpGet("put-key"),
		clientv3.OpPut("put-key", "456")}

	for _, op := range ops {
		if _, err := cli.Do(context.TODO(), op); err != nil {
			log.Fatal(err)
		}
	}
}

// watch相关

func ExampleWatcher_watch() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), "foo")
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
	// PUT "foo" : "bar"
}

func ExampleWatcher_watchWithPrefix() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), "foo", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
	// PUT "foo1" : "bar"
}

func ExampleWatcher_watchWithRange() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// watches within ['foo1', 'foo4'), in lexicographical order
	rch := cli.Watch(context.Background(), "foo1", clientv3.WithRange("foo4"))
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
	// PUT "foo1" : "bar"
	// PUT "foo2" : "bar"
	// PUT "foo3" : "bar"
}

func ExampleWatcher_watchWithProgressNotify() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}

	rch := cli.Watch(context.Background(), "foo", clientv3.WithProgressNotify())
	wresp := <-rch
	fmt.Printf("wresp.Header.Revision: %d\n", wresp.Header.Revision)
	fmt.Println("wresp.IsProgressNotify:", wresp.IsProgressNotify())
	// wresp.Header.Revision: 0
	// wresp.IsProgressNotify: true
}
