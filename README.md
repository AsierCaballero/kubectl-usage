# kubectl-usage

[![CI](https://img.shields.io/github/actions/workflow/status/AsierCaballero/kubectl-usage/ci.yml?label=CI&logo=github)](https://github.com/AsierCaballero/kubectl-usage/actions)
[![Go](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go)](https://go.dev)
[![Kubernetes](https://img.shields.io/badge/kubernetes-1.27+-326CE5?logo=kubernetes)](https://kubernetes.io)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A kubectl plugin that analyzes pod resource usage against configured requests and limits. Detects over-provisioned resources and provides downsizing recommendations.

```bash
kubectl usage -n prod
POD                  CPU_REQ  CPU_USED  WASTE  RECOMMEND
api-7d8f9c-2k3p      500m     120m      76%    250m
worker-6b1a2-9x4z    1        980m      2%     OK
```

---

## Features

- **Resource analysis** — compares CPU/memory requests against actual usage
- **Waste detection** — calculates over-provisioning percentage per pod
- **Recommendations** — suggests right-sized requests based on real usage
- **Flat 1.2x headroom** — recommended value = usage × 1.2
- **Deployment filter** — analyze pods belonging to a specific Deployment
- **Sort modes** — sort by waste, CPU, memory, name, or namespace
- **Output formats** — table (default) or JSON
- **Color-coded** — green/yellow/red based on waste severity
- **Kubectl plugin** — works with or without `kubectl plugin` mechanism

## Quick start

### Install

```bash
# As a kubectl plugin
kubectl krew install usage
# or manual
go install github.com/AsierCaballero/kubectl-usage@latest
```

### Usage

```bash
# Analyze default namespace
kubectl usage

# Analyze a specific namespace
kubectl usage -n prod

# Filter by deployment
kubectl usage -n staging -d api-gateway

# All namespaces, sorted by waste
kubectl usage -A -s waste

# JSON output
kubectl usage -n prod -o json

# Label selector
kubectl usage -l app=myapp,tier=frontend
```

## Example output

```
$ kubectl usage -n production -s waste
Pod                          Namespace  CPU Req  CPU Used  CPU Waste  Mem Req  Mem Used  Mem Waste  Rec CPU  Rec Mem
redis-5c3e1-7h8j             prod       250m     15m       94%        256Mi    80Mi      69%        100m     128Mi
api-7d8f9c-2k3p              prod       500m     120m      76%        512Mi    200Mi     61%        250m     256Mi
worker-6b1a2-9x4z            prod       1        980m      2%         1Gi      900Mi     10%        OK       OK
```

## Architecture

```
kubectl usage
     │
     ▼
  ┌─────────┐     ┌──────────────┐     ┌───────────┐
  │Collector │────▶│   Analyzer   │────▶│  Output   │
  │(K8s API) │     │ (waste calc) │     │ (table)   │
  └─────────┘     └──────────────┘     └───────────┘
```

## Author

**Asier Caballero** — Senior DevOps Engineer & Cloud Architect
asier.caballero1@gmail.com · [linkedin.com/in/asier-caballero](https://linkedin.com/in/asier-caballero)

