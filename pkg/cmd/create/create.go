package create

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
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

	// 读取Pod YAML文件
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read the YAML file
	yamlFile, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	fmt.Println(string(yamlFile))

	// 解析YAML文件
	//var obj unstructured.Unstructured
	//decoder := yaml.NewYAMLOrJSONDecoder(&yamlFile, len(yamlFile))
	//if err := decoder.Decode(&obj); err != nil {
	//	panic(err)
	//}
	//
	//// 将解析后的对象转换为Pod对象
	//var pod corev1.Pod
	//if err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, &pod); err != nil {
	//	panic(err)
	//}
	//
	//// 输出Pod的名称
	//fmt.Printf("Pod name: %s\n", pod.ObjectMeta.Name)

	// Send a POST request to kube-apiserver to create the pod
	resp, err := http.Post("http://localhost:8001/api/v1/namespaces/default/pods",
		"application/yaml",
		bytes.NewBuffer(yamlFile))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create pod: %v", resp.Status)
	}

	fmt.Println("pod created successfully")
	return nil

}
