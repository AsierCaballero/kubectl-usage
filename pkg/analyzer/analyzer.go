package analyzer

import (
	"fmt"

	"github.com/AsierCaballero/kubectl-usage/pkg/types"
)

type Analyzer struct{}

func New() *Analyzer {
	return &Analyzer{}
}

type Options struct {
	CPUThresholdPercent   float64
	MemoryThresholdPercent float64
}

var DefaultOptions = Options{
	CPUThresholdPercent:   50,
	MemoryThresholdPercent: 50,
}

func (a *Analyzer) Analyze(pods []types.PodResource, opts Options) []types.AnalysisResult {
	if opts.CPUThresholdPercent == 0 {
		opts.CPUThresholdPercent = DefaultOptions.CPUThresholdPercent
	}
	if opts.MemoryThresholdPercent == 0 {
		opts.MemoryThresholdPercent = DefaultOptions.MemoryThresholdPercent
	}

	results := make([]types.AnalysisResult, 0, len(pods))

	for _, pod := range pods {
		result := types.AnalysisResult{Pod: pod}
		result.CPUWastePercent = calculateWastePercent(
			pod.CPURequestMillis(),
			pod.CPUUsageMillis(),
		)
		result.MemoryWastePercent = calculateWastePercent(
			pod.MemoryRequestBytes(),
			pod.MemoryUsageBytes(),
		)

		result.CPURecommendation = recommendResource(
			pod.CPURequestMillis(),
			pod.CPUUsageMillis(),
			opts.CPUThresholdPercent,
			"cpu",
		)
		result.MemRecommendation = recommendResource(
			pod.MemoryRequestBytes(),
			pod.MemoryUsageBytes(),
			opts.MemoryThresholdPercent,
			"memory",
		)

		result.OverallWaste = overallWasteLevel(result.CPUWastePercent, result.MemoryWastePercent)

		results = append(results, result)
	}

	return results
}

func calculateWastePercent(request, usage int64) float64 {
	if request <= 0 {
		return 0
	}
	if usage <= 0 {
		return 100
	}
	waste := float64(request-usage) / float64(request) * 100
	if waste < 0 {
		return 0
	}
	return waste
}

func recommendResource(request, usage int64, thresholdPercent float64, resourceType string) string {
	if request <= 0 {
		return "no request set"
	}
	if usage <= 0 {
		if resourceType == "cpu" {
			return "remove request or set to 0"
		}
		return "remove request or set to 0"
	}

	wastePercent := calculateWastePercent(request, usage)

	if wastePercent >= thresholdPercent {
		recommended := float64(usage) * 1.2
		switch resourceType {
		case "cpu":
			return formatCPUMillis(int64(recommended))
		case "memory":
			return formatMemoryBytes(int64(recommended))
		}
	}

	return "OK"
}

func formatCPUMillis(millis int64) string {
	if millis >= 1000 {
		return fmt.Sprintf("%.1f", float64(millis)/1000)
	}
	return fmt.Sprintf("%dm", millis)
}

func formatMemoryBytes(bytes int64) string {
	if bytes >= 1<<30 {
		return fmt.Sprintf("%.1fGi", float64(bytes)/float64(1<<30))
	}
	if bytes >= 1<<20 {
		return fmt.Sprintf("%.0fMi", float64(bytes)/float64(1<<20))
	}
	return fmt.Sprintf("%.0fKi", float64(bytes)/float64(1<<10))
}

func overallWasteLevel(cpuWaste, memWaste float64) types.WasteLevel {
	maxWaste := cpuWaste
	if memWaste > maxWaste {
		maxWaste = memWaste
	}
	switch {
	case maxWaste >= 80:
		return types.WasteExtreme
	case maxWaste >= 60:
		return types.WasteHigh
	case maxWaste >= 30:
		return types.WasteMedium
	case maxWaste > 0:
		return types.WasteLow
	default:
		return types.WasteNone
	}
}
