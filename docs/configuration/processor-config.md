# Processor Configuration

This guide explains how to configure the Security Event Processor.

## Basic Configuration

The Security Event Processor is configured under the `processors` section:

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

## Processor Structure

```yaml
processors:
  securityevent:
    processors:
      # Processor-specific configurations
      openreports:
        enabled: true
        # Processor-specific options
```

## OpenReports Processor Configuration

### Basic Configuration

```yaml
processors:
  securityevent:
    processors:
      openreports:
        enabled: true
```

### Status Filtering

Filter which result statuses should be processed:

```yaml
processors:
  securityevent:
    processors:
      openreports:
        enabled: true
        status_filter:
          - "fail"    # Process failed policy evaluations
          - "error"   # Process error statuses
          # - "pass"  # Skip passing evaluations
          # - "skip"  # Skip skipped evaluations
```

**Valid Status Values**:
- `pass`: Policy evaluation passed
- `fail`: Policy evaluation failed
- `error`: Policy evaluation error
- `skip`: Policy evaluation skipped

If `status_filter` is empty or not specified, all statuses will be processed.

### Complete Example

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

## Processor in Pipeline

The processor must be included in the service pipeline:

```yaml
service:
  pipelines:
    logs:
      receivers: [k8sobjects]
      processors: [memory_limiter, securityevent, batch]
      exporters: [otlp, logging]
```

**Recommended Processor Order**:
1. `memory_limiter` - Limit memory usage
2. `securityevent` - Transform security logs
3. `batch` - Batch events for efficiency

## Configuration Validation

The processor validates configuration at startup. Invalid configurations will cause the collector to fail to start.

### Common Errors

**Processor not enabled**:
```yaml
# ❌ Wrong
processors:
  securityevent:
    processors:
      openreports:
        enabled: false  # Processor won't process logs
```

**Invalid status filter**:
```yaml
# ❌ Wrong
processors:
  securityevent:
    processors:
      openreports:
        status_filter:
          - "invalid_status"  # Will cause validation error
```

## Environment Variables

You can use environment variables in configuration:

```yaml
processors:
  securityevent:
    processors:
      openreports:
        enabled: true
        status_filter: ${OPENREPORTS_STATUS_FILTER:fail,error}
```

## Next Steps

- [OpenReports Configuration](openreports.md): Detailed OpenReports configuration
- [Receivers Configuration](receivers.md): Configure receivers for data collection
- [Exporters Configuration](exporters.md): Configure exporters for data output

