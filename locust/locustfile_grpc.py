"""
Teste de carga — gRPC
=====================
Pré-requisito: gerar os stubs Python primeiro:
  bash gerar_proto.sh

Java:  GRPC_HOST=localhost GRPC_PORT=8093 locust -f locustfile_grpc.py --headless -u 10 -r 2 -t 60s --csv=resultados_grpc_java
Go:    GRPC_HOST=localhost GRPC_PORT=8083 locust -f locustfile_grpc.py --headless -u 10 -r 2 -t 60s --csv=resultados_grpc_go

Nota: gRPC não usa HTTP, por isso --host é ignorado.
      O host/porta são controlados pelas variáveis de ambiente acima.
"""

import os
import random
import time

# Necessário para compatibilidade entre gRPC e gevent (usado pelo Locust)
import grpc.experimental.gevent as grpc_gevent
grpc_gevent.init_gevent()

import grpc
from locust import User, task, between, events

import music_pb2
import music_pb2_grpc

GRPC_HOST = os.getenv("GRPC_HOST", "localhost")
GRPC_PORT = int(os.getenv("GRPC_PORT", "8093"))

MUSICA_IDS   = list(range(1, 101))
PLAYLIST_IDS = list(range(1, 51))
USUARIO_IDS  = list(range(1, 21))


class GrpcUser(User):
    wait_time = between(0.5, 2)

    def on_start(self):
        self.channel = grpc.insecure_channel(f"{GRPC_HOST}:{GRPC_PORT}")
        self.stub = music_pb2_grpc.MusicServiceStub(self.channel)

    def on_stop(self):
        self.channel.close()

    def _chamar(self, nome: str, func, request):
        """Executa a chamada gRPC e reporta o tempo ao Locust."""
        inicio = time.perf_counter()
        try:
            resposta = func(request)
            elapsed = (time.perf_counter() - inicio) * 1000
            self.environment.events.request.fire(
                request_type="gRPC",
                name=nome,
                response_time=elapsed,
                response_length=0,
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

    # ── Músicas ──────────────────────────────────────────────────────────────

    @task(8)
    def listar_musicas(self):
        self._chamar(
            "ListMusics",
            self.stub.ListMusics,
            music_pb2.Empty()
        )

    @task(6)
    def buscar_musica(self):
        self._chamar(
            "GetMusic",
            self.stub.GetMusic,
            music_pb2.IdRequest(id=random.choice(MUSICA_IDS))
        )

    @task(2)
    def criar_musica(self):
        self._chamar(
            "CreateMusic",
            self.stub.CreateMusic,
            music_pb2.CreateMusicRequest(nome="Carga", artista="Locust")
        )

    # ── Playlists ─────────────────────────────────────────────────────────────

    @task(5)
    def listar_musicas_playlist(self):
        self._chamar(
            "ListPlaylistMusics",
            self.stub.ListPlaylistMusics,
            music_pb2.PlaylistMusicsRequest(playlist_id=random.choice(PLAYLIST_IDS))
        )

    @task(3)
    def listar_playlists(self):
        self._chamar(
            "ListPlaylists",
            self.stub.ListPlaylists,
            music_pb2.Empty()
        )

    # ── Utilizadores ──────────────────────────────────────────────────────────

    @task(2)
    def listar_usuarios(self):
        self._chamar(
            "ListUsers",
            self.stub.ListUsers,
            music_pb2.Empty()
        )

    @task(1)
    def buscar_usuario(self):
        self._chamar(
            "GetUser",
            self.stub.GetUser,
            music_pb2.IdRequest(id=random.choice(USUARIO_IDS))
        )
