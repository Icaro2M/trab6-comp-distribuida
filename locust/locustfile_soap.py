"""
Teste de carga — SOAP
=====================
Java:  locust -f locustfile_soap.py --host=http://localhost:8091
Go:    locust -f locustfile_soap.py --host=http://localhost:8081

Modo headless:
  locust -f locustfile_soap.py --host=http://localhost:8091 \
         --headless -u 10 -r 2 -t 60s --csv=resultados_soap_java
"""

from locust import HttpUser, task, between
import random

NS = "http://service.soap.musicaservice.com/"

def envelope(operacao: str, corpo: str) -> str:
    return (
        '<?xml version="1.0" encoding="UTF-8"?>'
        '<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"'
        f' xmlns:ser="{NS}">'
        "<soapenv:Body>"
        f"<ser:{operacao}>{corpo}</ser:{operacao}>"
        "</soapenv:Body>"
        "</soapenv:Envelope>"
    )

HEADERS = {"Content-Type": "text/xml;charset=UTF-8", "SOAPAction": '""'}


class SoapUser(HttpUser):
    wait_time = between(0.5, 2)

    MUSICA_IDS   = list(range(1, 101))
    PLAYLIST_IDS = list(range(1, 51))
    USUARIO_IDS  = list(range(1, 21))

    # ── Músicas ──────────────────────────────────────────────────────────────

    @task(8)
    def listar_musicas(self):
        self.client.post(
            "/soap",
            data=envelope("listarMusicas", ""),
            headers=HEADERS,
            name="SOAP listarMusicas"
        )

    @task(6)
    def buscar_musica(self):
        id = random.choice(self.MUSICA_IDS)
        self.client.post(
            "/soap",
            data=envelope("buscarMusica", f"<id>{id}</id>"),
            headers=HEADERS,
            name="SOAP buscarMusica"
        )

    @task(2)
    def criar_musica(self):
        self.client.post(
            "/soap",
            data=envelope("criarMusica", "<nome>Carga</nome><artista>Locust</artista>"),
            headers=HEADERS,
            name="SOAP criarMusica"
        )

    # ── Playlists ─────────────────────────────────────────────────────────────

    @task(5)
    def listar_musicas_playlist(self):
        id = random.choice(self.PLAYLIST_IDS)
        self.client.post(
            "/soap",
            data=envelope("listarMusicasPlaylist", f"<playlist_id>{id}</playlist_id>"),
            headers=HEADERS,
            name="SOAP listarMusicasPlaylist"
        )

    @task(3)
    def listar_playlists(self):
        self.client.post(
            "/soap",
            data=envelope("listarPlaylists", ""),
            headers=HEADERS,
            name="SOAP listarPlaylists"
        )

    # ── Utilizadores ──────────────────────────────────────────────────────────

    @task(2)
    def listar_usuarios(self):
        self.client.post(
            "/soap",
            data=envelope("listarUsuarios", ""),
            headers=HEADERS,
            name="SOAP listarUsuarios"
        )

    @task(1)
    def buscar_usuario(self):
        id = random.choice(self.USUARIO_IDS)
        self.client.post(
            "/soap",
            data=envelope("buscarUsuario", f"<id>{id}</id>"),
            headers=HEADERS,
            name="SOAP buscarUsuario"
        )
