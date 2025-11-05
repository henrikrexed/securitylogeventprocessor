# Testing GitHub Actions Workflows Locally

This guide explains how to test GitHub Actions workflows before pushing to GitHub.

## Prerequisites

### Install `act` (Recommended)

`act` is a tool that runs GitHub Actions locally using Docker.

**macOS:**
```bash
brew install act
```

**Linux/Other:**
```bash
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
```

**Manual installation:**
Visit https://github.com/nektos/act/releases

### Install Docker

`act` requires Docker to run. Make sure Docker is installed and running.

**macOS:**
```bash
brew install --cask docker
```

## Quick Start

### 1. Test the Simplified Workflow

```bash
./scripts/test-workflows.sh
```

Or manually:
```bash
act -W .github/workflows/test-local.yml workflow_dispatch
```

### 2. Test Specific CI Jobs

Test individual jobs from the main CI workflow:

```bash
# Test build and test job
act -j build-and-test -W .github/workflows/ci.yml

# Test code quality job
act -j code-quality -W .github/workflows/ci.yml

# Test security scan job
act -j security-scan -W .github/workflows/ci.yml
```

### 3. Test with Different Events

```bash
# Test push event
act push -W .github/workflows/ci.yml

# Test pull_request event
act pull_request -W .github/workflows/ci.yml
```

### 4. Dry Run (Validate Only)

Test workflow syntax without executing:

```bash
act -n -W .github/workflows/ci.yml
```

## Manual Testing Steps

### Validate YAML Syntax

```bash
# Using Python
python3 -c "import yaml; [yaml.safe_load(open(f)) for f in ['.github/workflows/ci.yml', '.github/workflows/release.yml']]"

# Using yamllint (if installed)
yamllint .github/workflows/*.yml
```

### Test Individual Commands

You can run the commands manually that the workflow would run:

```bash
# Set up Go environment (if needed)
export GO_VERSION=1.24

# Run tests
go test -v -race -coverprofile=coverage.out ./...

# Build
go build ./...

# Run linters
go vet ./...
go fmt -l .

# Run security scanner (if installed)
gosec ./...
govulncheck ./...
```

### Test Docker Build Locally

```bash
# Build using the same command as CI
make docker-build DOCKER_BIN=docker RELEASE=0.1 PLATFORM=linux/amd64
```

## Common `act` Options

| Option | Description |
|--------|-------------|
| `-W <workflow>` | Specify workflow file |
| `-j <job>` | Run specific job only |
| `-n, --dry-run` | Dry run (don't execute) |
| `-l, --list` | List all workflows |
| `-v, --verbose` | Verbose output |
| `--secret <name>=<value>` | Set secrets |
| `--env <name>=<value>` | Set environment variables |
| `--eventpath <file>` | Use custom event JSON |

## Limitations

`act` has some limitations compared to real GitHub Actions:

1. **Secrets**: You need to manually provide secrets
   ```bash
   act --secret GITHUB_TOKEN=your_token
   ```

2. **Docker-in-Docker**: Docker builds may not work exactly the same

3. **Matrix builds**: Matrix strategies may behave differently

4. **Service containers**: Limited support for service containers

5. **Artifacts**: Upload/download behavior may differ

## Troubleshooting

### Issue: `act` can't find Docker

**Solution**: Ensure Docker is running
```bash
docker ps
```

### Issue: Permission errors

**Solution**: Run with appropriate permissions or use Docker group
```bash
sudo usermod -aG docker $USER
```

### Issue: Out of disk space

**Solution**: Clean up Docker images
```bash
docker system prune -a
```

### Issue: Workflow fails in `act` but works on GitHub

This is normal - some actions behave differently. Focus on testing:
- YAML syntax
- Basic command execution
- Logic flow

The full integration should still be tested on GitHub.

## Resources

- [act GitHub Repository](https://github.com/nektos/act)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [act Usage Examples](https://github.com/nektos/act#example-commands)
