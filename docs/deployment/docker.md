# Docker Deployment

Deploy the Security Log Event Processor using Docker.

## Prerequisites

- Docker installed
- OpenReports CRD installed in your Kubernetes cluster (if using k8sobjects receiver)

## Build Docker Image

```bash
# Build the image
make docker-build

# Or using docker-compose
docker-compose -f docker/docker-compose.yml build
```

## Run Container

### Basic Run

```bash
docker run -d \
  --name otelcol-securityevent \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 13133:13133 \
  -v $(pwd)/config.yaml:/etc/otelcol/config.yaml:ro \
  otelcol-securityevent:latest
```

### Using Docker Compose

```bash
# Start services
docker-compose -f docker/docker-compose.yml up -d

# View logs
docker-compose -f docker/docker-compose.yml logs -f

# Stop services
docker-compose -f docker/docker-compose.yml down
```

## Configuration

Mount your configuration file:

```bash
docker run -v /path/to/config.yaml:/etc/otelcol/config.yaml:ro \
  otelcol-securityevent:latest
```

## Ports

- `4317`: OTLP gRPC receiver
- `4318`: OTLP HTTP receiver
- `13133`: Health check endpoint
- `55679`: zpages extension
- `1777`: pprof extension

## Health Check

```bash
curl http://localhost:13133/
```

## Logs

```bash
# View logs
docker logs otelcol-securityevent

# Follow logs
docker logs -f otelcol-securityevent
```

## Environment Variables

Set environment variables:

```bash
docker run -e OTEL_LOG_LEVEL=debug \
  otelcol-securityevent:latest
```

## Limitations

**Note**: Docker deployment cannot use the k8sobjects receiver directly because it requires Kubernetes API access. For k8sobjects receiver, use Kubernetes deployment.

For Docker, you can:
- Use OTLP receiver to receive logs from other sources
- Use filelog receiver to read logs from mounted volumes
- Run in Kubernetes-in-Docker (kind, minikube) for local development

## Next Steps

- [Kubernetes Deployment](kubernetes.md): Deploy to Kubernetes for k8sobjects support
- [Configuration Guide](../configuration/processor-config.md): Configure the processor

