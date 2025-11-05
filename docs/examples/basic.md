# Basic Example

Basic example configuration for the Security Log Event Processor.

## Minimal Configuration

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

processors:
  securityevent:
    processors:
      openreports:
        enabled: true

exporters:
  logging:
    loglevel: info

service:
  pipelines:
    logs:
      receivers: [k8sobjects]
      processors: [securityevent]
      exporters: [logging]
```

## With Status Filtering

```yaml
processors:
  securityevent:
    processors:
      openreports:
        enabled: true
        status_filter:
          - "fail"
          - "error"
```

## With OTLP Exporter

```yaml
exporters:
  otlp:
    endpoint: http://your-backend:4317

service:
  pipelines:
    logs:
      receivers: [k8sobjects]
      processors: [securityevent, batch]
      exporters: [otlp]
```

## Next Steps

- [OpenReports Example](openreports-k8sobjects.md): Complete example
- [Configuration Guide](../configuration/processor-config.md): Configuration options

