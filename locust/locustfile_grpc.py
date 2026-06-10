"""
Teste de carga â€” gRPC
=====================
Pré-requisito: gerar os stubs Python primeiro:
  python -m grpc_tools.protoc -I../go/grpc/proto --python_out=. --grpc_python_out=. ../go/grpc/proto/music.proto

Go:    GRPC_HOST=localhost GRPC_PORT=8083 locust -f locustfile_grpc.py --headless -u 50 -r 5 -t 60s --csv=resultados_grpc_go
Java:  GRPC_HOST=localhost GRPC_PORT=8093 locust -f locustfile_grpc.py --headless -u 50 -r 5 -t 60s --csv=resultados_grpc_java
"""

import os
import time

import grpc.experimental.gevent as grpc_gevent
grpc_gevent.init_gevent()

import grpc
from locust import User, task, between, constant

import music_pb2
import music_pb2_grpc

GRPC_HOST = os.getenv("GRPC_HOST", "localhost")
GRPC_PORT = int(os.getenv("GRPC_PORT", "8083"))


class GrpcUser(User):
    wait_time = constant(0.1)
    def on_start(self):
        self.channel = grpc.insecure_channel(f"{GRPC_HOST}:{GRPC_PORT}")
        self.stub = music_pb2_grpc.MusicServiceStub(self.channel)

    def on_stop(self):
        self.channel.close()

    def _chamar(self, nome: str, func, request):
        inicio = time.perf_counter()
        try:
            resposta = func(request)
            elapsed = (time.perf_counter() - inicio) * 1000
            response_length = resposta.ByteSize() if resposta is not None else 0
            self.environment.events.request.fire(
                request_type="gRPC",
                name=nome,
                response_time=elapsed,
                response_length=response_length,
                exception=None,
                context={}
            )
            return resposta
        except grpc.RpcError as e:
            elapsed = (time.perf_counter() - inicio) * 1000
            self.environment.events.request.fire(
                request_type="gRPC",
                name=nome,
                response_time=elapsed,
                response_length=0,
                exception=e,
                context={}
            )

    @task(5)
    def listar_musicas(self):
        self._chamar(
            "gRPC ListMusics",
            self.stub.ListMusics,
            music_pb2.Empty()
        )

    @task(3)
    def listar_playlists(self):
        self._chamar(
            "gRPC ListPlaylists",
            self.stub.ListPlaylists,
            music_pb2.Empty()
        )

    @task(2)
    def listar_usuarios(self):
        self._chamar(
            "gRPC ListUsers",
            self.stub.ListUsers,
            music_pb2.Empty()
        )