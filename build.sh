#!/bin/bash

# Specify the Go source file and the output directory
SOURCE_FILE="main.go"

# List of target operating systems and architectures
# You can customize this list as needed
TARGETS=(
    "linux/amd64"
    "linux/386"
)

# Build the binaries for each target
for target in "${TARGETS[@]}"; do
    # Extract the OS and architecture from the target
    OS=$(echo "$target" | cut -d'/' -f1)
    ARCH=$(echo "$target" | cut -d'/' -f2)
    
    # Set the output binary name
    BINARY_NAME="API_${OS}_${ARCH}"

    # Build the binary
    GOOS="$OS" GOARCH="$ARCH" go build -o "$BINARY_NAME" "$SOURCE_FILE"
    
    if [ $? -eq 0 ]; then
        echo "Built $BINARY_NAME"
    else
        echo "Failed to build $BINARY_NAME"
    fi
done
