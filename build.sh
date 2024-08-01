#!/bin/sh

VERSION=1.0.0
BINARY_NAME=dbcli

build_for_platform() {
    GOOS=$1
    GOARCH=$2
    output_name=$BINARY_NAME'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name=$output_name'.exe'
    fi

    echo "Building for $GOOS $GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name ./cmd
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
}

# Build for Windows
build_for_platform windows amd64

# Build for macOS
build_for_platform darwin amd64

# Build for Linux
build_for_platform linux amd64

echo "Build process completed."