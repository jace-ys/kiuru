#!/usr/bin/env bash

LANG=""
PROTOPATH=""
TARGETS=""

while getopts "l:p:t:" opt; do
  case $opt in
    l )
      LANG=$OPTARG
      ;;
    p )
      PROTOPATH=$OPTARG
      ;;
    t )
      TARGETS=$OPTARG
      ;;
  esac
done

if [ "$LANG" = "" ] || [ "$PROTOPATH" = "" ] || [ "$TARGETS" = "" ]; then
  echo "Usage: protogen.sh [-l] [-d] [-t]"
  exit 1
fi

IFS=','
read -ra ADDR <<< "$TARGETS"
mkdir -p api/grpc
for target in "${ADDR[@]}"; do
  protoc --proto_path=${PROTOPATH} --${LANG}_out=plugins=grpc:api/grpc ${PROTOPATH}/${target}.proto
done
