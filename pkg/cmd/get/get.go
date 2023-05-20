package get

import (
	"encoding/json"
	"fmt"
	"io"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
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
	Pods    []string `json:"Pods"`
	Message string   `json:"message"`
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
	if len(args) > 1 {
		values.Add("PodName", args[1])
		//body = bytes.NewBuffer([]byte(args[1]))
	}
	req, err := http.NewRequest("GET", apiconfig.Server_URL+apiconfig.POD_PATH+"?"+values.Encode(), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	// 打印响应结果
	fmt.Println("Response Status:", resp.Status)
	// 读取响应主体内容到字节数组
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	// 将字节数组转换为字符串并打印
	//fmt.Println("Response Body:", string(bodyBytes))
	var s GetRespond
	json.Unmarshal(bodyBytes, &s)
	fmt.Println(s.Pods)

	fmt.Println("pod Get successfully")
	return nil

}
