# Overview

The Security Log Event Processor is an OpenTelemetry Collector processor that transforms security-related logs into standardized security event logs. This document provides an overview of the processor and its capabilities.

## What is the Security Log Event Processor?

The Security Log Event Processor is a custom processor for the OpenTelemetry Collector that:

- **Transforms** security logs from various sources into a standardized security event format
- **Expands** single log records into multiple security events (e.g., one OpenReports log with multiple results becomes multiple security events)
- **Enriches** events with Kubernetes metadata and workload information
- **Filters** events based on configurable criteria (e.g., status, severity)

## Architecture

```
┌─────────────┐
│  Receivers  │  (k8sobjects, OTLP, etc.)
└──────┬──────┘
       │
       ▼
┌─────────────────────┐
│ Security Event      │
│ Processor           │
│  ├─ OpenReports     │
│  └─ (more coming)   │
└──────┬──────────────┘
       │
       ▼
┌─────────────┐
│  Exporters  │  (OTLP, logging, etc.)
└─────────────┘
```

## Key Concepts

### Processors

Processors are the core components that transform logs. Each processor handles a specific log format:

- **OpenReports Processor**: Processes OpenReports Custom Resource logs from Kubernetes

### Receivers

Receivers collect data from various sources. The processor you choose determines which receiver you need:

- **k8sobjects receiver**: Required for OpenReports processor (collects Kubernetes Custom Resources)
- **OTLP receiver**: Optional, for receiving logs via OTLP protocol
- **filelog receiver**: Optional, for reading logs from files

### Security Events

The processor outputs standardized security events with fields like:
- `event.id`: Unique event identifier
- `event.type`: Type of security event
- `dt.security.risk.level`: Risk level (CRITICAL, HIGH, MEDIUM, LOW)
- `compliance.status`: Compliance status (COMPLIANT, NON_COMPLIANT)
- `k8s.*`: Kubernetes metadata (workload, namespace, etc.)

## Use Cases

1. **Compliance Monitoring**: Transform policy evaluation results into security events
2. **Security Event Correlation**: Standardize security logs from multiple sources
3. **Kubernetes Security**: Monitor and alert on Kubernetes security events
4. **Audit Logging**: Create detailed audit trails from security logs

## Next Steps

- [Installation Guide](installation.md): Learn how to install the processor
- [Quick Start](quick-start.md): Get started in 5 minutes
- [Configuration](configuration/processor-config.md): Configure the processor

