# Installation

This guide explains how to install and set up the Security Log Event Processor.

## Prerequisites

- OpenTelemetry Collector (v0.138.0 or later)
- Go 1.21+ (if building from source)
- Docker (for containerized deployment)
- Kubernetes cluster (for Kubernetes deployment)

## Installation Methods

### Option 1: Using Pre-built Collector (Recommended)

Build the collector using OpenTelemetry Collector Builder (OCB) with the included manifest:

```bash
# Clone the repository
git clone https://github.com/dynatrace/securitylogeventprocessor.git
cd securitylogeventprocessor

# Build using OCB
make ocb-build
```

The collector binary will be available at `dist/otelcol-securityevent`.

### Option 2: Docker Image

Build the Docker image:

```bash
# Build Docker image
make docker-build

# Or using docker-compose
docker-compose -f docker/docker-compose.yml build
```

### Option 3: Kubernetes Deployment

Deploy to Kubernetes:

```bash
# Apply Kubernetes manifests
kubectl apply -f k8s/

# Or using kustomize
kubectl apply -k k8s/
```

## Verify Installation

### Verify Binary

```bash
./dist/otelcol-securityevent --version
```

### Verify Docker Image

```bash
docker run --rm otelcol-securityevent:latest --version
```

### Verify Kubernetes Deployment

```bash
kubectl get pods -n otelcol-securityevent
kubectl logs -n otelcol-securityevent -l app.kubernetes.io/name=otelcol-securityevent
```

## Next Steps

- [Quick Start Guide](quick-start.md): Get started with a simple configuration
- [Configuration Guide](../configuration/processor-config.md): Learn about configuration options

