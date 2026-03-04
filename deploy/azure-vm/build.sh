#!/bin/bash
set -e

# ============================================
# Gophish Build Script for Linux AMD64
# ============================================

# This script is intended to run on a Linux machine with Go installed.
# If running on Mac, the build may fail due to CGO requirements.

echo "========================================="
echo "  Building Gophish for Linux AMD64"
echo "========================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    echo "Please install Go 1.24+: https://go.dev/dl/"
    exit 1
fi

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo "Project directory: $PROJECT_DIR"
echo "Go version: $(go version)"

# Clean previous build
echo "Cleaning previous build..."
rm -f gophish

# Ensure dependencies are up to date
echo "Updating Go dependencies..."
go mod download
go mod tidy

# Check if we can build with CGO (Linux only)
if [ "$(uname)" = "Linux" ]; then
    echo "Building Gophish for Linux AMD64 (with CGO)..."
    GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o gophish .
else
    echo "Warning: Building on Mac may not work properly due to CGO."
    echo "Try: docker run --rm -v $(pwd):/workspace -w /workspace golang:1.24 bash -c 'GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o gophish .'"
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gophish . || true
fi

# Verify build
if [ -f "./gophish" ]; then
    echo "Build successful!"
    file ./gophish
    ls -lh ./gophish
else
    echo "Build failed!"
    exit 1
fi

echo ""
echo "========================================="
echo "  Build Complete!"
echo "========================================="
echo ""
echo "To deploy to Azure VM:"
echo "  1. Copy 'gophish' binary to your VM"
echo "  2. Copy 'deploy/azure-vm/install.sh' to your VM"
echo "  3. Run: sudo bash install.sh"
echo ""
