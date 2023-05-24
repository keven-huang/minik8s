package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	// 创建Docker客户端
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// 获取容器的实时统计信息
	stats, err := cli.ContainerStats(context.Background(), "home_Etcd_1", true)
	if err != nil {
		panic(err)
	}
	defer stats.Body.Close()

	// 读取和处理统计信息
	var result types.StatsJSON
	for {
		err := readStats(stats.Body, &result)
		if err != nil {
			if err == io.EOF {
				// 读取完所有统计信息后退出循环
				break
			}
			panic(err)
		}

		// 输出 CPU 使用率
		cpuUsage := result.CPUStats.CPUUsage.TotalUsage
		fmt.Printf("CPU Usage: %d\n", cpuUsage)

		// 输出内存使用量
		memoryUsage := result.MemoryStats.Usage
		fmt.Printf("Memory Usage: %d\n", memoryUsage)

		// 输出网络流量
		for interfaceName, networkStats := range result.Networks {
			rxBytes := networkStats.RxBytes
			txBytes := networkStats.TxBytes
			fmt.Printf("Interface: %s, RX Bytes: %d, TX Bytes: %d\n", interfaceName, rxBytes, txBytes)
		}

		// 等待一段时间再继续读取统计信息
		time.Sleep(1 * time.Second)
	}
}

func readStats(reader io.Reader, v interface{}) error {
	dec := json.NewDecoder(reader)
	return dec.Decode(v)
}
