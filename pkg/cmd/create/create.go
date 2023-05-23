package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/service"
	myJson "minik8s/pkg/util/json"
	"minik8s/pkg/util/web"
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
		Use:   "create TYPE [-f FILENAME]",
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
	if len(args) < 1 {
		fmt.Println("[kubectl] [create] [RunCreate] must have TYPE.")
		return nil
	}

	switch args[0] {
	case "pod":
		{
			return o.RunCreatePod(cmd, args)
		}
	case "replicaset":
		{
			return o.RunCreateReplicaSet(cmd, args)
		}
	case "service":
		{
			return o.RunCreateService(cmd, args)
		}
	case "dns":
		{
			return o.RunCreateDNS(cmd, args)
		}
	default:
		{
			fmt.Printf("[kubectl] [create] [RunCreate] %s is not supported.\n", args[0])
			return nil
		}
	}
}

func (o *CreateOptions) RunCreatePod(cmd *cobra.Command, args []string) error {
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
		fmt.Println("[kubectl] [create] [RunCreatePod] failed to marshal:", err)
	} else {
		//fmt.Println("[kubectl] [create] [RunCreatePod]\n", string(data))
	}

	// Send a POST request to kube-apiserver to create the pod
	// 创建 PUT 请求
	err = web.SendHttpRequest("PUT", apiconfig.Server_URL+apiconfig.POD_PATH,
		web.WithPrefix("[kubectl] [create] [RunCreate] "),
		web.WithBody(bytes.NewBuffer(data)),
		web.WithLog(true))
	if err != nil {
		return err
	}

	return nil
}

func (o *CreateOptions) RunCreateReplicaSet(cmd *cobra.Command, args []string) error {
	filename := o.Filename

	r := &core.ReplicaSet{}
	err := myJson.GetFromYaml(filename, r)
	if err != nil {
		return err
	}

	// 序列化
	// 调试用法，前面多两个空格，易于阅读
	//data, err := json.MarshalIndent(pod, "", "  ")
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Println("[kubectl] [create] [RunCreateReplicaSet] failed to marshal:", err)
	} else {
		//fmt.Println("[kubectl] [create] [RunCreateReplicaSet] ", string(data))
	}

	// Send a POST request to kube-apiserver to create the replicaset
	// 创建 PUT 请求
	err = web.SendHttpRequest("PUT", apiconfig.Server_URL+apiconfig.REPLICASET_PATH,
		web.WithPrefix("[kubectl] [create] [RunCreateReplicaSet] "),
		web.WithBody(bytes.NewBuffer(data)),
		web.WithLog(true))
	if err != nil {
		return err
	}

	return nil
}

func (o *CreateOptions) RunCreateService(cmd *cobra.Command, args []string) error {
	filename := o.Filename

	r := &service.Service{}
	err := myJson.GetFromYaml(filename, r)
	if err != nil {
		return err
	}
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Println("[kubectl] [create] [RunCreateService] failed to marshal:", err)
	} else {
		//fmt.Println("[kubectl] [create] [RunCreateReplicaSet] ", string(data))
	}
	err = web.SendHttpRequest("PUT", apiconfig.Server_URL+apiconfig.SERVICE_PATH,
		web.WithPrefix("[kubectl] [create] [RunCreatesService] "),
		web.WithBody(bytes.NewBuffer(data)),
		web.WithLog(true))
	if err != nil {
		return err
	}

	return nil
}

func (o *CreateOptions) RunCreateDNS(cmd *cobra.Command, args []string) error {
	filename := o.Filename

	r := &core.DNS{}
	err := myJson.GetFromYaml(filename, r)
	if err != nil {
		return err
	}
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Println("[kubectl] [create] [RunCreateDNS] failed to marshal:", err)
	} else {
		//fmt.Println("[kubectl] [create] [RunCreateReplicaSet] ", string(data))
	}
	err = web.SendHttpRequest("PUT", apiconfig.Server_URL+apiconfig.DNS_PATH,
		web.WithPrefix("[kubectl] [create] [RunCreatesDNS] "),
		web.WithBody(bytes.NewBuffer(data)),
		web.WithLog(true))
	if err != nil {
		return err
	}

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
