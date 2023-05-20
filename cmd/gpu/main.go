package main

import (
	"flag"
	"fmt"
	gpu "minik8s/cmd/gpu/server"
)

func main() {
	jobName := flag.String("jobname", "", "job名称")
	flag.Parse()
	fmt.Println("jobName:", *jobName)
	server := gpu.NewServer(*jobName)
	server.Run()
}
