#!/usr/bin/env bash

VER=3.11.2 

wget https://github.com/protocolbuffers/protobuf/releases/download/v${VER}/protoc-${VER}-linux-x86_64.zip -O protoc.zip
unzip -o protoc.zip -d protocfiles

cp protocfiles/bin/protoc $GOPATH/bin/protoc
rm -rf protoc.zip protocfiles

