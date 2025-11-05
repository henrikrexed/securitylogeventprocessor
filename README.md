# Security Log Event Processor

An OpenTelemetry Collector processor that transforms OpenTelemetry logs into security event logs conforming to a standardized security event schema.

## Overview

This processor processes OpenTelemetry logs and transforms them into security events with predefined processing types. Currently supports:

- **OpenReports**: Transforms OpenReports logs into security events

## Architecture

The processor is designed to be extensible, allowing new processing types to be added for different log sources while maintaining a consistent security event schema output.

## Structure

```
├── processor/              # Processor implementations
│   ├── securityevent/     # Main processor
│   └── openreports/        # OpenReports-specific transformation
├── schema/                # Security event schema definitions
└── internal/              # Internal utilities
```

## Usage

### Configuration Example

```yaml
processors:
  securityevent:
    processors:
      openreports:
        enabled: true
        # Optional: Filter which result statuses to process
        # Only results with these statuses will be transformed into security events
        # Valid values: "pass", "fail", "error", "skip"
        # If not specified or empty, all statuses will be processed
        status_filter:
          - "fail"
          - "error"
        # Example: Only process failures and errors, skip "pass" and "skip" results
```

#### Status Filter Options

The `status_filter` configuration allows you to control which OpenReports result statuses are transformed into security events:

- **"pass"**: Policy checks that passed
- **"fail"**: Policy checks that failed (violations)
- **"error"**: Policy checks that encountered errors
- **"skip"**: Policy checks that were skipped

If `status_filter` is not specified or is empty, all statuses will be processed. This is useful for:
- Only tracking violations (`["fail"]`)
- Tracking violations and errors (`["fail", "error"]`)
- Excluding skipped checks (`["pass", "fail", "error"]`)

## Development

### Prerequisites

- Go 1.21 or later
- OpenTelemetry Collector

### Building

```bash
go mod tidy
go build ./...
```

### Testing

```bash
go test ./...
```

