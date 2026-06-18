package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/AsierCaballero/kubectl-usage/pkg/types"
)

type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

var colorMap = map[types.WasteLevel]tablewriter.Colors{
	types.WasteNone:    {tablewriter.Normal, tablewriter.FgGreenColor},
	types.WasteLow:     {tablewriter.Normal, tablewriter.FgGreenColor},
	types.WasteMedium:  {tablewriter.Normal, tablewriter.FgYellowColor},
	types.WasteHigh:    {tablewriter.Normal, tablewriter.FgRedColor},
	types.WasteExtreme: {tablewriter.Bold, tablewriter.FgHiRedColor},
}

func Print(results []types.AnalysisResult, format Format, sortBy string) error {
	results = sortResults(results, sortBy)

	switch format {
	case FormatJSON:
		return printJSON(results)
	case FormatYAML:
		return printYAML(results)
	default:
		return printTable(results)
	}
}

func printTable(results []types.AnalysisResult) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Pod", "Namespace", "CPU Req", "CPU Used", "CPU Waste",
		"Mem Req", "Mem Used", "Mem Waste", "Rec CPU", "Rec Mem",
	})

	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
	})

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, r := range results {
		colors := colorMap[r.OverallWaste]

		cpuReq := formatCPUMetric(r.Pod.CPURequestMillis())
		cpuUsed := formatCPUMetric(r.Pod.CPUUsageMillis())
		memReq := formatMemoryMetric(r.Pod.MemoryRequestBytes())
		memUsed := formatMemoryMetric(r.Pod.MemoryUsageBytes())

		row := []string{
			r.Pod.Name,
			shortNamespace(r.Pod.Namespace),
			cpuReq,
			cpuUsed,
			fmt.Sprintf("%.0f%%", r.CPUWastePercent),
			memReq,
			memUsed,
			fmt.Sprintf("%.0f%%", r.MemoryWastePercent),
			r.CPURecommendation,
			r.MemRecommendation,
		}

		table.Rich(row, []tablewriter.Colors{
			colors, colors, colors, colors, colors,
			colors, colors, colors, colors, colors,
		})
	}

	table.Render()
	return nil
}

func printJSON(results []types.AnalysisResult) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func printYAML(results []types.AnalysisResult) error {
	return fmt.Errorf("yaml output not implemented")
}

type Summary struct {
	TotalPods   int     `json:"total_pods"`
	TotalWaste  float64 `json:"total_waste_percent"`
	OK          int     `json:"ok"`
	NeedsReview int     `json:"needs_review"`
	Overprovisioned int `json:"overprovisioned"`
}

func formatCPUMetric(millis int64) string {
	if millis >= 1000 {
		return fmt.Sprintf("%.1f", float64(millis)/1000)
	}
	return fmt.Sprintf("%dm", millis)
}

func formatMemoryMetric(bytes int64) string {
	if bytes >= 1<<30 {
		return fmt.Sprintf("%.1fGi", float64(bytes)/float64(1<<30))
	}
	if bytes >= 1<<20 {
		return fmt.Sprintf("%.0fMi", float64(bytes)/float64(1<<20))
	}
	return fmt.Sprintf("%.0fKi", float64(bytes)/float64(1<<10))
}

func shortNamespace(ns string) string {
	parts := strings.SplitN(ns, "-", 2)
	if len(parts) == 2 {
		return parts[0][:1] + "-" + parts[1]
	}
	if len(ns) > 12 {
		return ns[:12]
	}
	return ns
}
