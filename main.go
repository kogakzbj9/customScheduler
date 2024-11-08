package main

import (
	"os"

	"github.com/kogakzbj9/customScheduler/plugin"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
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
