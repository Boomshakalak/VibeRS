#!/bin/bash

# VibeRS Code Linting Script

set -e

echo "🔍 Running Go code checks..."

# Go vet
echo "📋 Running go vet..."
go vet ./...

# Go fmt check
echo "📝 Checking go fmt..."
UNFORMATTED=$(gofmt -l .)
if [ ! -z "$UNFORMATTED" ]; then
    echo "❌ The following files need formatting:"
    echo "$UNFORMATTED"
    echo "Run: gofmt -w ."
    exit 1
fi

# Go mod tidy check
echo "📦 Checking go mod tidy..."
go mod tidy
if ! git diff --exit-code go.mod go.sum; then
    echo "❌ go.mod or go.sum is not tidy. Run: go mod tidy"
    exit 1
fi

# Test compilation
echo "🔨 Testing compilation..."
go build ./...

echo "✅ All checks passed!" 