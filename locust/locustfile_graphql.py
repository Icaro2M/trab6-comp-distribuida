"""
Teste de carga — GraphQL
========================
Go:    locust -f locustfile_graphql.py --host=http://localhost:8082
Java:  locust -f locustfile_graphql.py --host=http://localhost:8092
"""

import os

from locust import HttpUser, task, between

WAIT_MIN = float(os.getenv("LOCUST_WAIT_MIN", "0.5"))
WAIT_MAX = float(os.getenv("LOCUST_WAIT_MAX", "2"))

HEADERS = {"Content-Type": "application/json"}


class GraphqlUser(HttpUser):
    wait_time = between(WAIT_MIN, WAIT_MAX)

    def gql(self, query: str, name: str):
        self.client.post(
            "/graphql",
            json={"query": query},
            headers=HEADERS,
            name=name
        )

    @task(5)
    def listar_musicas(self):
        self.gql(
            "{ musicas { id nome artista } }",
            "GQL listarMusicas"
        )

    @task(3)
    def listar_playlists(self):
        self.gql(
            "{ playlists { id nome usuario_id } }",
            "GQL listarPlaylists"
        )

    @task(2)
    def listar_usuarios(self):
        self.gql(
            "{ usuarios { id nome idade } }",
            "GQL listarUsuarios"
        )
