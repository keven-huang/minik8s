package main

import (
	"flag"
	"minik8s/pkg/scheduler"
)

func main() {
	Strategy := flag.String("strategy", scheduler.RRStrategy, "RRStrategy or RandomStrategy")
	flag.Parse()
	schedulerNow := scheduler.NewScheduler(Strategy)
	schedulerNow.Register()
	schedulerNow.Run()
	select {}
}
