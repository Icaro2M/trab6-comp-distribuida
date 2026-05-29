"""
Teste de carga — GraphQL
========================
Java:  locust -f locustfile_graphql.py --host=http://localhost:8092
Go:    locust -f locustfile_graphql.py --host=http://localhost:8082

Modo headless:
  locust -f locustfile_graphql.py --host=http://localhost:8092 \
         --headless -u 10 -r 2 -t 60s --csv=resultados_graphql_java
"""

from locust import HttpUser, task, between
import random

HEADERS = {"Content-Type": "application/json"}


class GraphqlUser(HttpUser):
    wait_time = between(0.5, 2)

    MUSICA_IDS   = list(range(1, 101))
    PLAYLIST_IDS = list(range(1, 51))
    USUARIO_IDS  = list(range(1, 21))

    def gql(self, query: str, name: str):
        self.client.post(
            "/graphql",
            json={"query": query},
            headers=HEADERS,
            name=name
        )

    # ── Músicas ──────────────────────────────────────────────────────────────

    @task(8)
    def listar_musicas(self):
        self.gql(
            "{ musicas { id nome artista } }",
            "GQL musicas"
        )

    @task(6)
    def buscar_musica(self):
        id = random.choice(self.MUSICA_IDS)
        self.gql(
            f"{{ musica(id: {id}) {{ id nome artista }} }}",
            "GQL musica(id)"
        )

    @task(2)
    def criar_musica(self):
        self.gql(
            'mutation { criarMusica(nome: "Carga", artista: "Locust") { id } }',
            "GQL criarMusica"
        )

    # ── Playlists ─────────────────────────────────────────────────────────────

    @task(5)
    def musicas_da_playlist(self):
        id = random.choice(self.PLAYLIST_IDS)
        self.gql(
            f"{{ musicasDaPlaylist(playlist_id: {id}) {{ id nome artista }} }}",
            "GQL musicasDaPlaylist"
        )

    @task(3)
    def listar_playlists(self):
        self.gql(
            "{ playlists { id nome usuario_id } }",
            "GQL playlists"
        )

    # ── Utilizadores ──────────────────────────────────────────────────────────

    @task(2)
    def listar_usuarios(self):
        self.gql(
            "{ usuarios { id nome idade } }",
            "GQL usuarios"
        )

    @task(1)
    def buscar_usuario(self):
        id = random.choice(self.USUARIO_IDS)
        self.gql(
            f"{{ usuario(id: {id}) {{ id nome idade }} }}",
            "GQL usuario(id)"
        )
