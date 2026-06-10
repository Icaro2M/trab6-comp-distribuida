"""
Teste de carga — SOAP GO
=====================
"""

from locust import HttpUser, task, between, constant

def envelope_go(operacao: str) -> str:
    # O Go exige que a tag seja enviada com a primeira letra maiúscula
    # e com a palavra "Request" no final. Ex: "ListarMusicasRequest"
    return (
        '<?xml version="1.0" encoding="UTF-8"?>'
        '<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">'
        "<soapenv:Body>"
        f"<{operacao}></{operacao}>"
        "</soapenv:Body>"
        "</soapenv:Envelope>"
    )

HEADERS = {"Content-Type": "text/xml", "SOAPAction": '""'}

class SoapUser(HttpUser):
    wait_time = constant(0.1)

    @task(5)
    def listar_musicas(self):
        self.client.post(
            "/soap",
            # Enviando a tag exatamente como o Go espera ler
            data=envelope_go("ListarMusicasRequest"), 
            headers=HEADERS,
            name="SOAP listarMusicas"
        )

    @task(3)
    def listar_playlists(self):
        self.client.post(
            "/soap",
            # Enviando a tag exatamente como o Go espera ler
            data=envelope_go("ListarPlaylistsRequest"),
            headers=HEADERS,
            name="SOAP listarPlaylists"
        )

    @task(2)
    def listar_usuarios(self):
        self.client.post(
            "/soap",
            # Enviando a tag exatamente como o Go espera ler
            data=envelope_go("ListarUsuariosRequest"),
            headers=HEADERS,
            name="SOAP listarUsuarios"
        )