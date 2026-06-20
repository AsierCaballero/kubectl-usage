package analyzer

import (
	"testing"

	"github.com/AsierCaballero/kubectl-usage/pkg/types"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestCalculateWastePercent(t *testing.T) {
	tests := []struct {
		name     string
		request  int64
		usage    int64
		expected float64
	}{
		{"no request", 0, 100, 0},
		{"no usage", 1000, 0, 100},
		{"full usage", 1000, 1000, 0},
		{"half usage", 1000, 500, 50},
		{"over request", 500, 1000, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateWastePercent(tt.request, tt.usage)
			if got != tt.expected {
				t.Errorf("calculateWastePercent(%d, %d) = %.0f; want %.0f",
					tt.request, tt.usage, got, tt.expected)
			}
		})
	}
}

func TestRecommendResource(t *testing.T) {
	tests := []struct {
		name            string
		request         int64
		usage           int64
		threshold       float64
		resourceType    string
		expectRecommend bool
	}{
		{"no request", 0, 100, 50, "cpu", false},
		{"no usage", 1000, 0, 50, "cpu", false},
		{"efficient cpu", 1000, 900, 50, "cpu", true},
		{"wasteful cpu", 1000, 100, 50, "cpu", false},
		{"wasteful memory", 1073741824, 104857600, 50, "memory", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := recommendResource(tt.request, tt.usage, tt.threshold, tt.resourceType)
			if tt.expectRecommend && got == "OK" {
				t.Errorf("expected recommendation, got OK")
			}
			if !tt.expectRecommend && got != "OK" {
				t.Errorf("expected OK, got %s", got)
			}
		})
	}
}

func TestOverallWasteLevel(t *testing.T) {
	tests := []struct {
		cpu    float64
		mem    float64
		want   types.WasteLevel
	}{
		{0, 0, types.WasteNone},
		{10, 20, types.WasteLow},
		{40, 10, types.WasteMedium},
		{70, 20, types.WasteHigh},
		{90, 0, types.WasteExtreme},
	}

	for _, tt := range tests {
		got := overallWasteLevel(tt.cpu, tt.mem)
		if got != tt.want {
			t.Errorf("overallWasteLevel(%.0f, %.0f) = %s; want %s",
				tt.cpu, tt.mem, got, tt.want)
		}
	}
}

func TestAnalyze(t *testing.T) {
	a := New()
	pods := []types.PodResource{
		{
			Name:       "test-pod",
			CPURequest: resource.NewMilliQuantity(1000, resource.DecimalSI),
			CPUUsage:   resource.NewMilliQuantity(200, resource.DecimalSI),
		},
	}

	results := a.Analyze(pods, Options{CPUThresholdPercent: 50, MemoryThresholdPercent: 50})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].CPUWastePercent != 80 {
		t.Errorf("expected 80%% CPU waste, got %.0f%%", results[0].CPUWastePercent)
	}
}
