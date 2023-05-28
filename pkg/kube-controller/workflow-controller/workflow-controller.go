package workflowcontroller

import (
	"encoding/json"
	"fmt"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	q "minik8s/pkg/util/concurrentqueue"
	"minik8s/pkg/util/web"
	"time"
)

type WorkflowController struct {
	WorkflowInformer informer.Informer
	queue            *q.ConcurrentQueue
	mark             map[string]bool
}

func NewHPAController() *WorkflowController {
	return &WorkflowController{
		WorkflowInformer: informer.NewInformer(apiconfig.WORKFLOW_PATH),
		queue:            q.NewConcurrentQueue(),
		mark:             make(map[string]bool),
	}
}

func (wfc *WorkflowController) Register() {
	wfc.WorkflowInformer.AddEventHandler(tool.Added, wfc.AddWorkflow)
	wfc.WorkflowInformer.AddEventHandler(tool.Deleted, wfc.DeleteWorkflow)
}

func (wfc *WorkflowController) Run() {
	go wfc.WorkflowInformer.Run()
	go worker(wfc)
	select {}
}

func worker(wfc *WorkflowController) {
	prefix := "[workflow] "
	for {
		if !wfc.queue.IsEmpty() {
			key := wfc.queue.Pop().(string)
			val, exist := (*wfc.WorkflowInformer.GetCache())[key]
			if !exist {
				fmt.Println("[ERROR] ", prefix, "cache doesn't have key:", key)
				continue
			}
			dag := &core.DAG{}
			err := json.Unmarshal([]byte(val), dag)
			if err != nil {
				fmt.Println("[ERROR] ", prefix, err)
				return
			}
			go DoWorkflowDAG(dag)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func DoWorkflowDAG(dag *core.DAG) {
	fmt.Println("DoWorkflow")
	startnode := dag.StartNode
	curnode := startnode.OutEdges[0].To
	var result string
	var err error
	for curnode.Type != core.StateTypeEnd {
		fmt.Println(curnode)
		result, err = TriggerFunc(curnode.Function, result)
		if err != nil {
			fmt.Println("[workflow] running func:", err)
			return
		}
		success := false
		for _, edge := range curnode.OutEdges {
			if edge.Condition == "true" {
				curnode = edge.To
				success = true
				break
			} else if evalCondition(edge.Condition, result) {
				curnode = edge.To
				success = true
				break
			}
		}
		if !success {
			fmt.Println("[workflow] no edge can be triggered")
			return
		}
	}
}

// TO DO when func end
func evalCondition(condition string, result string) bool {
	return true
}

// TO DO when func end
func TriggerFunc(f core.TMPfunction, input string) (string, error) {
	url := "http://localhost:10000/" // trigger function url
	bodyBytes := make([]byte, 0)
	web.SendHttpRequest("PUT", url, web.WithBodyBytes(&bodyBytes))
	return "result", nil
}

func (wfc *WorkflowController) AddWorkflow(event tool.Event) {
	prefix := "[Workflow] [AddWorkflow] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	wfc.WorkflowInformer.Set(event.Key, event.Val)
	wfc.queue.Push(event.Key)
}

func (wfc *WorkflowController) DeleteWorkflow(event tool.Event) {
	prefix := "[Workflow] [DeleteWorkflow] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	wfc.WorkflowInformer.Delete(event.Key)
}
