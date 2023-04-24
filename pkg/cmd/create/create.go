package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	myJson "minik8s/pkg/util/json"
	"net/http"
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
		Use:   "create -f FILENAME",
		Short: "Create a resource from a file or from stdin",
		Run: func(cmd *cobra.Command, args []string) {
			if o.Filename == "" {
				fmt.Print("Error: must specify one of -f\n\n")
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
	if len(args) != 0 {
		return fmt.Errorf("unexpected args: %v", args)
	}

	filename := o.Filename

	pod := &core.Pod{}
	err := myJson.GetFromYaml(filename, pod)
	if err != nil {
		return err
	}

	// 序列化
	// 调试用法，前面多两个空格，易于阅读
	//data, err := json.MarshalIndent(pod, "", "  ")
	data, err := json.Marshal(pod)
	if err != nil {
		fmt.Println("failed to marshal person:", err)
	} else {
		//fmt.Println(string(data))
	}

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

// "k8s.io/apimachinery/pkg/runtime/serializer/json"
//var s runtime.Serializer
//// Yaml decides whether yaml or json
//option := json.SerializerOptions{Yaml: false, Pretty: false, Strict: false}
//// json
//s = json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, option)
//
//// 将对象编码为字节数组
//buf := new(bytes.Buffer)
//if err := s.Encode(*pod, buf); err != nil {
//	panic(err)
//}

//obj, gvk, err := s.Decode([]byte(test.data), &schema.GroupVersionKind{Kind: "Test", Group: "other", Version: "blah"}, &core.Pod{})
//
//if !reflect.DeepEqual(test.expectedGVK, gvk) {
//	logTestCase(t, test)
//	t.Errorf("%d: unexpected GVK: %v", i, gvk)
//}
