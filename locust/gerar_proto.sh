#!/bin/bash
# Gera os stubs Python a partir do music.proto
# Necessário apenas para o teste gRPC
python3 -m grpc_tools.protoc \
  -I../go/grpc/proto \
  --python_out=. \
  --grpc_python_out=. \
  ../go/grpc/proto/music.proto

echo "Stubs gerados: music_pb2.py e music_pb2_grpc.py"
