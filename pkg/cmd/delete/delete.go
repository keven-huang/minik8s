package delete

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"net/http"
	"net/url"
)

// DeleteOptions is the commandline options for 'delete' sub command
type DeleteOptions struct {
	DeleteAll bool
}

// NewDeleteOptions returns an initialized DeleteOptions instance
func NewDeleteOptions() *DeleteOptions {
	return &DeleteOptions{
		DeleteAll: false,
	}
}

// NewCmdDelete returns new initialized instance of Delete sub command
func NewCmdDelete() *cobra.Command {
	o := NewDeleteOptions()

	cmd := &cobra.Command{
		Use:   "delete (TYPE [NAME | -all])",
		Short: "Delete a resource from stdin",
		Run: func(cmd *cobra.Command, args []string) {
			err := o.RunDelete(cmd, args)
			if err != nil {
				return
			}
		},
	}

	//usage := "to use to delete the resource"
	cmd.Flags().BoolVarP(&o.DeleteAll, "all", "a", false, "delete all.")

	return cmd
}

// RunDelete performs the creation
func (o *DeleteOptions) RunDelete(cmd *cobra.Command, args []string) error {
	// Send a POST request to kube-apiserver to delete the pod
	// 创建 PUT 请求
	if len(args) < 1 || args[0] != "pod" {
		fmt.Println("only support pod.")
		return nil
	}

	values := url.Values{}
	if o.DeleteAll {
		values.Add("all", "true")
	}
	if len(args) > 1 {
		values.Add("PodName", args[1])
		//body = bytes.NewBuffer([]byte(args[1]))
	}
	req, err := http.NewRequest("DELETE", apiconfig.Server_URL+apiconfig.POD_PATH+"?"+values.Encode(), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// 设置请求头，以指示请求体中包含表单数据
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
	fmt.Println("Response Body:", string(bodyBytes))

	fmt.Println("pod Delete successfully")
	return nil

}
