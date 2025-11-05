# OpenReports Processor

Detailed documentation for the OpenReports processor.

## Overview

The OpenReports processor transforms OpenReports Custom Resource logs from Kubernetes into standardized security events. Each policy evaluation result in an OpenReports log becomes a separate security event.

## Features

- ✅ Watches OpenReports Custom Resources via k8sobjects receiver
- ✅ Expands single log into multiple security events (one per result)
- ✅ Configurable status filtering (pass, fail, error, skip)
- ✅ Automatic Kubernetes workload identification
- ✅ Rich metadata extraction (namespace, labels, ownerReferences)

## Input Format

The processor expects OpenReports Custom Resource logs in this format:

```json
{
  "kind": "Report",
  "apiVersion": "openreports.io/v1alpha1",
  "metadata": {
    "name": "report-name",
    "namespace": "default"
  },
  "results": [
    {
      "policy": "policy-name",
      "rule": "rule-name",
      "status": "fail"
    }
  ]
}
```

## Output Format

Each result produces a security event:

```json
{
  "event": {
    "id": "uuid",
    "type": "policy-name",
    "description": "rule-name"
  },
  "dt.security.risk.level": "HIGH",
  "compliance.status": "NON_COMPLIANT",
  "k8s": {
    "namespace.name": "default",
    "resource.name": "report-name",
    "workload.name": "workload-name",
    "workload.kind": "Deployment"
  }
}
```

## Configuration

See [OpenReports Configuration](../configuration/openreports.md) for detailed configuration options.

## Field Mapping

See [Field Mapping Reference](../reference/field-mapping.md) for complete field mappings.

### Key Mappings

- `result.policy` → `event.type`
- `result.rule` → `event.description`
- `result.status` → `compliance.status`
- `metadata.namespace` → `k8s.namespace.name`

## Status Filtering

Filter which results are processed:

```yaml
openreports:
  enabled: true
  status_filter:
    - "fail"    # Process failures
    - "error"   # Process errors
```

**Valid Statuses**:
- `pass`: Policy evaluation passed
- `fail`: Policy evaluation failed
- `error`: Policy evaluation error
- `skip`: Policy evaluation skipped

## Workload Identification

The processor automatically identifies Kubernetes workloads from:

1. **ownerReferences** (primary): Extracts workload from `metadata.ownerReferences`
2. **Pod name inference** (fallback): Infers workload from pod name patterns

**Supported Workload Kinds**:
- Deployment
- StatefulSet
- DaemonSet
- ReplicaSet
- Job
- CronJob

## Example

### Input: OpenReports Log

```json
{
  "kind": "Report",
  "apiVersion": "openreports.io/v1alpha1",
  "metadata": {
    "name": "security-scan-001",
    "namespace": "production"
  },
  "results": [
    {
      "policy": "pod-security",
      "rule": "no-privileged-containers",
      "status": "fail"
    },
    {
      "policy": "network-policy",
      "rule": "require-network-policy",
      "status": "pass"
    }
  ]
}
```

### Output: Security Events (with status_filter: ["fail"])

**Event 1**:
```json
{
  "event": {
    "id": "uuid-1",
    "type": "pod-security",
    "description": "no-privileged-containers"
  },
  "compliance.status": "NON_COMPLIANT",
  "dt.security.risk.level": "HIGH",
  "k8s.namespace.name": "production"
}
```

**Event 2**: (filtered out because status is "pass")

## Use Cases

1. **Compliance Monitoring**: Track policy violations in real-time
2. **Security Alerting**: Generate alerts for failed policy evaluations
3. **Audit Logging**: Create detailed audit trails
4. **Compliance Reporting**: Aggregate compliance status

## Troubleshooting

### No Events Generated

1. Check if processor is enabled
2. Verify status_filter includes your result statuses
3. Check k8sobjects receiver is working
4. Enable debug logging

### Incorrect Workload Information

1. Verify ownerReferences in OpenReports metadata
2. Check pod name follows standard patterns
3. Review collector logs for extraction errors

## Next Steps

- [Configuration Guide](../configuration/openreports.md): Configure the processor
- [Field Mapping Reference](../reference/field-mapping.md): Complete field mappings
- [Examples](../examples/openreports-k8sobjects.md): Complete examples

