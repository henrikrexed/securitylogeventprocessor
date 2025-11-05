# Security Log Event Processor

Welcome to the Security Log Event Processor documentation!

The Security Log Event Processor is a custom OpenTelemetry Collector processor that transforms security-related logs into standardized security event logs. It supports multiple processor types, each designed to handle specific security log formats.

## Features

- 🔄 **Transform Security Logs**: Convert various security log formats into standardized security events
- 🎯 **Multiple Processors**: Support for different security log sources (OpenReports, etc.)
- 📊 **Telemetry**: Built-in metrics and logging for observability
- 🚀 **Easy Deployment**: Docker and Kubernetes deployment options
- ⚙️ **Flexible Configuration**: Configurable filtering and field mapping

## Supported Processors

### OpenReports Processor

Processes OpenReports Custom Resource logs from Kubernetes, transforming policy evaluation results into security events.

**Recommended Receiver**: `k8sobjects` receiver

**Features**:
- Watches OpenReports Custom Resources in Kubernetes
- Transforms each policy evaluation result into a separate security event
- Configurable status filtering (pass, fail, error, skip)
- Automatic Kubernetes workload identification

## Quick Start

```yaml
# Basic configuration
receivers:
  k8sobjects:
    auth_type: serviceAccount
    objects:
      - name: openreports
        mode: watch
        api_group: openreports.io
        api_version: v1alpha1
        resource: reports

processors:
  securityevent:
    processors:
      openreports:
        enabled: true
        status_filter:
          - "fail"
          - "error"

service:
  pipelines:
    logs:
      receivers: [k8sobjects]
      processors: [securityevent, batch]
      exporters: [otlp]
```

## Documentation Structure

- **[Getting Started](getting-started/overview.md)**: Learn the basics and get up and running
- **[Configuration](configuration/processor-config.md)**: Configure processors, receivers, and exporters
- **[Deployment](deployment/overview.md)**: Deploy using Docker or Kubernetes
- **[Processors](processors/overview.md)**: Detailed processor documentation
- **[Examples](examples/basic.md)**: Practical examples and use cases

## Need Help?

- Check the [Troubleshooting](troubleshooting.md) guide
- Review [Examples](examples/basic.md) for common configurations
- See [Reference](reference/api.md) for detailed API documentation

