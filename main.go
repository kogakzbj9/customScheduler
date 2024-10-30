package main

import (
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"k8s.io/kubernetes/pkg/scheduler/framework/runtime"
	"k8s.io/klog/v2"
	"os"
	"github.com/kogakzbj9/customScheduler/plugin"
)

func main() {
	command := app.NewSchedulerCommand(
		app.WithPlugin(plugin.Name, plugin.New),
	)

	if err := command.Execute(); err != nil {
		klog.Fatalf("Error running scheduler: %v", err)
		os.Exit(1)
	}
}
