# Multi-Processor Example

Example configuration using multiple processors (when available).

## Overview

This example shows how to configure multiple processors. Currently, only OpenReports processor is available, but this structure supports future processors.

## Configuration

```yaml
receivers:
  k8sobjects:
    auth_type: serviceAccount
    objects:
      - name: openreports
        mode: watch
        api_group: openreports.io
        api_version: v1alpha1
        resource: reports
  
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317

processors:
  securityevent:
    processors:
      openreports:
        enabled: true
        status_filter:
          - "fail"
          - "error"
      # Future processors will be added here

exporters:
  otlp:
    endpoint: http://your-backend:4317

service:
  pipelines:
    logs:
      receivers: [k8sobjects, otlp]
      processors: [securityevent, batch]
      exporters: [otlp]
```

## Processor Priority

Processors are checked in order. The first matching processor processes the log.

## Next Steps

- [Processors Overview](../processors/overview.md): Available processors
- [Configuration Guide](../configuration/processor-config.md): Configuration options

