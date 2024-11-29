#!/bin/bash

# Set script to exit immediately on error
set -e

# Define paths
PROJECT_ROOT="$(pwd)"
LAMBDA_SRC_DIR="src/processorlambda"
OUTPUT_DIR="${PROJECT_ROOT}/build"
LAMBDA_BINARY="${OUTPUT_DIR}/processorlambda"
ZIP_FILE="${OUTPUT_DIR}/lambda.zip"

# Ensure build directory exists
mkdir -p "$OUTPUT_DIR"

# Navigate to the Lambda source directory and build
echo "Compiling Go Lambda binary..."
cd "$LAMBDA_SRC_DIR"
GOOS=linux GOARCH=amd64 go build -o "$LAMBDA_BINARY"

# Navigate back to project root
cd "$PROJECT_ROOT"

# Create the zip file
echo "Creating zip file..."
zip -j "$ZIP_FILE" "$LAMBDA_BINARY"

echo "Lambda zip file created: $ZIP_FILE"

