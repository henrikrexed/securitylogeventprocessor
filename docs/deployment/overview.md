# Deployment Overview

This guide explains deployment options for the Security Log Event Processor.

## Deployment Methods

The Security Log Event Processor can be deployed using:

1. **[Docker](docker.md)**: Containerized deployment for local development and production
2. **[Kubernetes](kubernetes.md)**: Native Kubernetes deployment with full RBAC support

## Choosing a Deployment Method

### Use Docker If:
- Local development and testing
- Simple deployment scenarios
- Single-host deployment
- CI/CD pipelines

### Use Kubernetes If:
- Production deployments
- Need high availability
- Multi-cluster deployments
- Need RBAC integration for k8sobjects receiver

## Common Requirements

Regardless of deployment method, you'll need:

1. **OpenReports CRD**: The OpenReports Custom Resource Definition must be installed
2. **Configuration**: Collector configuration file
3. **RBAC** (for Kubernetes): ClusterRole and ClusterRoleBinding for k8sobjects receiver

## Quick Comparison

| Feature | Docker | Kubernetes |
|---------|--------|------------|
| Setup Complexity | Low | Medium |
| RBAC Support | No | Yes |
| High Availability | Manual | Built-in |
| Scaling | Manual | Automatic (HPA) |
| Resource Management | Manual | Automatic |
| Best For | Development | Production |

## Next Steps

- [Docker Deployment](docker.md): Deploy using Docker
- [Kubernetes Deployment](kubernetes.md): Deploy to Kubernetes
- [RBAC Setup](rbac.md): Configure RBAC permissions

