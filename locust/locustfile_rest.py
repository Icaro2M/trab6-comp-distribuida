"""
Teste de carga — REST
=====================
Go:    locust -f locustfile_rest.py --host=http://localhost:8080
Java:  locust -f locustfile_rest.py --host=http://localhost:8090
"""

import os

from locust import HttpUser, task, between

WAIT_MIN = float(os.getenv("LOCUST_WAIT_MIN", "0.5"))
WAIT_MAX = float(os.getenv("LOCUST_WAIT_MAX", "2"))


class RestUser(HttpUser):
    wait_time = between(WAIT_MIN, WAIT_MAX)

    @task(5)
    def listar_musicas(self):
        self.client.get("/musicas", name="GET /musicas")

    @task(3)
    def listar_playlists(self):
        self.client.get("/playlists", name="GET /playlists")

    @task(2)
    def listar_usuarios(self):
        self.client.get("/usuarios", name="GET /usuarios")
