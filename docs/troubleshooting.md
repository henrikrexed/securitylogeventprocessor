# Troubleshooting

Common issues and solutions for the Security Log Event Processor.

## Collector Not Starting

### Configuration Errors

**Symptoms**: Collector fails to start with configuration errors

**Solutions**:
1. Validate YAML syntax
2. Check required fields are present
3. Verify processor is enabled
4. Check receiver configuration

```bash
# Test configuration
./otelcol-securityevent --config=config.yaml --dry-run
```

### Missing Dependencies

**Symptoms**: Collector fails with "component not found" errors

**Solutions**:
1. Rebuild collector with OCB
2. Verify manifest includes required components
3. Check component versions match

## No Events Generated

### Processor Not Enabled

**Symptoms**: Logs received but no security events generated

**Solutions**:
1. Verify processor is enabled:
   ```yaml
   securityevent:
     processors:
       openreports:
         enabled: true  # Must be true
   ```

2. Check status filter:
   ```yaml
   status_filter:
     - "fail"  # Verify your logs have matching statuses
   ```

### Receiver Not Collecting

**Symptoms**: No logs received from receiver

**Solutions**:

**k8sobjects receiver**:
- Verify OpenReports CRD exists: `kubectl get crd reports.openreports.io`
- Check RBAC permissions
- Verify namespace filters match actual namespaces
- Check label selectors match resource labels

**OTLP receiver**:
- Verify endpoints are accessible
- Check firewall rules
- Verify client is sending to correct endpoint

## RBAC Issues

### Permission Denied

**Symptoms**: "forbidden" or "unauthorized" errors in logs

**Solutions**:
1. Verify ClusterRole exists:
   ```bash
   kubectl get clusterrole otelcol-securityevent
   ```

2. Verify ClusterRoleBinding exists:
   ```bash
   kubectl get clusterrolebinding otelcol-securityevent
   ```

3. Test permissions:
   ```bash
   kubectl auth can-i get reports -n default \
     --as=system:serviceaccount:otelcol-securityevent:otelcol-securityevent
   ```

4. Check ServiceAccount is configured in deployment

## Performance Issues

### High Memory Usage

**Symptoms**: Collector using excessive memory

**Solutions**:
1. Adjust memory limiter:
   ```yaml
   memory_limiter:
     limit_mib: 2000  # Reduce limit
   ```

2. Reduce batch size:
   ```yaml
   batch:
     send_batch_size: 512  # Reduce batch size
   ```

3. Increase resource limits in Kubernetes

### Slow Processing

**Symptoms**: Events processed slowly

**Solutions**:
1. Increase batch timeout:
   ```yaml
   batch:
     timeout: 5s  # Increase timeout
   ```

2. Reduce status filter if processing too many events
3. Scale up collector replicas

## Debugging

### Enable Debug Logging

```yaml
exporters:
  logging:
    loglevel: debug
```

### Check Collector Logs

**Docker**:
```bash
docker logs otelcol-securityevent
```

**Kubernetes**:
```bash
kubectl logs -n otelcol-securityevent \
  -l app.kubernetes.io/name=otelcol-securityevent -f
```

### Check Health Endpoint

```bash
curl http://localhost:13133/
```

### Check zpages

```bash
# Port forward
kubectl port-forward -n otelcol-securityevent \
  svc/otelcol-securityevent 55679:55679

# Open in browser
open http://localhost:55679
```

## Common Error Messages

### "component not found"

**Cause**: Component not included in OCB manifest

**Solution**: Rebuild collector with updated manifest

### "permission denied"

**Cause**: Missing RBAC permissions

**Solution**: Apply ClusterRole and ClusterRoleBinding

### "no matching processor"

**Cause**: Log format doesn't match any processor

**Solution**: Check log format matches expected format

### "invalid status filter"

**Cause**: Invalid status value in status_filter

**Solution**: Use valid statuses: pass, fail, error, skip

## Getting Help

1. Check logs with debug level enabled
2. Review configuration examples
3. Check RBAC permissions
4. Verify OpenReports CRD exists
5. Review [Examples](examples/basic.md)

## Next Steps

- [Configuration Guide](configuration/processor-config.md): Review configuration
- [Examples](examples/basic.md): Review examples
- [Reference](reference/api.md): API reference

