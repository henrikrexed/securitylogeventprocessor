# OpenReports with k8sobjects Receiver

Complete example configuration for OpenReports processor with k8sobjects receiver.

## Overview

This example shows how to:
1. Configure k8sobjects receiver to collect OpenReports Custom Resources
2. Configure OpenReports processor to transform logs
3. Export security events via OTLP

## Complete Configuration

```yaml
receivers:
  # k8sobjects receiver - collects OpenReports Custom Resources
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

processors:
  batch:
    timeout: 1s
    send_batch_size: 1024
  
  memory_limiter:
    check_interval: 1s
    limit_mib: 4000
  
  # Security Event Processor
  securityevent:
    processors:
      openreports:
        enabled: true
        status_filter:
          - "fail"
          - "error"

exporters:
  otlp:
    endpoint: http://your-backend:4317
    tls:
      insecure: false
  
  logging:
    loglevel: info

service:
  pipelines:
    logs:
      receivers: [k8sobjects]
      processors: [memory_limiter, securityevent, batch]
      exporters: [otlp, logging]
  
  extensions: [health_check, pprof, zpages]

extensions:
  health_check:
    endpoint: 0.0.0.0:13133
  pprof:
    endpoint: 0.0.0.0:1777
  zpages:
    endpoint: 0.0.0.0:55679
```

## Kubernetes Deployment

### 1. Apply RBAC

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/serviceaccount.yaml
kubectl apply -f k8s/clusterrole.yaml
kubectl apply -f k8s/clusterrolebinding.yaml
```

### 2. Create ConfigMap

```bash
kubectl create configmap otelcol-securityevent-config \
  --from-file=config.yaml=config.yaml \
  -n otelcol-securityevent
```

### 3. Deploy Collector

```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

## Verification

### Check Pods

```bash
kubectl get pods -n otelcol-securityevent
```

### Check Logs

```bash
kubectl logs -n otelcol-securityevent \
  -l app.kubernetes.io/name=otelcol-securityevent -f
```

You should see logs indicating:
- k8sobjects receiver started
- OpenReports processor enabled
- Security events being generated

### Check Health

```bash
kubectl port-forward -n otelcol-securityevent \
  svc/otelcol-securityevent 13133:13133

curl http://localhost:13133/
```

## Expected Output

When an OpenReports Custom Resource is created/updated, you should see:

1. **k8sobjects receiver** collects the resource
2. **OpenReports processor** transforms each result into a security event
3. **Exporter** sends events to configured backend

### Example Security Event

```json
{
  "event": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "pod-security-policy",
    "description": "no-privileged-containers",
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "dt.security.risk.level": "HIGH",
  "dt.security.risk.score": 8.5,
  "compliance.status": "NON_COMPLIANT",
  "k8s": {
    "namespace.name": "production",
    "resource.name": "security-scan-001",
    "workload.name": "my-app",
    "workload.kind": "Deployment",
    "workload.uid": "12345678-1234-1234-1234-123456789012"
  }
}
```

## Customization

### Filter by Namespace

```yaml
k8sobjects:
  objects:
    - name: openreports
      namespaces: [production, staging]
```

### Filter by Labels

```yaml
k8sobjects:
  objects:
    - name: openreports
      label_selector: environment=production
```

### Process All Statuses

```yaml
securityevent:
  processors:
    openreports:
      enabled: true
      # No status_filter = process all statuses
```

## Troubleshooting

### No Events Received

1. Check OpenReports CRD exists:
   ```bash
   kubectl get crd reports.openreports.io
   ```

2. Verify RBAC permissions:
   ```bash
   kubectl auth can-i get reports -n default \
     --as=system:serviceaccount:otelcol-securityevent:otelcol-securityevent
   ```

3. Check collector logs for errors

### Events Not Generated

1. Verify processor is enabled
2. Check status_filter matches your result statuses
3. Enable debug logging

## Next Steps

- [Configuration Guide](../configuration/openreports.md): Detailed configuration
- [RBAC Setup](../deployment/rbac.md): RBAC configuration
- [Field Mapping](../reference/field-mapping.md): Field mappings

