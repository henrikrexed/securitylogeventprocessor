# Processors Overview

Overview of available processors in the Security Log Event Processor.

## Available Processors

### OpenReports Processor

Transforms OpenReports Custom Resource logs into security events.

**Status**: ✅ Available  
**Required Receiver**: `k8sobjects`  
**Documentation**: [OpenReports Processor](openreports.md)

**Features**:
- Watches OpenReports Custom Resources
- Expands single log into multiple security events (one per result)
- Configurable status filtering
- Automatic Kubernetes workload identification

## Processor Architecture

```
Input Log (OpenReports)
    │
    ├─ Result 1 → Security Event 1
    ├─ Result 2 → Security Event 2
    └─ Result 3 → Security Event 3
```

Each processor:
1. Receives logs from configured receivers
2. Identifies matching log format
3. Transforms logs into security events
4. Applies filters (if configured)
5. Outputs security events to exporters

## Processor Configuration Pattern

All processors follow the same configuration pattern:

```yaml
processors:
  securityevent:
    processors:
      <processor-name>:
        enabled: true
        # Processor-specific options
```

## Choosing a Processor

| Use Case | Processor | Receiver |
|----------|-----------|----------|
| OpenReports CR logs | OpenReports | k8sobjects |
| More coming soon... | - | - |

## Next Steps

- [OpenReports Processor](openreports.md): Detailed OpenReports documentation
- [Configuration Guide](../configuration/processor-config.md): Configure processors

