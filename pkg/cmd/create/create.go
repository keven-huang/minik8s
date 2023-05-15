package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	myJson "minik8s/pkg/util/json"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// CreateOptions is the commandline options for 'create' sub command
type CreateOptions struct {
	Filename string
}

// NewCreateOptions returns an initialized CreateOptions instance
func NewCreateOptions() *CreateOptions {
	return &CreateOptions{
		Filename: "",
	}
}

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdCreate() *cobra.Command {
	o := NewCreateOptions()

	cmd := &cobra.Command{
		Use:   "create [-f FILENAME]",
		Short: "Create a resource from a file or from stdin",
		Run: func(cmd *cobra.Command, args []string) {
			if o.Filename == "" {
				fmt.Print("Error: must specify one -f\n\n")
				return
			}
			//fmt.Println(o.Filename, args)
			err := o.RunCreate(cmd, args)
			if err != nil {
				return
			}
		},
	}

	usage := "to use to create the resource"
	cmd.Flags().StringVarP(&o.Filename, "filename", "f", "", "filename "+usage)

	return cmd
}

// RunCreate performs the creation
func (o *CreateOptions) RunCreate(cmd *cobra.Command, args []string) error {
	fmt.Println("in run create")

	if len(args) != 0 {
		return fmt.Errorf("unexpected args: %v", args)
	}

	filename := o.Filename
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open file %s", filename)
	}
	defer file.Close()
	// Read the YAML file
	var dataMap map[string]interface{}

	yamlFile, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read yaml error")
	}

	err = yaml.Unmarshal(yamlFile, &dataMap)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// get kind
	var kind string
	kind = dataMap["kind"].(string)
	fmt.Println("kind:", kind)
	switch kind {
	case "Pod":
		err = createPod(yamlFile)
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func createPod(yamlfile []byte) error {
	// 解析文件
	pod := &core.Pod{}
	err := yaml.Unmarshal(yamlfile, pod)
	if err != nil {
		return fmt.Errorf("unmarshal yaml error")
	}
	// 序列化
	// 调试用法，前面多两个空格，易于阅读
	//data, err := json.MarshalIndent(pod, "", "  ")
	data, err := json.Marshal(pod)
	if err != nil {
		fmt.Println("failed to marshal person:", err)
	} else {
		// fmt.Println(string(data))
	}

	// FIX（hjm） : 这两步后续可以去掉，只是为了调试
	// 反序列化
	pod2 := &core.Pod{}
	err = json.Unmarshal(data, pod2)
	if err != nil {
		return err
	}
	fmt.Println(pod2)

	// 检查序列化和反序列化的结果
	myJson.CheckDeepEqual(pod, pod2)

	// Send a POST request to kube-apiserver to create the pod
	// 创建 PUT 请求
	req, err := http.NewRequest("PUT", apiconfig.Server_URL+apiconfig.POD_PATH, bytes.NewBuffer(data))
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
	fmt.Println("Response Body:", string(bodyBytes))

	fmt.Println("pod created successfully")
	return nil
}
