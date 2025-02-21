#!/bin/bash
set -e

APP_NAME="reddittui"
BUILD_DIR="build"
GO_MAIN_FILE="main.go"
INSTALL_DIR="/usr/local/bin"

# Build reddittui
echo "Building reddittui application..."
mkdir -p "$BUILD_DIR"
go build -o "$BUILD_DIR/$APP_NAME" "$GO_MAIN_FILE"

# Install reddittui
echo "Installing reddittui..."
echo "Copying binary to $INSTALL_DIR (may require sudo)..."
sudo install -m 0755 "$BUILD_DIR/$APP_NAME" "$INSTALL_DIR/$APP_NAME"

echo "Installation complete. You can now run $APP_NAME' from your terminal."
