# RBAC Setup

Configure RBAC permissions for the k8sobjects receiver.

## Overview

The k8sobjects receiver requires Kubernetes RBAC permissions to access OpenReports Custom Resources. The deployment includes ClusterRole and ClusterRoleBinding resources.

## Required Permissions

The collector needs the following permissions:

- `get`, `list`, `watch` on `reports` resource in `openreports.io` API group

## Included RBAC Resources

The Kubernetes manifests include:

1. **ServiceAccount** (`k8s/serviceaccount.yaml`): Service account for the collector
2. **ClusterRole** (`k8s/clusterrole.yaml`): Defines permissions
3. **ClusterRoleBinding** (`k8s/clusterrolebinding.yaml`): Binds permissions to ServiceAccount

## ClusterRole

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: otelcol-securityevent
rules:
  - apiGroups: ["openreports.io"]
    resources: ["reports"]
    verbs: ["get", "list", "watch"]
```

## ClusterRoleBinding

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: otelcol-securityevent
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: otelcol-securityevent
subjects:
  - kind: ServiceAccount
    name: otelcol-securityevent
    namespace: otelcol-securityevent
```

## Verification

### Check RBAC Resources

```bash
# Check ClusterRole
kubectl get clusterrole otelcol-securityevent

# Check ClusterRoleBinding
kubectl get clusterrolebinding otelcol-securityevent

# Check ServiceAccount
kubectl get serviceaccount otelcol-securityevent -n otelcol-securityevent
```

### Test Permissions

```bash
# Test get permission
kubectl auth can-i get reports -n default \
  --as=system:serviceaccount:otelcol-securityevent:otelcol-securityevent

# Test list permission
kubectl auth can-i list reports -n default \
  --as=system:serviceaccount:otelcol-securityevent:otelcol-securityevent

# Test watch permission
kubectl auth can-i watch reports -n default \
  --as=system:serviceaccount:otelcol-securityevent:otelcol-securityevent
```

All commands should return `yes`.

## Custom Permissions

### Additional Resources

To collect additional Kubernetes resources, add to ClusterRole:

```yaml
rules:
  - apiGroups: ["openreports.io"]
    resources: ["reports"]
    verbs: ["get", "list", "watch"]
  # Additional resources
  - apiGroups: [""]
    resources: ["events", "pods"]
    verbs: ["get", "list", "watch"]
```

### Namespace-Scoped Permissions

If you only need namespace-scoped permissions, use Role and RoleBinding instead:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: otelcol-securityevent
  namespace: <namespace>
rules:
  - apiGroups: ["openreports.io"]
    resources: ["reports"]
    verbs: ["get", "list", "watch"]
```

## Troubleshooting

### Permission Denied Errors

If you see permission denied errors:

1. **Verify RBAC resources exist**:
   ```bash
   kubectl get clusterrole otelcol-securityevent
   kubectl get clusterrolebinding otelcol-securityevent
   ```

2. **Check ServiceAccount**:
   ```bash
   kubectl get pod -n otelcol-securityevent \
     -o jsonpath='{.items[0].spec.serviceAccountName}'
   ```

3. **Test permissions manually** (see above)

4. **Check OpenReports CRD exists**:
   ```bash
   kubectl get crd reports.openreports.io
   ```

### Collector Logs

Check collector logs for RBAC errors:

```bash
kubectl logs -n otelcol-securityevent \
  -l app.kubernetes.io/name=otelcol-securityevent | grep -i "permission\|forbidden\|unauthorized"
```

## Next Steps

- [Kubernetes Deployment](kubernetes.md): Complete deployment guide
- [k8sobjects Receiver](../configuration/receivers.md#k8sobjects-receiver): Receiver configuration

