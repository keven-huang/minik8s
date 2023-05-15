package core

import (
	"fmt"
	metav1 "minik8s/pkg/apis/meta/v1"
	"strings"
)

// user visible types
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
	Nodes           int    `json:"nodes" yaml:"nodes"`
	NTasks          int    `json:"ntasks" yaml:"ntasks"`
	NTasks_Per_Node int    `json:"ntasks-per-node" yaml:"ntasks-per-node"`
	Cpus_Per_Task   int    `json:"cpus-per-task" yaml:"cpus-per-task"`
	Gpu             int    `json:"gpu" yaml:"gpu"`
	Error           string `json:"error" yaml:"error"`
	Output          string `json:"output" yaml:"output"`
	MailType        string `json:"mail-type" yaml:"mail-type"`
	MailUser        string `json:"mail-user" yaml:"mail-user"`
	Program         string `json:"program" yaml:"program"`
}

// system visible types
type JobUpload struct {
	JobName string `json:"jobname"`
	Slurm   []byte `json:"slurm"`
	Program []byte `json:"program"`
}

// JobStatus defines the observed state of Job
type JobStatus struct {
	JobName string
	Status  string
}

func (jc *JobConfig) GenerateSlurm() []byte {
	var slurm []string
	slurm = append(slurm, "#!/bin/bash")
	slurm = append(slurm, fmt.Sprintf("#SBATCH --job-name=%s", jc.JobName))
	slurm = append(slurm, fmt.Sprintf("#SBATCH --partition=%s", jc.Partition))
	if jc.Cpus_Per_Task > 0 {
		slurm = append(slurm, fmt.Sprintf("#SBATCH --cpus-per-task=%d", jc.Cpus_Per_Task))
	}
	if jc.Nodes > 0 {
		slurm = append(slurm, fmt.Sprintf("#SBATCH --nodes=%d", jc.Nodes))
	}
	if jc.NTasks > 0 {
		slurm = append(slurm, fmt.Sprintf("#SBATCH -n %d", jc.NTasks))
	}
	if jc.NTasks_Per_Node > 0 {
		slurm = append(slurm, fmt.Sprintf("#SBATCH --ntasks-per-node=%d", jc.NTasks_Per_Node))
	}
	if jc.Output != "" {
		slurm = append(slurm, fmt.Sprintf("#SBATCH --output=%s", jc.Output))
	}
	if jc.Error != "" {
		slurm = append(slurm, fmt.Sprintf("#SBATCH --error=%s", jc.Error))
	}
	if jc.MailType != "" {
		slurm = append(slurm, fmt.Sprintf("#SBATCH --mail-type=%s", jc.MailType))
	}
	if jc.MailUser != "" {
		slurm = append(slurm, fmt.Sprintf("#SBATCH --mail-user=%s", jc.MailUser))
	}
	slurm = append(slurm, fmt.Sprintf("#SBATCH --gres=%d", jc.Gpu))

	slurm = append(slurm, "module load gcc/8.3.0 cuda/10.1.243-gcc-8.3.0")
	slurm = append(slurm, fmt.Sprintf("nvcc %s -o job -lcublas", jc.Program))
	slurm = append(slurm, "./job")
	return []byte(strings.Join(slurm, "\n"))
}
