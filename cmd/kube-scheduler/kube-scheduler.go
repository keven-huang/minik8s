package main

import "minik8s/pkg/scheduler"

func main() {
	scheduler := scheduler.NewScheduler()
	scheduler.Register()
	scheduler.Run()
	select {}
}
