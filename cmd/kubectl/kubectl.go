package main

import (
	"fmt"
	"minik8s/pkg/cmd"
	"os"
)

func main() {
	var rootCmd = cmd.NewKubectlCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
