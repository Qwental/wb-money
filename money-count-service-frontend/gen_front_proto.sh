#!/bin/bash

# Script to generate protobuf and gRPC-Web code for money_service.proto

# Exit on any error
set -e

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed. Please install Protocol Buffers (e.g., 'brew install protobuf' on macOS)."
    exit 1
fi

# Check if protoc-gen-grpc-web is installed
if ! command -v protoc-gen-grpc-web &> /dev/null; then
    echo "Error: protoc-gen-grpc-web is not installed. Install it from https://github.com/grpc/grpc-web/releases."
    exit 1
fi

# Ensure gen folder exists
mkdir -p gen

# Run protoc to generate JavaScript and gRPC-Web code
protoc \
  -I=../money-count-service/proto \
  --js_out=import_style=commonjs:gen \
  --grpc-web_out=import_style=commonjs,mode=grpcwebtext:gen \
  ../money-count-service/proto/money_service.proto

echo "Successfully generated protobuf files in gen/:
- money_service_pb.js
- money_service_grpc_web_pb.js"