# Field Mapping Reference

Complete field mapping from source logs to security events.

## OpenReports to Security Event Mapping

### Event Fields

| Source Field | Target Field | Type | Description |
|--------------|-------------|------|-------------|
| Generated UUID | `event.id` | string | Unique event identifier |
| `result.policy` | `event.type` | string | Policy name |
| `result.rule` | `event.description` | string | Rule description |
| `metadata.timestamp` | `event.timestamp` | timestamp | Event timestamp |

### Security Fields

| Source Field | Target Field | Type | Description |
|--------------|-------------|------|-------------|
| Calculated from status | `dt.security.risk.level` | string | CRITICAL, HIGH, MEDIUM, LOW |
| Calculated from status | `dt.security.risk.score` | number | 0-10 risk score |
| `result.status` | `compliance.status` | string | COMPLIANT, NON_COMPLIANT |

### Kubernetes Fields

| Source Field | Target Field | Type | Description |
|--------------|-------------|------|-------------|
| `metadata.namespace` | `k8s.namespace.name` | string | Kubernetes namespace |
| `metadata.name` | `k8s.resource.name` | string | Resource name |
| Extracted from ownerReferences | `k8s.workload.name` | string | Workload name |
| Extracted from ownerReferences | `k8s.workload.kind` | string | Workload kind |
| Extracted from ownerReferences | `k8s.workload.uid` | string | Workload UID |
| `metadata.namespace` | `k8s.workload.namespace` | string | Workload namespace |

## Risk Level Calculation

Risk levels are calculated based on result status:

| Status | Risk Level | Risk Score |
|--------|-----------|------------|
| `fail` | HIGH | 8.5 |
| `error` | CRITICAL | 9.5 |
| `pass` | LOW | 2.0 |
| `skip` | MEDIUM | 5.0 |

## Compliance Status Mapping

| Source Status | Compliance Status |
|---------------|-------------------|
| `pass` | COMPLIANT |
| `fail` | NON_COMPLIANT |
| `error` | NON_COMPLIANT |
| `skip` | NON_COMPLIANT |

## Example Mapping

### Input: OpenReports Log

```json
{
  "kind": "Report",
  "apiVersion": "openreports.io/v1alpha1",
  "metadata": {
    "name": "scan-001",
    "namespace": "production",
    "ownerReferences": [
      {
        "kind": "Deployment",
        "name": "my-app",
        "uid": "12345678-1234-1234-1234-123456789012"
      }
    ]
  },
  "results": [
    {
      "policy": "pod-security",
      "rule": "no-privileged-containers",
      "status": "fail"
    }
  ]
}
```

### Output: Security Event

```json
{
  "event": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "pod-security",
    "description": "no-privileged-containers",
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "dt.security.risk.level": "HIGH",
  "dt.security.risk.score": 8.5,
  "compliance.status": "NON_COMPLIANT",
  "k8s.namespace.name": "production",
  "k8s.resource.name": "scan-001",
  "k8s.workload.name": "my-app",
  "k8s.workload.kind": "Deployment",
  "k8s.workload.uid": "12345678-1234-1234-1234-123456789012",
  "k8s.workload.namespace": "production"
}
```

## Next Steps

- [OpenReports Processor](../processors/openreports.md): Processor details
- [Configuration Guide](../configuration/openreports.md): Configuration options

