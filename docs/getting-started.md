# Getting started

## Prerequisites

- Go 1.22+ (for building from source)
- kubectl configured with cluster access
- Metrics Server installed in the cluster

## Install Metrics Server

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

Verify:

```bash
kubectl top pods
```

## Install kubectl-usage

### From source

```bash
git clone https://github.com/AsierCaballero/kubectl-usage.git
cd kubectl-usage
make build
sudo make install
```

### Using Go install

```bash
go install github.com/AsierCaballero/kubectl-usage@latest
```

## Usage

```bash
# Analyze all pods in default namespace
kubectl-usage

# Analyze production namespace
kubectl-usage -n production

# Filter by deployment
kubectl-usage -n staging -d api-gateway

# Show all namespaces
kubectl-usage -A

# JSON output
kubectl-usage -n prod -o json

# Sort by CPU waste
kubectl-usage -n prod -s cpu

# Label selector
kubectl-usage -l app=myapp,tier=frontend
```

## Output

The table shows:
- **CPU/Mem Req**: configured requests
- **CPU/Mem Used**: actual usage from metrics-server
- **Waste**: over-provisioning percentage
- **Rec**: recommended request value (usage × 1.2)
- **OK**: no recommendation needed
