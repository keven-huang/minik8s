package scheduler

import (
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/client/informer"
)

type Scheduler struct {
	PodInformer informer.Informer
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		PodInformer: informer.NewInformer(apiconfig.POD_PATH),
	}
}

func (s *Scheduler) Run() {

}
