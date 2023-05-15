package core

import (
	metav1 "minik8s/pkg/apis/meta/v1"
)

type Job struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	// JobConfig JobConfig `json:`
	Spec JobSpec `json:"spec" yaml:"spec"`
}

type JobSpec struct {
	Jobs []JobConfig `json:"jobs" yaml:"jobs"`
}

type JobConfig struct {
	JobName         string `json:"name" yaml:"name"`
	Partition       string `json:"partition" yaml:"partition"`
	Nodes           string `json:"nodes" yaml:"nodes"`
	NTasks_Per_Node int    `json:"ntasks-per-node" yaml:"ntasks-per-node"`
	Cpus_Per_Task   int    `json:"cpus-per-task" yaml:"cpus-per-task"`
	Gpu             int    `json:"gpu" yaml:"gpu"`
	MailType        string `json:"mail-type" yaml:"mail-type"`
	MailUser        string `json:"mail-user" yaml:"mail-user"`
	Program         string `json:"program" yaml:"program"`
}

// JobStatus defines the observed state of Job
type JobStatus struct {
	JobName string
	Status  string
}
