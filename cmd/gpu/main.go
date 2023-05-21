package main

import (
	"flag"
	gpu "minik8s/cmd/gpu/server"
)

func main() {
	jobName := flag.String("jobname", "", "job名称")
	errFile := flag.String("errfile", "", "错误文件")
	outFile := flag.String("outfile", "", "输出文件")
	flag.Parse()
	server := gpu.NewServer(*jobName, *outFile, *errFile)
	server.Run()
}
