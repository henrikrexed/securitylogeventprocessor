# Metrics

The Security Event Processor exposes metrics through the OpenTelemetry Collector's telemetry system.

## Available Metrics

### `processor_securityevent_incoming_logs_total`
- **Type**: Counter (Int64)
- **Description**: Total number of incoming logs processed by the processor
- **Unit**: 1 (count)
- **Labels**: None

### `processor_securityevent_outgoing_logs_total`
- **Type**: Counter (Int64)
- **Description**: Total number of outgoing logs produced by the processor
- **Unit**: 1 (count)
- **Labels**: None

**Note**: This metric counts:
- Logs that pass through unchanged (not OpenReports logs or filtered out)
- Each expanded security event log (when one OpenReports log expands to multiple security events)

### `processor_securityevent_dropped_logs_total`
- **Type**: Counter (Int64)
- **Description**: Total number of logs dropped during processing
- **Unit**: 1 (count)
- **Labels**: None

**Note**: Logs are counted as dropped when:
- Processing errors occur (e.g., JSON parsing failures)
- Errors during transformation

### `processor_securityevent_processing_errors_total`
- **Type**: Counter (Int64)
- **Description**: Total number of processing errors encountered
- **Unit**: 1 (count)
- **Labels**:
  - `error_type`: Type of error (e.g., "processing_error")

## Metric Relationships

- **Incoming Logs** = **Outgoing Logs** + **Dropped Logs**
  - When logs are expanded (e.g., 1 OpenReports log → 3 security events), outgoing count will be higher
  - When logs are dropped, they don't appear in outgoing count

## Example Scenarios

### Scenario 1: Normal Processing
- 10 incoming logs (all pass through unchanged)
- **Result**: `incoming_logs_total = 10`, `outgoing_logs_total = 10`, `dropped_logs_total = 0`

### Scenario 2: OpenReports Expansion
- 1 incoming OpenReports log with 5 results (all pass status filter)
- **Result**: `incoming_logs_total = 1`, `outgoing_logs_total = 5`, `dropped_logs_total = 0`

### Scenario 3: Mixed Processing
- 5 incoming logs: 2 OpenReports (with 3 results each), 3 regular logs
- **Result**: `incoming_logs_total = 5`, `outgoing_logs_total = 9` (6 expanded + 3 unchanged), `dropped_logs_total = 0`

### Scenario 4: With Errors
- 10 incoming logs: 8 successful, 2 errors
- **Result**: `incoming_logs_total = 10`, `outgoing_logs_total = 8`, `dropped_logs_total = 2`, `processing_errors_total = 2`

## Accessing Metrics

Metrics are exposed through the OpenTelemetry Collector's telemetry system and can be:

1. **Exported via Prometheus**: If Prometheus exporter is configured
2. **Viewed in Collector metrics endpoint**: If metrics are enabled
3. **Collected by observability tools**: Any OTLP-compatible metrics receiver

## Configuration

No additional configuration is needed - metrics are automatically enabled when the processor is used.

