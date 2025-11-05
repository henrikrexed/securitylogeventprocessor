# Kubernetes Deployment

Deploy the Security Log Event Processor to Kubernetes.

## Prerequisites

- Kubernetes cluster (v1.19+)
- `kubectl` configured
- OpenReports CRD installed
- Container registry access (for pushing images)

## Quick Start

```bash
# Apply all manifests
kubectl apply -f k8s/

# Or using kustomize
kubectl apply -k k8s/
```

## Build and Push Image

```bash
# Build the image
make docker-build

# Tag for your registry
docker tag otelcol-securityevent:latest \
  your-registry/otelcol-securityevent:0.1.0

# Push to registry
docker push your-registry/otelcol-securityevent:0.1.0
```

## Update Image in Deployment

Edit `k8s/deployment.yaml`:

```yaml
containers:
- name: otelcol
  image: your-registry/otelcol-securityevent:0.1.0
  imagePullPolicy: IfNotPresent
```

## Verify Deployment

```bash
# Check pods
kubectl get pods -n otelcol-securityevent

# Check services
kubectl get svc -n otelcol-securityevent

# View logs
kubectl logs -n otelcol-securityevent \
  -l app.kubernetes.io/name=otelcol-securityevent -f
```

## Configuration

### Update ConfigMap

Edit `k8s/configmap.yaml` to customize configuration:

```yaml
data:
  config.yaml: |
    receivers:
      k8sobjects:
        # Your configuration
    processors:
      securityevent:
        # Your configuration
```

Apply changes:

```bash
kubectl apply -f k8s/configmap.yaml
kubectl rollout restart deployment/otelcol-securityevent -n otelcol-securityevent
```

### Environment Variables

Set environment variables in `k8s/deployment.yaml`:

```yaml
env:
- name: OTEL_LOG_LEVEL
  value: "debug"
```

## RBAC

The deployment includes RBAC resources for k8sobjects receiver. See [RBAC Setup](rbac.md) for details.

## Scaling

### Manual Scaling

```bash
kubectl scale deployment otelcol-securityevent \
  -n otelcol-securityevent --replicas=3
```

### Horizontal Pod Autoscaler

Create an HPA:

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
```

## Monitoring

### Health Check

```bash
# Port forward
kubectl port-forward -n otelcol-securityevent \
  svc/otelcol-securityevent 13133:13133

# Check health
curl http://localhost:13133/
```

### Metrics

The collector exposes metrics on the health check endpoint. Configure Prometheus ServiceMonitor if using Prometheus Operator.

## Troubleshooting

### Pod Not Starting

```bash
# Check pod events
kubectl describe pod -n otelcol-securityevent <pod-name>

# Check logs
kubectl logs -n otelcol-securityevent <pod-name>
```

### RBAC Issues

```bash
# Check RBAC
kubectl get clusterrole otelcol-securityevent
kubectl get clusterrolebinding otelcol-securityevent

# Test permissions
kubectl auth can-i get reports -n default \
  --as=system:serviceaccount:otelcol-securityevent:otelcol-securityevent
```

## Next Steps

- [RBAC Setup](rbac.md): Configure RBAC permissions
- [Configuration Guide](../configuration/processor-config.md): Configure the processor

