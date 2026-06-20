package output

import (
	"os"
	"strings"
	"testing"

	"github.com/AsierCaballero/kubectl-usage/pkg/types"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestSortResults(t *testing.T) {
	results := []types.AnalysisResult{
		{Pod: types.PodResource{Name: "b-pod"}, CPUWastePercent: 50},
		{Pod: types.PodResource{Name: "a-pod"}, CPUWastePercent: 80},
	}

	sorted := sortResults(results, "cpu")
	if sorted[0].Pod.Name != "a-pod" {
		t.Errorf("expected a-pod first by cpu waste, got %s", sorted[0].Pod.Name)
	}

	sorted = sortResults(results, "name")
	if sorted[0].Pod.Name != "a-pod" {
		t.Errorf("expected a-pod first by name, got %s", sorted[0].Pod.Name)
	}
}

func TestFormatCPUMetric(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{100, "100m"},
		{1000, "1.0"},
		{1500, "1.5"},
	}

	for _, tt := range tests {
		got := formatCPUMetric(tt.input)
		if got != tt.expected {
			t.Errorf("formatCPUMetric(%d) = %s; want %s", tt.input, got, tt.expected)
		}
	}
}

func TestFormatMemoryMetric(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{1024, "1.0Ki"},
		{1048576, "1.0Mi"},
		{1073741824, "1.0Gi"},
	}

	for _, tt := range tests {
		got := formatMemoryMetric(tt.input)
		if got != tt.expected {
			t.Errorf("formatMemoryMetric(%d) = %s; want %s", tt.input, got, tt.expected)
		}
	}
}

func TestPrintJSON(t *testing.T) {
	results := []types.AnalysisResult{
		{
			Pod: types.PodResource{
				Name:      "test-pod",
				Namespace: "default",
				CPURequest: resource.NewMilliQuantity(500, resource.DecimalSI),
			},
			CPUWastePercent: 50,
		},
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := printJSON(results)
	if err != nil {
		t.Fatalf("printJSON: %v", err)
	}

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	var readBuf [4096]byte
	n, _ := r.Read(readBuf[:])
	buf.Write(readBuf[:n])

	if !strings.Contains(buf.String(), "test-pod") {
		t.Error("JSON output should contain pod name")
	}
}

func TestShortNamespace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"production", "production"},
		{"my-long-namespace", "m-long-namespace"},
	}

	for _, tt := range tests {
		got := shortNamespace(tt.input)
		if got != tt.expected {
			t.Errorf("shortNamespace(%s) = %s; want %s", tt.input, got, tt.expected)
		}
	}
}
