# API Reference

API reference for the Security Log Event Processor.

## Processor Configuration

### SecurityEvent Processor

```yaml
processors:
  securityevent:
    processors:
      <processor-name>:
        enabled: <boolean>
        # Processor-specific options
```

## OpenReports Processor API

### Configuration

```yaml
processors:
  securityevent:
    processors:
      openreports:
        enabled: boolean          # Required: Enable/disable processor
        status_filter: []string   # Optional: Filter by status
```

### Status Filter Values

- `pass`: Policy evaluation passed
- `fail`: Policy evaluation failed
- `error`: Policy evaluation error
- `skip`: Policy evaluation skipped

## Next Steps

- [Configuration Guide](../configuration/processor-config.md): Detailed configuration
- [Field Mapping](field-mapping.md): Field mappings

