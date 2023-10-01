#!/bin/bash

# Check if protoc is installed
if ! [ -x "$(command -v protoc)" ]; then
    echo "Installing protoc"
    # Check has sudo
    if ! [ -x "$(command -v sudo)" ]; then
        apt-get install -y protobuf-compiler
    else
        sudo apt-get install -y protobuf-compiler
    fi
else
    echo "protoc is already installed"
fi