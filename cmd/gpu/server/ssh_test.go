package gpu_server

import (
	"fmt"
	"testing"
)

func TestSShcmd(t *testing.T) {
	client := NewSSHClient(User, Pwd, Host, Port)
	cmd := "ls"
	backinfo, err := client.Run(cmd)
	if err != nil {
		fmt.Printf("failed to run shell,err=[%v]\n", err)
		return
	}
	fmt.Printf("%v back info: \n[%v]\n", cmd, backinfo)
}

func TestSShUpload(t *testing.T) {
	client := NewSSHClient(User, Pwd, Host, Port)
	cmd := "ls"
	backinfo, err := client.Run(cmd)
	if err != nil {
		fmt.Printf("failed to run shell,err=[%v]\n", err)
		return
	}
	fmt.Printf("%v back info: \n[%v]\n", cmd, backinfo)
}
