#!/bin/bash
# Script to test GitHub Actions workflows locally

set -e

echo "🧪 Testing GitHub Actions Workflows Locally"
echo "==========================================="

# Check if act is installed
if ! command -v act &> /dev/null; then
    echo "❌ 'act' is not installed."
    echo ""
    echo "Install it with:"
    echo "  brew install act  # macOS"
    echo "  # or visit: https://github.com/nektos/act"
    exit 1
fi

echo "✅ act is installed"
echo ""

# Validate YAML syntax
echo "📋 Validating YAML syntax..."
if command -v python3 &> /dev/null; then
    python3 -c "import yaml, sys; [yaml.safe_load(open(f)) for f in sys.argv[1:]]" .github/workflows/*.yml && echo "✅ YAML syntax is valid"
else
    echo "⚠️  Python3 not found, skipping YAML validation"
fi
echo ""

# List available workflows
echo "📝 Available workflows:"
ls -1 .github/workflows/*.yml | sed 's/^/  - /'
echo ""

# Test workflow options
echo "🔧 Usage examples:"
echo ""
echo "1. Test a specific workflow:"
echo "   act -W .github/workflows/ci.yml"
echo ""
echo "2. Test a specific job:"
echo "   act -j build-and-test -W .github/workflows/ci.yml"
echo ""
echo "3. Test with dry-run (no actual execution):"
echo "   act -n -W .github/workflows/ci.yml"
echo ""
echo "4. Test with specific event:"
echo "   act push -W .github/workflows/ci.yml"
echo ""
echo "5. Test pull_request event:"
echo "   act pull_request -W .github/workflows/ci.yml"
echo ""

# Ask user what they want to test
read -p "Would you like to test the simplified workflow now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🚀 Running test workflow..."
    act -W .github/workflows/test-local.yml workflow_dispatch
fi
