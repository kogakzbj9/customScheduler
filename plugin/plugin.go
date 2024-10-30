package plugin

import (
	"context"
	"fmt"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"k8s.io/kubernetes/pkg/scheduler/framework/runtime"
)

const (
	Name = "CustomSchedulerPlugin"
)

type CustomSchedulerPlugin struct {
	handle framework.Handle
}

func (pl *CustomSchedulerPlugin) Name() string {
	return Name
}

func (pl *CustomSchedulerPlugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	// Implement your custom filter logic here
	return framework.NewStatus(framework.Success, "")
}

func (pl *CustomSchedulerPlugin) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	// Implement your custom scoring logic here
	return 0, framework.NewStatus(framework.Success, "")
}

func (pl *CustomSchedulerPlugin) ScoreExtensions() framework.ScoreExtensions {
	return nil
}

func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	return &CustomSchedulerPlugin{handle: h}, nil
}

func init() {
	framework.RegisterPlugin(Name, New)
}
