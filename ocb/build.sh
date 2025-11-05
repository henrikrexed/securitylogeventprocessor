#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building OpenTelemetry Collector with Security Event Processor${NC}"

# Check if OCB is installed
if ! command -v otelcolbuilder &> /dev/null; then
    echo -e "${YELLOW}otelcolbuilder not found. Installing...${NC}"
    go install go.opentelemetry.io/collector/cmd/builder@latest
    export PATH=$PATH:$(go env GOPATH)/bin
fi

# Clean previous builds
echo -e "${YELLOW}Cleaning previous builds...${NC}"
rm -rf dist build

# Build with OCB
MANIFEST=${1:-manifest.yaml}
echo -e "${GREEN}Building collector with manifest: ${MANIFEST}${NC}"

otelcolbuilder --manifest ${MANIFEST} --output-path ./dist

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Collector built successfully!${NC}"
    echo -e "${GREEN}Binary location: ./dist/$(ls dist/ | grep otelcol-securityevent | head -1)${NC}"
else
    echo -e "${RED}✗ Build failed!${NC}"
    exit 1
fi

