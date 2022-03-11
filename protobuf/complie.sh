#!/usr/bin/env bash

outDir="gen"

if [ ! -d "$outDir" ]; then
  mkdir $outDir
fi

protoc --proto_path=./rpc/ --proto_path=./im/ --go_out=./../../ ./rpc/*.proto ./im/*.proto
