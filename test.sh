#!/bin/bash
set -e

echo "Running local tests inside Docker using Go 1.26..."
docker build -t eightbyeight-test -f Dockerfile.test .
echo "Local Docker tests passed!"
