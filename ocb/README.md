# OpenTelemetry Collector Builder (OCB) Configuration

This directory contains OCB manifest files and configuration for building a custom OpenTelemetry Collector distribution that includes the Security Event Processor.

## Files

- `manifest.yaml` - Full OCB manifest with all components
- `manifest-simple.yaml` - Minimal OCB manifest with essential components only
- `config.yaml` - Default collector configuration
- `build.sh` - Script to build the collector using OCB

## Prerequisites

1. **Install OCB (OpenTelemetry Collector Builder)**
   ```bash
   go install go.opentelemetry.io/collector/cmd/builder@latest
   ```

2. **Verify installation**
   ```bash
   otelcolbuilder --version
   ```

## Building the Collector

### Using the build script (Recommended)

```bash
# Make script executable
chmod +x ocb/build.sh

# Build with default manifest
./ocb/build.sh

# Build with specific manifest
./ocb/build.sh manifest-simple.yaml
```

### Manual build

```bash
# Build with default manifest
otelcolbuilder --manifest manifest.yaml --output-path ./dist

# Build with simple manifest
otelcolbuilder --manifest manifest-simple.yaml --output-path ./dist
```

## Output

After building, the collector binary will be in the `dist/` directory:
- `dist/otelcol-securityevent_linux_amd64` (or other platform-specific names)

## Manifest Files

### `manifest.yaml`
Full-featured collector with:
- Multiple receivers (OTLP, Filelog)
- Essential processors (Batch, MemoryLimiter, Resource, Transform)
- **Security Event Processor** (custom)
- Multiple exporters (OTLP, OTLP HTTP, Logging, Debug)
- Extensions (HealthCheck, pprof, zpages)

### `manifest-simple.yaml`
Minimal collector with:
- OTLP receiver
- Essential processors (Batch, MemoryLimiter)
- **Security Event Processor** (custom)
- Basic exporters (OTLP, Logging)
- HealthCheck extension

## Configuration

The `config.yaml` file provides a default configuration that:
- Enables the Security Event Processor
- Configures OpenReports processor with status filter (fail, error)
- Sets up OTLP receivers and exporters
- Configures health check and debugging extensions

### Customizing Configuration

Edit `config.yaml` to customize:
- Processor settings
- Status filters
- Export destinations
- Log levels

## Including in Docker Build

The OCB build is integrated into the Docker build process (see `docker/Dockerfile.ocb`), which:
1. Uses OCB to build the collector
2. Creates a minimal Alpine-based runtime image
3. Includes the collector binary and configuration

## Troubleshooting

### OCB not found
```bash
# Install OCB
go install go.opentelemetry.io/collector/cmd/builder@latest

# Add to PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### Build failures
- Check Go version (requires Go 1.21+)
- Verify all dependencies with `go mod tidy`
- Check manifest syntax

### Module not found
- Ensure `go.mod` is in the project root
- Run `go mod download` before building
- Verify the local processor path in manifest.yaml

## Next Steps

1. Build the collector: `./ocb/build.sh`
2. Test locally: `./dist/otelcol-securityevent_linux_amd64 --config=ocb/config.yaml`
3. Build Docker image: `./docker/build.sh`
4. Deploy: Use Docker Compose or Kubernetes manifests

