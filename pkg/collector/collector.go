package collector

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/AsierCaballero/kubectl-usage/pkg/types"
)

type Collector struct {
	clientset    kubernetes.Interface
	metricsClient metricsclient.Interface
}

func New(clientset kubernetes.Interface, metricsClient metricsclient.Interface) *Collector {
	return &Collector{
		clientset:    clientset,
		metricsClient: metricsClient,
	}
}

type Options struct {
	Namespace     string
	LabelSelector string
	Deployment    string
	AllNamespaces bool
}

func (c *Collector) Collect(ctx context.Context, opts Options) ([]types.PodResource, error) {
	pods, err := c.listPods(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("list pods: %w", err)
	}

	usage, err := c.fetchPodMetrics(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("fetch metrics: %w", err)
	}

	usageMap := buildUsageMap(usage)

	resources := make([]types.PodResource, 0, len(pods))
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			r := types.PodResource{
				Name:         pod.Name,
				Namespace:    pod.Namespace,
				PodIP:        pod.Status.PodIP,
				NodeName:     pod.Spec.NodeName,
				Phase:        pod.Status.Phase,
				CPURequest:   container.Resources.Requests.Cpu(),
				CPULimit:     container.Resources.Limits.Cpu(),
				MemoryRequest: container.Resources.Requests.Memory(),
				MemoryLimit:  container.Resources.Limits.Memory(),
				ContainerName: container.Name,
				RestartCount: getRestartCount(pod),
				Age:          formatAge(pod.CreationTimestamp.Time),
			}

			if metrics, ok := usageMap[pod.Namespace+"/"+pod.Name]; ok {
				for _, mc := range metrics.Containers {
					if mc.Name == container.Name {
						r.CPUUsage = mc.Usage.Cpu()
						r.MemoryUsage = mc.Usage.Memory()
						break
					}
				}
			}

			resources = append(resources, r)
		}
	}

	return resources, nil
}

func (c *Collector) listPods(ctx context.Context, opts Options) ([]corev1.Pod, error) {
	var listOpts metav1.ListOptions{}
	if opts.LabelSelector != "" {
		listOpts.LabelSelector = opts.LabelSelector
	}

	namespace := opts.Namespace
	if opts.AllNamespaces {
		namespace = ""
	}

	podList, err := c.clientset.CoreV1().Pods(namespace).List(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	if opts.Deployment != "" {
		sel, err := c.deploymentSelector(ctx, opts.Deployment, opts.Namespace)
		if err != nil {
			return nil, err
		}
		if sel != nil {
			filtered := make([]corev1.Pod, 0, len(podList.Items))
			for _, pod := range podList.Items {
				if labels.SelectorFromSet(sel).Matches(labels.Set(pod.Labels)) {
					filtered = append(filtered, pod)
				}
			}
			return filtered, nil
		}
	}

	return podList.Items, nil
}

func (c *Collector) deploymentSelector(ctx context.Context, name, namespace string) (map[string]string, error) {
	deploy, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get deployment %s: %w", name, err)
	}
	return deploy.Spec.Selector.MatchLabels, nil
}

func (c *Collector) fetchPodMetrics(ctx context.Context, opts Options) ([]v1beta1.PodMetrics, error) {
	namespace := opts.Namespace
	if opts.AllNamespaces {
		namespace = ""
	}

	listOpts := metav1.ListOptions{}
	if opts.LabelSelector != "" {
		listOpts.LabelSelector = opts.LabelSelector
	}

	metricsList, err := c.metricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, listOpts)
	if err != nil {
		return nil, err
	}
	return metricsList.Items, nil
}

func buildUsageMap(metrics []v1beta1.PodMetrics) map[string]v1beta1.PodMetrics {
	m := make(map[string]v1beta1.PodMetrics, len(metrics))
	for _, pm := range metrics {
		m[pm.Namespace+"/"+pm.Name] = pm
	}
	return m
}

func getRestartCount(pod corev1.Pod) int32 {
	var total int32
	for _, cs := range pod.Status.ContainerStatuses {
		total += cs.RestartCount
	}
	return total
}

func formatAge(timestamp metav1.Time) string {
	// Handled in output formatting using PodResource directly
	_ = timestamp
	return ""
}
