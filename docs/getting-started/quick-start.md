# Quick Start

Get the Security Log Event Processor up and running in 5 minutes!

## Step 1: Build the Collector

```bash
# Build the collector
make ocb-build
```

Or build the Docker image:

```bash
make docker-build
```

## Step 2: Create Configuration File

Create a file `config.yaml`:

```yaml
receivers:
  # k8sobjects receiver for collecting OpenReports Custom Resources
  k8sobjects:
    auth_type: serviceAccount
    objects:
      - name: openreports
        mode: watch
        api_group: openreports.io
        api_version: v1alpha1
        resource: reports

processors:
  batch:
    timeout: 1s
    send_batch_size: 1024
  
  memory_limiter:
    check_interval: 1s
    limit_mib: 4000
  
  # Security Event Processor
  securityevent:
    processors:
      openreports:
        enabled: true
        status_filter:
          - "fail"
          - "error"

exporters:
  logging:
    loglevel: info

service:
  pipelines:
    logs:
      receivers: [k8sobjects]
      processors: [memory_limiter, securityevent, batch]
      exporters: [logging]
  
  extensions: [health_check]

extensions:
  health_check:
    endpoint: 0.0.0.0:13133
```

## Step 3: Run the Collector

### Using Binary

```bash
./dist/otelcol-securityevent --config=config.yaml
```

### Using Docker

```bash
docker run -d \
  --name otelcol-securityevent \
  -p 13133:13133 \
  -v $(pwd)/config.yaml:/etc/otelcol/config.yaml:ro \
  otelcol-securityevent:latest
```

### Using Docker Compose

```bash
docker-compose -f docker/docker-compose.yml up -d
```

### Using Kubernetes

```bash
# Apply manifests
kubectl apply -f k8s/

# Check status
kubectl get pods -n otelcol-securityevent
```

## Step 4: Verify It's Working

### Check Health

```bash
curl http://localhost:13133/
```

### Check Logs

```bash
# Docker
docker logs otelcol-securityevent

# Kubernetes
kubectl logs -n otelcol-securityevent -l app.kubernetes.io/name=otelcol-securityevent -f
```

You should see logs indicating the collector is processing OpenReports Custom Resources.

## What's Next?

- [Configuration Guide](../configuration/processor-config.md): Learn about all configuration options
- [OpenReports Processor](../processors/openreports.md): Detailed documentation for OpenReports
- [Examples](../examples/basic.md): More examples and use cases

