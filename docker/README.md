# Docker Build for Security Event Processor

This directory contains Docker-related files for building and running the OpenTelemetry Collector with the Security Event Processor.

## Files

- `Dockerfile` - Multi-stage Dockerfile using OCB (OpenTelemetry Collector Builder)
- `docker-compose.yml` - Docker Compose configuration for local development
- `build.sh` - Script to build the Docker image
- `otlp-exporter-config.yaml` - Test configuration for OTLP exporter

## Building the Docker Image

The Dockerfile builds the collector inside Docker using OCB, so no local build is required.

### Option 1: Using the build script (Recommended)

```bash
# Make script executable
chmod +x docker/build.sh

# Build the image
./docker/build.sh
```

This script will:
1. Build the collector using OCB inside Docker
2. Create the final Docker image

### Option 2: Manual build

```bash
docker build -f docker/Dockerfile -t otelcol-securityevent:latest .
```

### Option 3: Using Docker Compose

```bash
# Build and start services
docker-compose -f docker/docker-compose.yml up -d --build

# Or just build
docker-compose -f docker/docker-compose.yml build
```

## Running the Container

### Using Docker

```bash
docker run -d \
  --name otelcol-securityevent \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 55679:55679 \
  -p 13133:13133 \
  -v $(pwd)/ocb/config.yaml:/etc/otelcol/config.yaml:ro \
  otelcol-securityevent:latest
```

### Using Docker Compose

```bash
docker-compose -f docker/docker-compose.yml up -d
```

## Configuration

The default configuration is located at `ocb/config.yaml`. You can override it by mounting a custom configuration:

```bash
docker run -v /path/to/your/config.yaml:/etc/otelcol/config.yaml:ro otelcol-securityevent:latest
```

## Ports

- `4317` - OTLP gRPC receiver
- `4318` - OTLP HTTP receiver
- `55679` - zpages extension (debugging)
- `13133` - Health check extension
- `1777` - pprof extension (profiling)

## Health Check

The container includes a health check endpoint:

```bash
curl http://localhost:13133/
```

## Testing

Send test logs to the collector:

```bash
# Using curl (HTTP)
curl -X POST http://localhost:4318/v1/logs \
  -H "Content-Type: application/json" \
  -d @test-log.json

# Using otel-cli or otel-collector-contrib
```

## Troubleshooting

### View logs
```bash
docker logs otelcol-securityevent
```

### Access zpages
```bash
# Open in browser
open http://localhost:55679
```

### Debug mode
```bash
docker run -e OTEL_LOG_LEVEL=debug otelcol-securityevent:latest
```

