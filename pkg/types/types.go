package types

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type PodResource struct {
	Name      string
	Namespace string
	PodIP     string
	NodeName  string
	Phase     corev1.PodPhase

	CPURequest    *resource.Quantity
	CPULimit      *resource.Quantity
	CPUUsage      *resource.Quantity
	CPUUsageAvg   *resource.Quantity
	CPUUsageMax   *resource.Quantity

	MemoryRequest    *resource.Quantity
	MemoryLimit      *resource.Quantity
	MemoryUsage      *resource.Quantity
	MemoryUsageAvg   *resource.Quantity
	MemoryUsageMax   *resource.Quantity

	ContainerName string
	RestartCount  int32
	Age           string
}

type AnalysisResult struct {
	Pod               PodResource
	CPUWastePercent   float64
	MemoryWastePercent float64
	CPURecommendation string
	MemRecommendation string
	OverallWaste      WasteLevel
}

type WasteLevel string

const (
	WasteNone   WasteLevel = "none"
	WasteLow    WasteLevel = "low"
	WasteMedium WasteLevel = "medium"
	WasteHigh   WasteLevel = "high"
	WasteExtreme WasteLevel = "extreme"
)

func (p PodResource) CPURequestMillis() int64 {
	if p.CPURequest == nil {
		return 0
	}
	return p.CPURequest.MilliValue()
}

func (p PodResource) CPULimitMillis() int64 {
	if p.CPULimit == nil {
		return 0
	}
	return p.CPULimit.MilliValue()
}

func (p PodResource) CPUUsageMillis() int64 {
	if p.CPUUsage == nil {
		return 0
	}
	return p.CPUUsage.MilliValue()
}

func (p PodResource) MemoryRequestBytes() int64 {
	if p.MemoryRequest == nil {
		return 0
	}
	return p.MemoryRequest.Value()
}

func (p PodResource) MemoryUsageBytes() int64 {
	if p.MemoryUsage == nil {
		return 0
	}
	return p.MemoryUsage.Value()
}
