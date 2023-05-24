package delete

import (
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/util/web"
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
		Use:   "delete TYPE [NAME | -all]",
		Short: "Delete a resource from stdin",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.RunDelete(cmd, args)
			if err != nil {
				return err
			}
			return nil
		},
	}

	//usage := "to use to delete the resource"
	cmd.Flags().BoolVarP(&o.DeleteAll, "all", "a", false, "delete all.")

	return cmd
}

// RunDelete performs the creation
func (o *DeleteOptions) RunDelete(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		fmt.Println("[kubectl] [delete] [RunDelete] must have TYPE.")
		return nil
	}

	switch args[0] {
	case "pod":
		{
			return o.RunDeletePod(cmd, args)
		}
	case "replicaset":
		{
			return o.RunDeleteReplicaSet(cmd, args)
		}
	case "service":
		{
			return o.RunDeleteService(cmd, args)
		}
	default:
		{
			fmt.Printf("[kubectl] [delete] [RunDelete] %s is not supported.\n", args[0])
			return nil
		}
	}

}

func (o *DeleteOptions) RunDeletePod(cmd *cobra.Command, args []string) error {
	prefix := "[kubectl] [delete] [RunDeletePod] "
	values := url.Values{}
	if o.DeleteAll {
		values.Add("all", "true")
	}
	if len(args) > 1 {
		values.Add("PodName", args[1])
		//body = bytes.NewBuffer([]byte(args[1]))
	}

	err := web.SendHttpRequest("DELETE", apiconfig.Server_URL+apiconfig.POD_PATH+"?"+values.Encode(),
		web.WithPrefix(prefix),
		web.WithLog(true))
	if err != nil {
		return err
	}

	fmt.Println("pod Delete successfully")
	return nil
}

func (o *DeleteOptions) RunDeleteReplicaSet(cmd *cobra.Command, args []string) error {
	prefix := "[kubectl] [delete] [RunDelete] "
	values := url.Values{}
	if o.DeleteAll {
		values.Add("all", "true")
	}
	if len(args) > 1 {
		values.Add("Name", args[1])
		//body = bytes.NewBuffer([]byte(args[1]))
	}

	err := web.SendHttpRequest("DELETE", apiconfig.Server_URL+apiconfig.REPLICASET_PATH+"?"+values.Encode(),
		web.WithPrefix(prefix),
		web.WithLog(true))
	if err != nil {
		return err
	}

	fmt.Println("Delete ReplicaSet successfully")
	return nil
}

func (o *DeleteOptions) RunDeleteService(cmd *cobra.Command, args []string) error {
	// 创建 PUT 请求
	if len(args) < 1 || args[0] != "service" {
		fmt.Println("only support service.")
		return nil
	}

	values := url.Values{}
	if o.DeleteAll {
		values.Add("all", "true")
	}
	if len(args) > 1 {
		values.Add("ServiceName", args[1])
		//body = bytes.NewBuffer([]byte(args[1]))
	}

	err := web.SendHttpRequest("DELETE", apiconfig.Server_URL+apiconfig.SERVICE_PATH+"?"+values.Encode(),
		web.WithPrefix("[kubectl] [delete] [DeleteService] "),
		web.WithLog(true))
	if err != nil {
		return err
	}

	fmt.Println("service Delete successfully")
	return nil
}
