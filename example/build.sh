#!/bin/bash

# 自动生成并构建
echo "Generating annotations..."
go generate ./...

echo "Building..."
go build -o app .

echo "Done!"