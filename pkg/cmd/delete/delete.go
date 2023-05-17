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

	err := web.SendHttpRequest("DELETE", apiconfig.Server_URL+apiconfig.POD_PATH+"?"+values.Encode(),
		web.WithPrefix("[kubectl] [delete] [RunDelete] "),
		web.WithLog(true))
	if err != nil {
		return err
	}

	fmt.Println("pod Delete successfully")
	return nil

}
