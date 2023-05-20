package file

import (
	"fmt"
	"os"
	"strings"
)

func MakeFile(val []byte, name string, path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	file, err := os.OpenFile(path+"/"+name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(val)
	if err != nil {
		return err
	}
	return nil
}

func GetBasePath(paths string) string {
	path := strings.Split(paths, "/")
	fmt.Println(path)
	return path[len(path)-1]
}
