# Exporters Configuration

This guide explains how to configure exporters for the Security Log Event Processor.

## Overview

Exporters send processed security events to various destinations. You can configure multiple exporters to send events to different systems.

## Common Exporters

### OTLP Exporter

Send events via OTLP protocol (recommended for most use cases).

```yaml
exporters:
  otlp:
    endpoint: http://your-backend:4317
    tls:
      insecure: false
      cert_file: /path/to/cert.pem
      key_file: /path/to/key.pem
```

**Configuration Options**:
- `endpoint`: OTLP endpoint URL
- `tls.insecure`: Skip TLS verification (not recommended for production)
- `tls.cert_file`: Client certificate file
- `tls.key_file`: Client private key file

### Logging Exporter

Output events to collector logs (useful for debugging).

```yaml
exporters:
  logging:
    loglevel: info  # debug, info, warn, error
```

### OTLP HTTP Exporter

Send events via OTLP over HTTP.

```yaml
exporters:
  otlphttp:
    endpoint: http://your-backend:4318
    tls:
      insecure: false
```

## Complete Example

```yaml
exporters:
  otlp:
    endpoint: http://dynatrace-backend:4317
    tls:
      insecure: false
  
  logging:
    loglevel: info

service:
  pipelines:
    logs:
      receivers: [k8sobjects]
      processors: [securityevent, batch]
      exporters: [otlp, logging]
```

## Next Steps

- [Processor Configuration](processor-config.md): Configure processors
- [Examples](../examples/basic.md): Complete examples

