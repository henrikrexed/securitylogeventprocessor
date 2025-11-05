# Security Event Field Mapping

This document describes how OpenReports log fields are mapped to the Security Event schema.

## Field Mappings

### Event Fields

| Security Event Field | Source/Mapping | Notes |
|---------------------|----------------|-------|
| `event.id` | Generated UUID | Unique identifier for each security event |
| `event.version` | Hardcoded `"1.309"` | Fixed version |
| `event.category` | Hardcoded `"COMPLIANCE"` | Fixed category |
| `event.description` | Generated | Format: "Policy violation on {pod} for rule {rule}" (or appropriate message based on result) |
| `event.name` | Hardcoded `"Compliance finding event"` | Fixed name |
| `event.type` | Hardcoded `"COMPLIANCE_FINDING"` | Fixed type |

### Product Fields

| Security Event Field | Source/Mapping | Notes |
|---------------------|----------------|-------|
| `product.name` | Empty string | Currently not mapped |
| `product.vendor` | Empty string | Currently not mapped |

### Smartscape Fields

| Security Event Field | Source/Mapping | Notes |
|---------------------|----------------|-------|
| `smartscape.type` | `"K8S_POD"` if `scope.kind == "Pod"` | Set only for Pod resources |

### Security Risk Fields

| Security Event Field | Source/Mapping | Notes |
|---------------------|----------------|-------|
| `dt.security.risk.level` | Mapped from `finding.severity` | Mapping: critical→CRITICAL, high→HIGH, medium→MEDIUM, low→LOW, default→MEDIUM |
| `dt.security.risk.score` | Calculated from risk level | CRITICAL=10.0, HIGH=8.9, MEDIUM=6.9, LOW=3.9, default=0.0 |

### Object Fields

| Security Event Field | Source/Mapping | Notes |
|---------------------|----------------|-------|
| `object.id` | `scope.uid` | Kubernetes resource UID |
| `object.type` | `scope.kind` | Kubernetes resource kind (e.g., "Pod") |

### Finding Fields

| Security Event Field | Source/Mapping | Notes |
|---------------------|----------------|-------|
| `finding.description` | `result.message` | Message from the policy evaluation result |
| `finding.id` | Generated UUID | Unique identifier for the finding |
| `finding.severity` | `result.severity` | Original severity from result |
| `finding.time.created` | `result.timestamp` | Timestamp from result, formatted as RFC3339Nano |
| `finding.title` | `{policy} - {rule}` | Policy and rule combined |
| `finding.type` | `result.policy` | Policy name |
| `finding.url` | Empty string | Currently not mapped |

### Compliance Fields

| Security Event Field | Source/Mapping | Notes |
|---------------------|----------------|-------|
| `compliance.control` | `result.rule` | Rule name from result |
| `compliance.requirements` | `result.policy` | Policy name from result |
| `compliance.standards` | `result.category` (if available) | Category from result, or omitted |
| `compliance.status` | Mapped from `result.result` | Mapping: pass→COMPLIANT, fail/error/skip/unknown→NON_COMPLIANT |

### Kubernetes Fields

| Security Event Field | Source/Mapping | Notes |
|---------------------|----------------|-------|
| `k8s.*` | All `k8s.*` fields from original log | Copied from original OpenReports log attributes |
| `k8s.pod.name` | `scope.name` (if scope.kind is Pod) | Pod name |
| `k8s.namespace.name` | `scope.namespace` | Namespace name |
| `k8s.resource.kind` | `scope.kind` | Resource kind |
| `k8s.resource.uid` | `scope.uid` | Resource UID |
| `k8s.workload.name` | Extracted from `metadata.ownerReferences` or inferred from pod name | Workload name (Deployment, StatefulSet, etc.) |
| `k8s.workload.kind` | Extracted from `metadata.ownerReferences` or defaults to "Deployment" | Workload kind |
| `k8s.workload.namespace` | `scope.namespace` | Workload namespace (same as pod namespace) |
| `k8s.workload.uid` | Extracted from `metadata.ownerReferences` | Workload UID |
| `k8s.deployment.name` | `k8s.workload.name` (if workload.kind is Deployment) | Deployment name |
| `k8s.statefulset.name` | `k8s.workload.name` (if workload.kind is StatefulSet) | StatefulSet name |
| `k8s.daemonset.name` | `k8s.workload.name` (if workload.kind is DaemonSet) | DaemonSet name |

## Result Status Mapping

The `result.result` field from OpenReports is mapped to `compliance.status`:

- `"pass"` → `"COMPLIANT"`
- `"fail"` → `"NON_COMPLIANT"`
- `"error"` → `"NON_COMPLIANT"`
- `"skip"` → `"NON_COMPLIANT"`
- Unknown → `"NON_COMPLIANT"`

## Severity to Risk Level Mapping

The `finding.severity` field is mapped to `dt.security.risk.level`:

- `"critical"` → `"CRITICAL"`
- `"high"` → `"HIGH"`
- `"medium"` → `"MEDIUM"`
- `"low"` → `"LOW"`
- Unknown/empty → `"MEDIUM"` (default)

## Risk Score Calculation

Risk scores are calculated based on risk level:

- `CRITICAL` → 10.0
- `HIGH` → 8.9
- `MEDIUM` → 6.9
- `LOW` → 3.9
- Unknown → 0.0

