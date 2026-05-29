"""
Teste de carga — REST
=====================
Java:  locust -f locustfile_rest.py --host=http://localhost:8090
Go:    locust -f locustfile_rest.py --host=http://localhost:8080

Modo headless (sem interface):
  locust -f locustfile_rest.py --host=http://localhost:8090 \
         --headless -u 10 -r 2 -t 60s --csv=resultados_rest_java

Descrição dos pesos das tarefas:
  listar_musicas       → 40% dos pedidos (leitura mais comum)
  buscar_musica_por_id → 30% (leitura por id)
  listar_musicas_playlist → 15% (query com join)
  listar_usuarios      → 10%
  criar_musica         →  5% (escrita)
"""

from locust import HttpUser, task, between
import random


class RestUser(HttpUser):
    wait_time = between(0.5, 2)

    # IDs presentes na carga 1
    MUSICA_IDS  = list(range(1, 101))
    PLAYLIST_IDS = list(range(1, 51))
    USUARIO_IDS = list(range(1, 21))

    # ── Músicas ──────────────────────────────────────────────────────────────

    @task(8)
    def listar_musicas(self):
        self.client.get("/musicas", name="GET /musicas")

    @task(6)
    def buscar_musica_por_id(self):
        id = random.choice(self.MUSICA_IDS)
        self.client.get(f"/musicas/{id}", name="GET /musicas/{id}")

    @task(2)
    def criar_musica(self):
        self.client.post(
            "/musicas",
            json={"nome": "Carga Locust", "artista": "Teste"},
            name="POST /musicas"
        )

    # ── Playlists ─────────────────────────────────────────────────────────────

    @task(5)
    def listar_musicas_playlist(self):
        id = random.choice(self.PLAYLIST_IDS)
        self.client.get(f"/playlists/{id}/musicas", name="GET /playlists/{id}/musicas")

    @task(3)
    def listar_playlists(self):
        self.client.get("/playlists", name="GET /playlists")

    # ── Utilizadores ──────────────────────────────────────────────────────────

    @task(2)
    def listar_usuarios(self):
        self.client.get("/usuarios", name="GET /usuarios")

    @task(1)
    def buscar_usuario_por_id(self):
        id = random.choice(self.USUARIO_IDS)
        self.client.get(f"/usuarios/{id}", name="GET /usuarios/{id}")
