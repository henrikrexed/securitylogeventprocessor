# OpenReports Processor Configuration

Detailed configuration guide for the OpenReports processor.

## Overview

The OpenReports processor transforms OpenReports Custom Resource logs into security events. Each policy evaluation result in an OpenReports log becomes a separate security event.

## Required Receiver

The OpenReports processor requires the **k8sobjects receiver** to collect OpenReports Custom Resources from Kubernetes.

See [Receivers Configuration](receivers.md#k8sobjects-receiver) for detailed receiver setup.

## Basic Configuration

```yaml
processors:
  securityevent:
    processors:
      openreports:
        enabled: true
```

## Configuration Options

### `enabled`

Enable or disable the OpenReports processor.

**Type**: `boolean`  
**Default**: `false`  
**Required**: Yes

```yaml
openreports:
  enabled: true
```

### `status_filter`

Array of result statuses to process. Only results with statuses in this list will be transformed into security events.

**Type**: `array of strings`  
**Default**: `[]` (all statuses)  
**Required**: No

**Valid Values**:
- `pass`: Policy evaluation passed
- `fail`: Policy evaluation failed
- `error`: Policy evaluation error
- `skip`: Policy evaluation skipped

```yaml
openreports:
  enabled: true
  status_filter:
    - "fail"
    - "error"
```

**Examples**:

Process only failures:
```yaml
status_filter:
  - "fail"
```

Process failures and errors:
```yaml
status_filter:
  - "fail"
  - "error"
```

Process all statuses (default):
```yaml
# No status_filter specified, or empty array
status_filter: []
```

## Complete Configuration Example

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

## Field Mapping

The processor maps OpenReports fields to security event fields. See [Field Mapping Reference](../reference/field-mapping.md) for details.

### Key Mappings

- `result.policy` → `event.type`
- `result.rule` → `event.description`
- `result.status` → `compliance.status`
- `metadata.name` → `k8s.resource.name`
- `metadata.namespace` → `k8s.namespace.name`

## Kubernetes Workload Identification

The processor automatically identifies Kubernetes workloads from:
1. `metadata.ownerReferences` (primary)
2. Pod name inference (fallback)

Workload information is added to security events:
- `k8s.workload.name`
- `k8s.workload.kind`
- `k8s.workload.uid`
- `k8s.workload.namespace`

## Output

Each OpenReports log with multiple results produces multiple security events:

**Input**: 1 OpenReports log with 3 results  
**Output**: 3 security event logs

## Troubleshooting

### Processor Not Processing Logs

1. **Check if processor is enabled**:
   ```yaml
   openreports:
     enabled: true  # Must be true
   ```

2. **Check status filter**:
   ```yaml
   # If status_filter is too restrictive, logs may be filtered out
   status_filter:
     - "fail"  # Only processes failures
   ```

3. **Check receiver configuration**:
   - Verify k8sobjects receiver is configured
   - Verify receiver is in the pipeline

4. **Check logs**:
   ```bash
   # Enable debug logging
   exporters:
     logging:
       loglevel: debug
   ```

### No Security Events Generated

- Check if OpenReports logs match the expected format
- Verify status_filter includes the statuses in your logs
- Check collector logs for processing errors

## Next Steps

- [Receivers Configuration](receivers.md): Configure k8sobjects receiver
- [Examples](../examples/openreports-k8sobjects.md): Complete examples
- [Field Mapping Reference](../reference/field-mapping.md): Detailed field mappings

