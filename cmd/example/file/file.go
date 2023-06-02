package main

import "minik8s/pkg/util/file"

func main() {
	a := "hello world"
	file.MakeFile([]byte(a), "test.txt", ".")
	for {

	}
}
