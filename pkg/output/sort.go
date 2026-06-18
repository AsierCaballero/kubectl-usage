package output

import (
	"sort"

	"github.com/AsierCaballero/kubectl-usage/pkg/types"
)

type sortFunc func(a, b types.AnalysisResult) bool

var sorters = map[string]sortFunc{
	"name":      func(a, b types.AnalysisResult) bool { return a.Pod.Name < b.Pod.Name },
	"namespace": func(a, b types.AnalysisResult) bool { return a.Pod.Namespace < b.Pod.Namespace },
	"cpu":       func(a, b types.AnalysisResult) bool { return a.CPUWastePercent > b.CPUWastePercent },
	"memory":    func(a, b types.AnalysisResult) bool { return a.MemoryWastePercent > b.MemoryWastePercent },
	"waste":     func(a, b types.AnalysisResult) bool { return a.OverallWaste > b.OverallWaste },
}

func sortResults(results []types.AnalysisResult, sortBy string) []types.AnalysisResult {
	sorter, ok := sorters[sortBy]
	if !ok {
		sorter = sorters["waste"]
	}

	sorted := make([]types.AnalysisResult, len(results))
	copy(sorted, results)

	sort.SliceStable(sorted, func(i, j int) bool {
		return sorter(sorted[i], sorted[j])
	})

	return sorted
}
