#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo -e "${GREEN}Building Docker image for Security Event Processor${NC}"

# Change to project root
cd "$PROJECT_ROOT"

# Build Docker image (OCB will build the collector inside Docker)
echo -e "${YELLOW}Building Docker image with OCB (collector will be built inside Docker)...${NC}"
docker build -f docker/Dockerfile -t otelcol-securityevent:latest .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Docker image built successfully!${NC}"
    echo -e "${GREEN}Image: otelcol-securityevent:latest${NC}"
    echo -e "${YELLOW}To run: docker run -p 4317:4317 -p 4318:4318 otelcol-securityevent:latest${NC}"
else
    echo -e "${RED}✗ Docker build failed!${NC}"
    exit 1
fi

