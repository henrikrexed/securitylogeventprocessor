# Kubernetes Deployment Manifests

This directory contains Kubernetes manifests for deploying the OpenTelemetry Collector with Security Event Processor.

## Files

- `namespace.yaml` - Creates the `otelcol-securityevent` namespace
- `serviceaccount.yaml` - Service account for the collector pods
- `clusterrole.yaml` - RBAC permissions for k8sobjects receiver (OpenReports CR access)
- `clusterrolebinding.yaml` - Binds ClusterRole to ServiceAccount
- `configmap.yaml` - Collector configuration (includes k8sobjects receiver)
- `deployment.yaml` - Deployment with 2 replicas, health checks, and resource limits
- `service.yaml` - ClusterIP service and headless service
- `kustomization.yaml` - Kustomize configuration for managing all resources

## Prerequisites

1. **Docker Image**: Build and push the Docker image to your registry
   ```bash
   # Build the image
   make docker-build
   
   # Tag for your registry
   docker tag otelcol-securityevent:latest your-registry/otelcol-securityevent:0.1.0
   
   # Push to registry
   docker push your-registry/otelcol-securityevent:0.1.0
   ```

2. **Kubernetes Cluster**: Access to a Kubernetes cluster (v1.19+)

3. **kubectl**: Kubernetes CLI tool configured

4. **OpenReports CRD**: The OpenReports Custom Resource Definition must be installed in your cluster
   ```bash
   # Verify OpenReports CRD exists
   kubectl get crd reports.openreports.io
   ```

## Deployment Options

### Option 1: Using kubectl (Individual Files)

```bash
# Apply all manifests in order
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/serviceaccount.yaml
kubectl apply -f k8s/clusterrole.yaml
kubectl apply -f k8s/clusterrolebinding.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

### Option 2: Using kubectl (All at Once)

```bash
kubectl apply -f k8s/
```

### Option 3: Using kustomize

```bash
# Build and apply
kubectl apply -k k8s/

# Or build first to see what will be applied
kubectl kustomize k8s/ | kubectl apply -f -
```

## Configuration

### Update Image

Before deploying, update the image in `deployment.yaml`:

```yaml
containers:
- name: otelcol
  image: your-registry/otelcol-securityevent:0.1.0  # Update this
  imagePullPolicy: IfNotPresent
```

### Customize Configuration

Edit `configmap.yaml` to customize:
- **k8sobjects receiver**: Configure which OpenReports to collect
  - Namespace filters
  - Label selectors
  - API group/version/resource
- OTLP endpoint settings
- Processor configuration
- Status filters for OpenReports
- Log levels

#### k8sobjects Receiver Configuration

The collector uses the `k8sobjects` receiver to watch OpenReports Custom Resources. You can customize the collection:

```yaml
k8sobjects:
  auth_type: serviceAccount
  objects:
    - name: openreports
      mode: watch
      api_group: openreports.io
      api_version: v1alpha1
      resource: reports
      # Filter by namespace (optional)
      namespaces: [default, production]
      # Filter by label selector (optional)
      label_selector: environment in (production,staging)
```

### Environment Variables

You can set environment variables in `deployment.yaml`:
- `OTEL_LOG_LEVEL` - Log level (default: "info")
- `OTLP_ENDPOINT` - OTLP exporter endpoint
- `OTLP_TLS_INSECURE` - Disable TLS (default: "false")

### Resource Limits

Default resources in `deployment.yaml`:
- Requests: 100m CPU, 256Mi memory
- Limits: 1000m CPU, 2Gi memory

Adjust based on your workload.

## Verification

### Check Deployment Status

```bash
# Check namespace
kubectl get namespace otelcol-securityevent

# Check pods
kubectl get pods -n otelcol-securityevent

# Check services
kubectl get svc -n otelcol-securityevent

# Check deployment
kubectl get deployment -n otelcol-securityevent
```

### View Logs

```bash
# All pods
kubectl logs -n otelcol-securityevent -l app.kubernetes.io/name=otelcol-securityevent

# Specific pod
kubectl logs -n otelcol-securityevent <pod-name>

# Follow logs
kubectl logs -n otelcol-securityevent -l app.kubernetes.io/name=otelcol-securityevent -f
```

### Check Health

```bash
# Port forward to health endpoint
kubectl port-forward -n otelcol-securityevent svc/otelcol-securityevent 13133:13133

# Check health
curl http://localhost:13133/
```

### Access zpages

```bash
# Port forward to zpages
kubectl port-forward -n otelcol-securityevent svc/otelcol-securityevent 55679:55679

# Open in browser
open http://localhost:55679
```

## Scaling

### Scale Deployment

```bash
# Scale to 3 replicas
kubectl scale deployment otelcol-securityevent -n otelcol-securityevent --replicas=3
```

### Horizontal Pod Autoscaler (HPA)

Create an HPA for automatic scaling:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: otelcol-securityevent-hpa
  namespace: otelcol-securityevent
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: otelcol-securityevent
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

## Monitoring

### Prometheus ServiceMonitor (if using Prometheus Operator)

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: otelcol-securityevent
  namespace: otelcol-securityevent
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: otelcol-securityevent
  endpoints:
  - port: healthcheck
    path: /metrics
    interval: 30s
```

## Troubleshooting

### Pod Not Starting

```bash
# Check pod events
kubectl describe pod -n otelcol-securityevent <pod-name>

# Check pod logs
kubectl logs -n otelcol-securityevent <pod-name>
```

### Configuration Issues

```bash
# View configmap
kubectl get configmap otelcol-securityevent-config -n otelcol-securityevent -o yaml

# Check if config is mounted
kubectl exec -n otelcol-securityevent <pod-name> -- cat /etc/otelcol/config.yaml
```

### Image Pull Errors

```bash
# Check image pull secrets if using private registry
kubectl get secrets -n otelcol-securityevent

# Update deployment with imagePullSecrets if needed
```

### k8sobjects Receiver Issues

If the collector is not receiving OpenReports logs:

```bash
# Check if RBAC permissions are correct
kubectl get clusterrole otelcol-securityevent -o yaml
kubectl get clusterrolebinding otelcol-securityevent -o yaml

# Verify ServiceAccount is being used
kubectl get pod -n otelcol-securityevent -o jsonpath='{.items[0].spec.serviceAccountName}'

# Check collector logs for k8sobjects errors
kubectl logs -n otelcol-securityevent -l app.kubernetes.io/name=otelcol-securityevent | grep k8sobjects

# Verify OpenReports CRD exists
kubectl get crd reports.openreports.io

# Test RBAC permissions manually
kubectl auth can-i get reports -n default --as=system:serviceaccount:otelcol-securityevent:otelcol-securityevent
kubectl auth can-i list reports -n default --as=system:serviceaccount:otelcol-securityevent:otelcol-securityevent
kubectl auth can-i watch reports -n default --as=system:serviceaccount:otelcol-securityevent:otelcol-securityevent
```

## Cleanup

```bash
# Delete all resources
kubectl delete -f k8s/

# Or using kustomize
kubectl delete -k k8s/

# Delete namespace (deletes everything in namespace)
kubectl delete namespace otelcol-securityevent
```

## Production Considerations

1. **Image Registry**: Use a proper container registry (not local images)
2. **Resource Limits**: Adjust based on production workload
3. **Replicas**: Consider higher replica count for HA
4. **Storage**: Add persistent volumes if needed
5. **Network Policies**: Implement network policies for security
6. **Secrets**: Use Kubernetes secrets for sensitive data
7. **Monitoring**: Set up proper monitoring and alerting
8. **Backup**: Implement configuration backup strategy

