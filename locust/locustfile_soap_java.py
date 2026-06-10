"""
Teste de carga — SOAP
=====================
Go:    locust -f locustfile_soap.py --host=http://localhost:8081
Java:  locust -f locustfile_soap.py --host=http://localhost:8091
"""

from locust import HttpUser, task, between, constant
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
    wait_time = constant(0.1)

    @task(5)
    def listar_musicas(self):
        self.client.post(
            "/soap",
            data=envelope("listarMusicas", ""),
            headers=HEADERS,
            name="SOAP listarMusicas"
        )

    @task(3)
    def listar_playlists(self):
        self.client.post(
            "/soap",
            data=envelope("listarPlaylists", ""),
            headers=HEADERS,
            name="SOAP listarPlaylists"
        )

    @task(2)
    def listar_usuarios(self):
        self.client.post(
            "/soap",
            data=envelope("listarUsuarios", ""),
            headers=HEADERS,
            name="SOAP listarUsuarios"
        )
