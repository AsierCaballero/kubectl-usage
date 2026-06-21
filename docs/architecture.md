# Architecture

## Overview

kubectl-usage is a CLI tool that connects to the Kubernetes API server, fetches pod specifications and resource metrics, and analyzes resource utilization.

## Components

### Collector

The collector uses `client-go` to query:
- Core API: Pod list with container resource requests/limits
- Metrics API: Pod CPU and memory usage

### Analyzer

The analyzer calculates for each container:
- Waste percentage: `(request - usage) / request * 100`
- Recommendation: `usage * 1.2` if waste exceeds threshold
- Overall waste level: based on max waste percentage

### Output

Formats analysis results:
- Table: color-coded terminal output via tablewriter
- JSON: structured output for programmatic use

## Data flow

```
User CLI → Cobra command → Collector (K8s API)
                                ↓
                         Analyzer (waste calc)
                                ↓
                         Output (table/JSON)
```
