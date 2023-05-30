package get

import (
	"encoding/json"
	"fmt"
	"log"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/kube-apiserver/etcd"
	"minik8s/pkg/util/web"
	"net/url"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

// GetOptions is the commandline options for 'get' sub command
type GetOptions struct {
	GetAll       bool
	WithNoPrefix bool
}

// NewGetOptions returns an initialized GetOptions instance
func NewGetOptions() *GetOptions {
	return &GetOptions{
		GetAll:       false,
		WithNoPrefix: false,
	}
}

// NewCmdGet returns new initialized instance of Get sub command
func NewCmdGet() *cobra.Command {
	o := NewGetOptions()

	cmd := &cobra.Command{
		Use:   "get TYPE [NAME | -all]",
		Short: "Get a resource from stdin",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.RunGet(cmd, args)
			if err != nil {
				return err
			}
			return nil
		},
	}

	//usage := "to use to get the resource"
	cmd.Flags().BoolVarP(&o.GetAll, "all", "a", false, "get all.")
	cmd.Flags().BoolVarP(&o.WithNoPrefix, "no-prefix", "n", false, "get with no prefix query")

	return cmd
}

type GetRespond struct {
	Results []etcd.ListRes `json:"Results"`
	Message string         `json:"message"`
}

// RunGet performs the creation
func (o *GetOptions) RunGet(cmd *cobra.Command, args []string) error {
	prefix := "[kubectl] [get] [RunGet] "

	if len(args) < 1 {
		fmt.Println(prefix, "must have TYPE.")
		return nil
	}

	switch args[0] {
	case "pod":
		{
			return o.RunGetPod(cmd, args)
		}
	case "replicaset":
		{
			return o.RunGetReplicaSet(cmd, args)
		}
	case "job":
		{
			return o.RunGetJob(cmd, args)
		}
	case "workflow":
		{
			return o.RunGetWorkflow(cmd, args)
		}
	case "function":
		{
			return o.RunGetFunction(cmd, args)
		}
	default:
		{
			fmt.Printf(prefix, "%s is not supported.\n", args[0])
			return nil
		}
	}

}

func (o *GetOptions) RunGetPod(cmd *cobra.Command, args []string) error {
	prefix := "[kubectl] [get] [RunGetPod] "
	values := url.Values{}
	if o.GetAll {
		values.Add("all", "true")
	}

	if o.WithNoPrefix {
		values.Add("prefix", "false")
	} else {
		values.Add("prefix", "true")
	}

	if len(args) > 1 {
		values.Add("Name", args[1])
		//body = bytes.NewBuffer([]byte(args[1]))
	}

	bodyBytes := make([]byte, 0)

	err := web.SendHttpRequest("GET", apiconfig.Server_URL+apiconfig.POD_PATH+"?"+values.Encode(),
		web.WithPrefix(prefix),
		web.WithLog(false),
		web.WithBodyBytes(&bodyBytes))
	if err != nil {
		return err
	}

	// 将字节数组转换为字符串并打印
	//var s GetRespond
	//json.Unmarshal(bodyBytes, &s)
	var res []etcd.ListRes
	json.Unmarshal(bodyBytes, &res)
	fmt.Println(prefix, "Pod Get successfully. Here are the results:")

	fmt.Println("total number:", len(res))

	table := uitable.New()
	table.MaxColWidth = 100
	table.RightAlign(10)
	table.AddRow("NAME", "NODE", "POD_IP", "STATUS", "Owner", "CreationTimestamp")
	for _, val := range res {
		pod := core.Pod{}
		err := json.Unmarshal([]byte(val.Value), &pod)
		if err != nil {
			log.Println(prefix, err)
			return err
		}
		var owner string
		if len(pod.OwnerReferences) > 0 {
			owner = pod.OwnerReferences[0].Name
		} else {
			owner = ""
		}

		table.AddRow(color.RedString(pod.Name),
			color.WhiteString(pod.Spec.NodeName),
			color.GreenString(pod.Status.PodIP),
			color.BlueString(string(pod.Status.Phase)),
			color.GreenString(owner),
			color.YellowString(pod.CreationTimestamp.Format(time.UnixDate)))
	}
	fmt.Println(table)

	return nil
}

func (o *GetOptions) RunGetReplicaSet(cmd *cobra.Command, args []string) error {
	prefix := "[kubectl] [get] [RunGetPod] "
	values := url.Values{}
	if o.GetAll {
		values.Add("all", "true")
	}

	if o.WithNoPrefix {
		values.Add("prefix", "false")
	} else {
		values.Add("prefix", "true")
	}

	if len(args) > 1 {
		values.Add("Name", args[1])
		//body = bytes.NewBuffer([]byte(args[1]))
	}

	bodyBytes := make([]byte, 0)

	err := web.SendHttpRequest("GET", apiconfig.Server_URL+apiconfig.REPLICASET_PATH+"?"+values.Encode(),
		web.WithPrefix(prefix),
		web.WithLog(false),
		web.WithBodyBytes(&bodyBytes))
	if err != nil {
		return err
	}

	// 将字节数组转换为字符串并打印
	//var s GetRespond
	//json.Unmarshal(bodyBytes, &s)
	var res []etcd.ListRes
	json.Unmarshal(bodyBytes, &res)
	fmt.Println(prefix, "ReplicaSet Get successfully. Here are the results:")

	fmt.Println("total number:", len(res))

	table := uitable.New()
	table.MaxColWidth = 100
	table.RightAlign(10)
	table.AddRow("NAME", "Desired", "Current", "CreationTimestamp")
	for _, val := range res {
		r := core.ReplicaSet{}
		err := json.Unmarshal([]byte(val.Value), &r)
		if err != nil {
			log.Println(prefix, err)
			return err
		}
		fmt.Println(r.Name, *r.Spec.Replicas, r.Status.Replicas, r.CreationTimestamp.Format(time.UnixDate))
		table.AddRow(color.RedString(r.Name),
			color.BlueString(strconv.Itoa(int(*r.Spec.Replicas))),
			color.BlueString(strconv.Itoa(int(r.Status.Replicas))),
			color.YellowString(r.CreationTimestamp.Format(time.UnixDate)))
	}
	fmt.Println(table)

	return nil
}

func (o *GetOptions) RunGetJob(cmd *cobra.Command, args []string) error {
	prefix := "[kubectl] [get] [RunGetJob] "
	values := url.Values{}
	if o.GetAll {
		values.Add("all", "true")
	}

	if o.WithNoPrefix {
		values.Add("prefix", "false")
	} else {
		values.Add("prefix", "true")
	}
	if len(args) > 1 {
		values.Add("Name", args[1])
		//body = bytes.NewBuffer([]byte(args[1]))
	}

	bodyBytes := make([]byte, 0)
	fmt.Println("ask cmd:", apiconfig.Server_URL+apiconfig.JOB_PATH+"?"+values.Encode())
	err := web.SendHttpRequest("GET", apiconfig.Server_URL+apiconfig.JOB_PATH+"?"+values.Encode(),
		web.WithPrefix(prefix),
		web.WithLog(false),
		web.WithBodyBytes(&bodyBytes))
	if err != nil {
		return err
	}
	var res []etcd.ListRes
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		log.Println(prefix, err)
		return err
	}
	fmt.Println(prefix, "Job Get successfully. Here are the results:")
	fmt.Println("total number:", len(res))

	table := uitable.New()
	table.MaxColWidth = 100
	table.RightAlign(10)
	table.AddRow("NAME", "STATUS")
	for _, val := range res {
		job := core.JobStatus{}
		err := json.Unmarshal([]byte(val.Value), &job)
		if err != nil {
			log.Println(prefix, err)
			return err
		}
		fmt.Println("job:", job)
		table.AddRow(color.RedString(job.JobName),
			color.BlueString(string(job.Status)))
	}
	fmt.Println(table)

	return nil
}

func (o *GetOptions) RunGetWorkflow(cmd *cobra.Command, args []string) error {
	prefix := "[kubectl] [get] [RunGetWorkflow] "
	values := url.Values{}
	if o.GetAll {
		values.Add("all", "true")
	}

	if o.WithNoPrefix {
		values.Add("prefix", "false")
	} else {
		values.Add("prefix", "true")
	}
	if len(args) > 1 {
		values.Add("Name", args[1])
		//body = bytes.NewBuffer([]byte(args[1]))
	}

	bodyBytes := make([]byte, 0)
	fmt.Println("ask cmd:", apiconfig.Server_URL+apiconfig.WORKFLOW_PATH+"?"+values.Encode())
	err := web.SendHttpRequest("GET", apiconfig.Server_URL+apiconfig.WORKFLOW_PATH+"?"+values.Encode(),
		web.WithPrefix(prefix),
		web.WithLog(false),
		web.WithBodyBytes(&bodyBytes))
	if err != nil {
		return err
	}
	var res []etcd.ListRes
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		log.Println(prefix, err)
		return err
	}
	fmt.Println(prefix, "Workflow Get successfully. Here are the results:")
	fmt.Println("total number:", len(res))

	table := uitable.New()
	table.MaxColWidth = 100
	table.RightAlign(10)
	table.AddRow("NAME", "RESULT")
	for _, val := range res {
		w := core.Workflow{}
		err := json.Unmarshal([]byte(val.Value), &w)
		if err != nil {
			log.Println(prefix, err)
			return err
		}
		fmt.Println("workflow:", w)
		table.AddRow(color.RedString(w.Name),
			color.BlueString(string(w.Spec.Result)))
	}
	fmt.Println(table)

	return nil
}

func (o *GetOptions) RunGetFunction(cmd *cobra.Command, args []string) error {
	prefix := "[kubectl] [get] [RunGetFunction] "
	values := url.Values{}
	if o.GetAll {
		values.Add("all", "true")
	}

	if o.WithNoPrefix {
		values.Add("prefix", "false")
	} else {
		values.Add("prefix", "true")
	}

	if len(args) > 1 {
		values.Add("Name", args[1])
		//body = bytes.NewBuffer([]byte(args[1]))
	}

	bodyBytes := make([]byte, 0)

	err := web.SendHttpRequest("GET", apiconfig.Server_URL+apiconfig.FUNCTION_PATH+"?"+values.Encode(),
		web.WithPrefix(prefix),
		web.WithLog(false),
		web.WithBodyBytes(&bodyBytes))
	if err != nil {
		return err
	}

	// 将字节数组转换为字符串并打印
	//var s GetRespond
	//json.Unmarshal(bodyBytes, &s)
	var res []etcd.ListRes
	json.Unmarshal(bodyBytes, &res)
	fmt.Println(prefix, "Pod Get successfully. Here are the results:")

	fmt.Println("total number:", len(res))

	table := uitable.New()
	table.MaxColWidth = 100
	table.RightAlign(10)
	table.AddRow("NAME", "InvokeTimes", "Image")
	for _, val := range res {
		function := core.Function{}
		err := json.Unmarshal([]byte(val.Value), &function)
		if err != nil {
			log.Println(prefix, err)
			return err
		}

		table.AddRow(color.RedString(function.Name),
			color.WhiteString(strconv.Itoa(function.Spec.InvokeTimes)),
			color.GreenString(function.Spec.Image))
	}
	fmt.Println(table)

	return nil
}
