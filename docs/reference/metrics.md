# Metrics Reference

Telemetry metrics exposed by the Security Log Event Processor.

## Available Metrics

### Incoming Logs

**Metric**: `securityevent.logs.incoming`  
**Type**: Counter  
**Description**: Total number of incoming log records  
**Labels**: None

### Outgoing Logs

**Metric**: `securityevent.logs.outgoing`  
**Type**: Counter  
**Description**: Total number of outgoing log records (after processing)  
**Labels**: None

### Dropped Logs

**Metric**: `securityevent.logs.dropped`  
**Type**: Counter  
**Description**: Total number of logs dropped due to processing errors  
**Labels**: None

### Processing Errors

**Metric**: `securityevent.processing.errors`  
**Type**: Counter  
**Description**: Total number of processing errors  
**Labels**:
- `error_type`: Type of error (e.g., "processing_error")

## Metric Collection

Metrics are exposed via the collector's telemetry system and can be scraped by Prometheus or other monitoring systems.

### Example Query

```promql
# Rate of incoming logs
rate(securityevent.logs.incoming[5m])

# Rate of dropped logs
rate(securityevent.logs.dropped[5m])

# Error rate
rate(securityevent.processing.errors[5m])
```

## Next Steps

- [Configuration Guide](../configuration/processor-config.md): Configuration options
- [Troubleshooting](../troubleshooting.md): Debugging guide

