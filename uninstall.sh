#!/bin/bash
set -e

APP_NAME="reddittui"
INSTALL_DIR="/usr/local/bin"
BINARY_PATH="$INSTALL_DIR/$APP_NAME"

if [[ -f $BINARY_PATH ]]; then
    echo "Uninstalling reddittui..."
    echo "Removing binary from $INSTALL_DIR (may require sudo)..."
    sudo rm "$BINARY_PATH"
    echo "Uninstallation complete."
else
    echo "reddittui is not installed. Nothing to do."
fi
