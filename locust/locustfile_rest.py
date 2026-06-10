"""
Teste de carga — REST
=====================
Go:    locust -f locustfile_rest.py --host=http://localhost:8080
Java:  locust -f locustfile_rest.py --host=http://localhost:8090
"""

from locust import HttpUser, task, between, constant
import random


class RestUser(HttpUser):
    wait_time = constant(0.1)

    musica_ids   = []
    playlist_ids = []
    usuario_ids  = []

    def on_start(self):
        # Busca IDs reais do banco via API
        r = self.client.get("/musicas", name="[setup] GET /musicas")
        if r.status_code == 200:
            data = r.json()
            self.musica_ids = [m["id"] for m in data] if data else list(range(1, 101))

        r = self.client.get("/playlists", name="[setup] GET /playlists")
        if r.status_code == 200:
            data = r.json()
            self.playlist_ids = [p["id"] for p in data] if data else list(range(1, 201))

        r = self.client.get("/usuarios", name="[setup] GET /usuarios")
        if r.status_code == 200:
            data = r.json()
            self.usuario_ids = [u["id"] for u in data] if data else list(range(1, 101))

    @task(5)
    def listar_musicas(self):
        self.client.get("/musicas", name="GET /musicas")

    @task(3)
    def listar_playlists(self):
        self.client.get("/playlists", name="GET /playlists")

    @task(2)
    def listar_usuarios(self):
        self.client.get("/usuarios", name="GET /usuarios")
