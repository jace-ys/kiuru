#!/usr/bin/env bash

PROTOC_VERSION=3.8.0
curl -OL https://github.com/google/protobuf/releases/download/v$PROTOC_VERSION/protoc-$PROTOC_VERSION-linux-x86_64.zip
sudo unzip protoc-$PROTOC_VERSION-linux-x86_64.zip -d /usr/local
rm protoc-$PROTOC_VERSION-linux-x86_64.zip

case "$1" in
  go )
    GO111MODULE=off go get github.com/golang/protobuf/protoc-gen-go
    GO111MODULE=off go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
    ;;
esac
