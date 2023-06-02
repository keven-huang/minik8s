package apply

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/util/web"
	"net/url"
	"os"
	"strconv"
)

type ApplyOptions struct {
	Filename string
}

func NewApplyOptions() *ApplyOptions {
	return &ApplyOptions{
		Filename: "",
	}
}

func NewCmdApply() *cobra.Command {
	o := NewApplyOptions()

	cmd := &cobra.Command{
		Use:   "apply [-f FILENAME]",
		Short: "Apply a resource from a file or from stdin",
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.Filename == "" {
				fmt.Print("Error: must specify one -f\n")
				return errors.New("must specify one -f")
			}
			//fmt.Println(o.Filename, args)
			err := o.RunApply(cmd, args)
			if err != nil {
				return err
			}
			return nil
		},
	}

	usage := "to use to apply the resource"
	cmd.Flags().StringVarP(&o.Filename, "filename", "f", "", "filename "+usage)

	return cmd
}

func (o *ApplyOptions) RunApply(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("[kubectl] [apply] [RunApply] unexpected args: %v", args)
	}

	filename := o.Filename
	fmt.Printf("filename: %s\n", filename)

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
	case "ReplicaSet":
		err = o.RunApplyReplicaSet(cmd, args, yamlFile)
	}

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

	return nil
}

func (o *ApplyOptions) RunApplyReplicaSet(cmd *cobra.Command, args []string, yamlFile []byte) error {
	r := &core.ReplicaSet{}
	err := yaml.Unmarshal(yamlFile, r)
	if err != nil {
		return err
	}

	err = ApplyReplicaSet(r)
	if err != nil {
		return err
	}
	return nil
}

func ApplyReplicaSet(r *core.ReplicaSet) error {
	// 序列化 调试用法，前面多两个空格，易于阅读
	//data, err := json.MarshalIndent(pod, "", "  ")
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Println("[kubectl] [apply] [RunApplyReplicaSet] failed to marshal:", err)
	} else {
		//fmt.Println("[kubectl] [apply] [RunApplyReplicaSet] ", string(data))
	}

	values := url.Values{}
	values.Add("replicas", strconv.Itoa(int(*r.Spec.Replicas)))
	err = web.SendHttpRequest("POST", apiconfig.Server_URL+apiconfig.REPLICASET_PATH+"?"+values.Encode(),
		web.WithPrefix("[kubectl] [apply] [RunApplyReplicaSet] "),
		web.WithBody(bytes.NewBuffer(data)),
		web.WithLog(true))
	if err != nil {
		return err
	}

	return nil
}
