#!/bin/bash

# Build server
echo "Building server..."
go build -o bin/server cmd/server/main.go

# Build database tools
echo "Building database tools..."
go build -o bin/dbtools cmd/tools/main.go

echo "Build complete!" 