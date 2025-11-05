# Multi-stage Dockerfile for OpenTelemetry Collector with Security Event Processor
# This Dockerfile creates a minimal, secure image for production use

# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Configure Go module settings for proper resolution
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org
ENV GO111MODULE=on

# Set working directory
WORKDIR /app

# Copy all source code first (including go.mod, go.sum, and all source files)
COPY . .

# Build the collector using OCB
# OCB will handle dependency management, so we don't need go mod download here
RUN go install go.opentelemetry.io/collector/cmd/builder@v0.139.0
RUN mv $(go env GOPATH)/bin/builder $(go env GOPATH)/bin/ocb

# Initialize the local module to ensure Go can resolve it
# This helps OCB find the local module when processing the replace directive
RUN go mod download || true

# Run OCB with the manifest (with verbose output for debugging)
# First, generate code without compiling to inspect what's generated
RUN ocb --config ocb/manifest.yaml  --verbose || true

# Show generated files for debugging
RUN ls -la dist/ || true
RUN find dist/ -name "*.go" -type f | head -20 || true

# Now try to compile (this will show the actual error)
RUN ocb --config ocb/manifest.yaml --verbose

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    wget \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN adduser -D -s /bin/sh -u 10001 otelcol

# Create directories
RUN mkdir -p /etc/otelcol /var/lib/otelcol /tmp/otelcol \
    && chown -R otelcol:otelcol /etc/otelcol /var/lib/otelcol /tmp/otelcol

# Copy the collector binary from build stage
# OCB outputs platform-specific binaries - copy the entire dist directory first
COPY --from=builder /app/dist /tmp/dist

# Find and copy the binary (OCB may output with or without platform suffix)
RUN if [ -f /tmp/dist/otelcol-securityevent_linux_amd64 ]; then \
        cp /tmp/dist/otelcol-securityevent_linux_amd64 /otelcol-securityevent; \
    elif [ -f /tmp/dist/otelcol-securityevent ]; then \
        cp /tmp/dist/otelcol-securityevent /otelcol-securityevent; \
    else \
        echo "Error: Collector binary not found. Available files:"; \
        ls -la /tmp/dist/; \
        exit 1; \
    fi && \
    rm -rf /tmp/dist

# Set ownership
RUN chown otelcol:otelcol /otelcol-securityevent \
    && chmod +x /otelcol-securityevent

# Switch to non-root user
USER otelcol

# Expose ports
# 4317: OTLP gRPC receiver
# 4318: OTLP HTTP receiver
EXPOSE 4317 4318

# Set working directory
WORKDIR /etc/otelcol

# Default configuration
COPY ocb/config.yaml /etc/otelcol/config.yaml

# Set ownership of config file
USER root
RUN chown otelcol:otelcol /etc/otelcol/config.yaml
USER otelcol

# Set entrypoint
ENTRYPOINT ["/otelcol-securityevent"]
CMD ["--config=/etc/otelcol/config.yaml"]

# Labels
LABEL maintainer="security-event-processor-team"
LABEL description="OpenTelemetry Collector with Security Event Processor"
LABEL version="0.1.0"
LABEL collector.version="0.139.0"

