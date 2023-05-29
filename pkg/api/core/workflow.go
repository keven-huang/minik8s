package core

import (
	"fmt"
	metav1 "minik8s/pkg/apis/meta/v1"
)

// workflow
type Workflow struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec              WorkflowSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type WorkflowSpec struct {
	States []State `json:"states,omitempty" yaml:"states,omitempty"`
	Result string  `json:"result,omitempty" yaml:"result,omitempty"`
}

type State struct {
	Name      string    `json:"Name" yaml:"Name"`
	Type      StateType `json:"Type" yaml:"Type"`
	Choices   []Choice  `json:"Choices,omitempty" yaml:"Choices,omitempty"`
	Resource  string    `json:"Resource,omitempty" yaml:"Resource,omitempty"`
	Next      string    `json:"Next,omitempty" yaml:"Next,omitempty"`
	End       bool      `json:"End,omitempty" yaml:"End,omitempty"`
	InputData string    `json:"InputData,omitempty" yaml:"InputData,omitempty"`
}

type StateType string

const (
	StateTypeTask   StateType = "Task"
	StateTypeChoice StateType = "Choice"
	StateTypeInput  StateType = "Input"
	StateTypeEnd    StateType = "End"
)

type Choice struct {
	Condition string `json:"Condition" yaml:"Condition"`
	Next      string `json:"Next" yaml:"Next"`
}

// DAG

type DAG struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	StartNode         DAGNode
	Nodes             []DAGNode `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	Edges             []DAGEdge `json:"edges,omitempty" yaml:"edges,omitempty"`
}

// 占位符
type TMPfunction struct {
	Name string `json:"name" yaml:"name"`
}

type DAGNode struct {
	Type     StateType `json:"type" yaml:"type"`
	Function TMPfunction
	OutEdges []DAGEdge
}

type DAGEdge struct {
	From      DAGNode `json:"from" yaml:"from"`
	To        DAGNode `json:"to" yaml:"to"`
	Condition string  `json:"condition" yaml:"condition"`
}

func GetFunc(resource string) TMPfunction {
	return TMPfunction{
		Name: resource,
	}
}

func (w *Workflow) Workflow2DAG() (*DAG, error) {
	var dag DAG
	mapState := make(map[string]State)
	mapNode := make(map[string]*DAGNode)
	for _, state := range w.Spec.States {
		var node *DAGNode
		_, ok := mapState[state.Name]
		if ok {
			return nil, fmt.Errorf("state name %s is duplicate", state.Name)
		} else {
			mapState[state.Name] = state
		}
		if state.Type == StateTypeInput {
			dag.StartNode = DAGNode{
				Type: StateTypeInput,
			}
			node = &dag.StartNode
		}
		if state.Type == StateTypeTask {
			nodeType := StateTypeTask
			if state.End {
				nodeType = StateTypeEnd
			}
			node = &DAGNode{
				Type:     nodeType,
				Function: GetFunc(state.Resource),
			}
		}
		if state.Type == StateTypeChoice {
			node = &DAGNode{
				Type: StateTypeChoice,
			}
		}

		mapNode[state.Name] = node
	}
	for _, state := range w.Spec.States {
		curNode := mapNode[state.Name]
		if state.Type != StateTypeChoice {
			if state.Next != "" {
				if state.End {
					return nil, fmt.Errorf("state name %s is end state, next is not allowed", state.Name)
				}
				nextNode, ok := mapNode[state.Next]
				if ok {
					edge := DAGEdge{
						From:      *curNode,
						To:        *nextNode,
						Condition: "true",
					}
					dag.Edges = append(dag.Edges, edge)
					curNode.OutEdges = append(curNode.OutEdges, edge)
				} else {
					return nil, fmt.Errorf("state name %s is not exist", state.Next)
				}
			}
		} else {
			for _, choice := range state.Choices {
				nextNode, ok := mapNode[choice.Next]
				if ok {
					edge := DAGEdge{
						From:      *curNode,
						To:        *nextNode,
						Condition: choice.Condition,
					}
					dag.Edges = append(dag.Edges, edge)
					curNode.OutEdges = append(curNode.OutEdges, edge)
				} else {
					return nil, fmt.Errorf("state name %s is not exist", choice.Next)
				}
			}
		}
	}
	for _, node := range mapNode {
		dag.Nodes = append(dag.Nodes, *node)
	}
	return &dag, nil
}
