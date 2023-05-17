package get

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/kube-apiserver/etcd"
	"minik8s/pkg/util/web"
	"net/url"
)

// GetOptions is the commandline options for 'get' sub command
type GetOptions struct {
	GetAll bool
}

// NewGetOptions returns an initialized GetOptions instance
func NewGetOptions() *GetOptions {
	return &GetOptions{
		GetAll: false,
	}
}

// NewCmdGet returns new initialized instance of Get sub command
func NewCmdGet() *cobra.Command {
	o := NewGetOptions()

	cmd := &cobra.Command{
		Use:   "get (TYPE [NAME | -all])",
		Short: "Get a resource from stdin",
		Run: func(cmd *cobra.Command, args []string) {
			err := o.RunGet(cmd, args)
			if err != nil {
				return
			}
		},
	}

	//usage := "to use to get the resource"
	cmd.Flags().BoolVarP(&o.GetAll, "all", "a", false, "get all.")

	return cmd
}

type GetRespond struct {
	Results []etcd.ListRes `json:"Results"`
	Message string         `json:"message"`
}

// RunGet performs the creation
func (o *GetOptions) RunGet(cmd *cobra.Command, args []string) error {
	// Send a POST request to kube-apiserver to get the pod
	// 创建 PUT 请求
	if len(args) < 1 || args[0] != "pod" {
		fmt.Println("only support pod.")
		return nil
	}

	values := url.Values{}
	if o.GetAll {
		values.Add("all", "true")
	}

	values.Add("prefix", "true")

	if len(args) > 1 {
		values.Add("Name", args[1])
		//body = bytes.NewBuffer([]byte(args[1]))
	}

	bodyBytes := make([]byte, 0)

	err := web.SendHttpRequest("GET", apiconfig.Server_URL+apiconfig.POD_PATH+"?"+values.Encode(),
		web.WithPrefix("[kubectl] [get] [RunGet] "),
		web.WithLog(true),
		web.WithBodyBytes(&bodyBytes))
	if err != nil {
		return err
	}

	// 将字节数组转换为字符串并打印
	var s GetRespond
	json.Unmarshal(bodyBytes, &s)
	fmt.Println("[kubectl] [get] [RunGet] Results:", s.Results)

	fmt.Println("[kubectl] [get] [RunGet] pod Get successfully")

	return nil
}
