# Receivers Configuration

This guide explains which receivers to use with the Security Event Processor and how to configure them.

## Overview

Receivers collect data from various sources. The processor you choose determines which receiver you need.

## Processor-Receiver Mapping

| Processor | Required Receiver | Optional Receivers |
|-----------|-------------------|-------------------|
| **OpenReports** | `k8sobjects` | `otlp`, `filelog` |

## k8sobjects Receiver

**Required for**: OpenReports processor

The k8sobjects receiver collects Kubernetes Custom Resources, making it ideal for collecting OpenReports logs.

### Basic Configuration

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
```

### Configuration Options

#### `auth_type`

Authentication type for Kubernetes API access.

**Type**: `string`  
**Default**: `serviceAccount`  
**Required**: Yes

**Valid Values**:
- `serviceAccount`: Use Kubernetes ServiceAccount (recommended for production)
- `kubeConfig`: Use kubeconfig file

```yaml
k8sobjects:
  auth_type: serviceAccount
```

#### `objects`

Array of Kubernetes objects to collect.

**Type**: `array`  
**Required**: Yes

Each object configuration:

```yaml
objects:
  - name: <unique-name>
    mode: <watch|pull>
    api_group: <api-group>
    api_version: <api-version>
    resource: <resource-name>
    # Optional filters
    namespaces: [<namespace-list>]
    label_selector: <label-selector>
    field_selector: <field-selector>
```

**Fields**:

- `name`: Unique identifier for this object collection
- `mode`: Collection mode
  - `watch`: Watch for changes (real-time)
  - `pull`: Periodically pull objects
- `api_group`: Kubernetes API group (e.g., `openreports.io`)
- `api_version`: API version (e.g., `v1alpha1`)
- `resource`: Resource name (e.g., `reports`)
- `namespaces`: (Optional) List of namespaces to watch
- `label_selector`: (Optional) Label selector filter
- `field_selector`: (Optional) Field selector filter

### OpenReports Configuration

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
        # Optional: Filter by namespace
        namespaces: [default, production]
        # Optional: Filter by labels
        label_selector: environment in (production,staging)
```

### Namespace Filtering

Watch OpenReports in specific namespaces:

```yaml
objects:
  - name: openreports
    mode: watch
    api_group: openreports.io
    api_version: v1alpha1
    resource: reports
    namespaces: [default, production, staging]
```

### Label Selector Filtering

Filter by labels:

```yaml
objects:
  - name: openreports
    mode: watch
    api_group: openreports.io
    api_version: v1alpha1
    resource: reports
    label_selector: environment in (production,staging)
```

### Kubernetes RBAC Requirements

The k8sobjects receiver requires RBAC permissions. See [RBAC Setup](../deployment/rbac.md) for details.

**Required Permissions**:
- `get`, `list`, `watch` on `reports` resource in `openreports.io` API group

## OTLP Receiver

**Optional**: Can be used alongside k8sobjects for receiving logs via OTLP protocol

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318
```

Useful for:
- Receiving logs from other collectors
- Testing with OTLP clients
- Receiving logs from applications

## filelog Receiver

**Optional**: Can be used for reading OpenReports logs from files

```yaml
receivers:
  filelog:
    include:
      - /var/log/openreports/*.log
    exclude: []
    start_at: end
```

## Complete Example

```yaml
receivers:
  # Primary receiver for OpenReports
  k8sobjects:
    auth_type: serviceAccount
    objects:
      - name: openreports
        mode: watch
        api_group: openreports.io
        api_version: v1alpha1
        resource: reports
  
  # Optional: OTLP receiver for additional log sources
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  securityevent:
    processors:
      openreports:
        enabled: true

service:
  pipelines:
    logs:
      receivers: [k8sobjects, otlp]
      processors: [securityevent, batch]
      exporters: [otlp]
```

## Troubleshooting

### k8sobjects Receiver Issues

**"Permission denied" errors**:
- Check RBAC permissions (ClusterRole, ClusterRoleBinding)
- Verify ServiceAccount is configured
- See [RBAC Setup](../deployment/rbac.md)

**No logs received**:
- Verify OpenReports CRD exists: `kubectl get crd reports.openreports.io`
- Check namespace filters match actual namespaces
- Verify label selectors match resource labels
- Enable debug logging: `loglevel: debug`

**Receiver not starting**:
- Check API group/version/resource are correct
- Verify authentication is configured
- Check collector logs for errors

## Next Steps

- [RBAC Setup](../deployment/rbac.md): Configure Kubernetes permissions
- [OpenReports Configuration](openreports.md): Configure OpenReports processor
- [Examples](../examples/openreports-k8sobjects.md): Complete examples

