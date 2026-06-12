#!/bin/bash

set -e

echo "Building app..."

mkdir -p bin
go build -o bin/panaszlada

echo "Build completed."
echo "Starting app..."

exec ./bin/panaszlada