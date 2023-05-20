package jobcontroller

import (
	"encoding/json"
	"fmt"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	metav1 "minik8s/pkg/apis/meta/v1"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	q "minik8s/pkg/util/concurrentqueue"
	"time"
)

type JobController struct {
	JobInformer informer.Informer
	queue       *q.ConcurrentQueue
}

func NewJobController() *JobController {
	return &JobController{
		JobInformer: informer.NewInformer(apiconfig.JOB_PATH),
		queue:       q.NewConcurrentQueue(),
	}
}

func (jc *JobController) Register() {
	jc.JobInformer.AddEventHandler(tool.Added, jc.AddJob)
}

func (jc *JobController) AddJob(event tool.Event) {
	var job *core.Job
	err := json.Unmarshal([]byte(event.Val), &job)
	if err != nil {
		fmt.Println("[jobcontroller] add job error")
		return
	}
	jc.queue.Push(job)
}

func (jc *JobController) worker() {
	for {
		if !jc.queue.IsEmpty() {
			job := jc.queue.Pop()
			jc.RunJob(job.(*core.Job))
		} else {
			time.Sleep(time.Second)
		}

	}
}

func (jc *JobController) RunJob(job *core.Job) {
	cmd := fmt.Sprintf("./gpuserver --jobname=%s", job.Name)
	jobcontainer := core.Container{
		Name:  "gpu-job",
		Image: "gpu-job-image",
		VolumeMounts: []core.VolumeMount{
			{
				Name:      "job-volume",
				MountPath: "/home/job",
			},
		},
		Command: []string{cmd},
	}

	pod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: job.Name,
		},
		Spec: core.PodSpec{
			Volumes: []core.Volume{
				{
					Name:         "job-volume",
					VolumeSource: core.VolumeSource{},
				},
			},
			Containers: []core.Container{
				jobcontainer,
			},
			GPUJob:     true,
			GPUJobName: job.Name,
		},
	}

	tool.AddPod(pod)
}

func (jc *JobController) Run() {
	go jc.Register()
	go jc.worker()
	select {}
}
