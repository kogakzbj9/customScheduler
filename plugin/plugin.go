package plugin

import (
	"context"
	"fmt"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"k8s.io/kubernetes/pkg/scheduler/framework/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
	"time"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"k8s.io/client-go/util/retry"
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

func (pl *CustomSchedulerPlugin) Permit(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) *framework.Status {
	// Get the CPU usage of the target node
	clientset, err := kubernetes.NewForConfig(rest.InClusterConfig())
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("Failed to create clientset: %v", err))
	}

	metricsClientset, err := versioned.NewForConfig(rest.InClusterConfig())
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("Failed to create metrics clientset: %v", err))
	}

	nodeMetrics, err := metricsClientset.MetricsV1beta1().NodeMetricses().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("Failed to get node metrics: %v", err))
	}

	allocatedCPU := nodeMetrics.Usage.Cpu().MilliValue()
	currentCPUUsage := nodeMetrics.Usage.Cpu().MilliValue()
	cpuUsagePercentage := (float64(currentCPUUsage) / float64(allocatedCPU)) * 100

	fmt.Printf("CPU usage of node %s: %d%%\n", nodeName, int(cpuUsagePercentage))

	config, err := clientset.CoreV1().ConfigMaps("default").Get(ctx, "custom-scheduler-config", metav1.GetOptions{})
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("Failed to get config map: %v", err))
	}

	cpuThreshold := 50
	waitTime := 10

	if val, ok := config.Data["cpuThreshold"]; ok {
		fmt.Sscanf(val, "%d", &cpuThreshold)
	}

	if val, ok := config.Data["waitTime"]; ok {
		fmt.Sscanf(val, "%d", &waitTime)
	}

	if cpuUsagePercentage > float64(cpuThreshold) {
		return framework.NewStatus(framework.Wait, "", time.Duration(waitTime)*time.Second)
	}

	return framework.NewStatus(framework.Success, "")
}

func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	return &CustomSchedulerPlugin{handle: h}, nil
}

func init() {
	framework.RegisterPlugin(Name, New)
}
