package plugin

import (
	"context"
	"fmt"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
	"time"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"strconv"
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
	config, err := rest.InClusterConfig()
	if err != nil {
    	// エラーハンドリング
		return framework.NewStatus(framework.Error, fmt.Sprintf("Failed creating in-cluster config: %v", err))
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("Failed to create clientset: %v", err))
	}

	metricsClientset, err := versioned.NewForConfig(config)
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

	customSchedulerConfig, err := clientset.CoreV1().ConfigMaps("kube-system").Get(ctx, "custom-scheduler-config", metav1.GetOptions{})
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("Failed to get config map: %v", err))
	}

	cpuThreshold := 50
	waitTime := 10

	if val, ok := customSchedulerConfig.Data["cpuThreshold"]; ok {
		fmt.Sscanf(val, "%d", &cpuThreshold)
	}

	if val, ok := customSchedulerConfig.Data["waitTime"]; ok {
		fmt.Sscanf(val, "%d", &waitTime)
	}

	// Retrieve the cpuSpike annotation from the Pod
	cpuSpikeAnnotation := pod.Annotations["cpuSpike"]
	cpuSpikeValue := 0
	if cpuSpikeAnnotation != "" {
		cpuSpikeValue, err = strconv.Atoi(cpuSpikeAnnotation)
		if err != nil {
			return framework.NewStatus(framework.Error, fmt.Sprintf("Invalid cpuSpike annotation value: %v", err))
		}
	}

	if cpuUsagePercentage+float64(cpuSpikeValue) > float64(cpuThreshold) {
		return framework.NewStatus(framework.Wait, fmt.Sprintf("%v", time.Duration(waitTime)*time.Second))
	}

	return framework.NewStatus(framework.Success, "")
}

func New(ctx context.Context, configuration runtime.Object, h framework.Handle) (framework.Plugin, error)  {
	return &CustomSchedulerPlugin{handle: h}, nil
}

// func init() {
// 	framework.RegisterPlugin(Name, New)
// }
