#!/bin/bash
PROTO_DIR="../ndk/"
GO_OUT_DIR="../pkg/ndk"

go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
# Also need to install protoc
go mod download google.golang.org/grpc/cmd/protoc-gen-go-grpc

mkdir -p $GO_OUT_DIR
for f in ${PROTO_DIR}/*.proto ; do
    # Generate Go bindings
    protoc --proto_path=$PROTO_DIR --go-grpc_out=$GO_OUT_DIR --go_out=$GO_OUT_DIR $f 
done
