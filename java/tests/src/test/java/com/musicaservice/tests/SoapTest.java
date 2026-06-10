package com.musicaservice.tests;

import org.junit.jupiter.api.*;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

import static org.junit.jupiter.api.Assertions.*;

@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
@DisplayName("SOAP — http://HOST:8091/soap")
public class SoapTest {

    private static final String HOST = System.getProperty("test.host", "localhost");
    private static final String ENDPOINT = "http://" + HOST + ":8091/soap";
    private static final HttpClient http = HttpClient.newHttpClient();

    private static int musicaIdCriada;
    private static int usuarioIdCriado;
    private static int playlistIdCriada;

    // ─── Utilitário ─────────────────────────────────────────────────────────

    private String soap(String operacao, String parametros) throws Exception {
        String envelope =
            "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" +
            "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\"" +
            "  xmlns:ser=\"http://service.soap.musicaservice.com/\">" +
            "  <soapenv:Body>" +
            "    <ser:" + operacao + ">" + parametros + "</ser:" + operacao + ">" +
            "  </soapenv:Body>" +
            "</soapenv:Envelope>";

        HttpResponse<String> r = http.send(
            HttpRequest.newBuilder(URI.create(ENDPOINT))
                .header("Content-Type", "text/xml;charset=UTF-8")
                .header("SOAPAction", "\"\"")
                .POST(HttpRequest.BodyPublishers.ofString(envelope)).build(),
            HttpResponse.BodyHandlers.ofString());

        assertEquals(200, r.statusCode(), "SOAP request falhou para operação: " + operacao);
        return r.body();
    }

    private String extrairTag(String xml, String tag) {
        int start = xml.indexOf("<" + tag + ">");
        int end   = xml.indexOf("</" + tag + ">");
        if (start == -1 || end == -1) return "";
        return xml.substring(start + tag.length() + 2, end);
    }

    // ─── WSDL ────────────────────────────────────────────────────────────────

    @Test @Order(1)
    @DisplayName("WSDL disponível em /soap?wsdl")
    void wsdlDisponivel() throws Exception {
        var r = http.send(
            HttpRequest.newBuilder(URI.create(ENDPOINT + "?wsdl")).GET().build(),
            HttpResponse.BodyHandlers.ofString());
        assertEquals(200, r.statusCode());
        assertTrue(r.body().contains("wsdl:definitions"), "Resposta não é um WSDL válido");
        System.out.println("  [OK] WSDL disponível");
    }

    // ─── Músicas ─────────────────────────────────────────────────────────────

    @Test @Order(2)
    @DisplayName("listarMusicas → resposta com elementos")
    void listarMusicas() throws Exception {
        String resp = soap("listarMusicas", "");
        assertTrue(resp.contains("<nome>"), "Resposta não contém músicas");
        System.out.println("  [OK] listarMusicas retornou músicas");
    }

    @Test @Order(3)
    @DisplayName("buscarMusica id=1 → retorna música")
    void buscarMusica() throws Exception {
        String resp = soap("buscarMusica", "<id>1</id>");
        String nome = extrairTag(resp, "nome");
        assertFalse(nome.isEmpty(), "Nome não encontrado na resposta");
        System.out.println("  [OK] Música: " + nome);
    }

    @Test @Order(4)
    @DisplayName("criarMusica → retorna id gerado")
    void criarMusica() throws Exception {
        String resp = soap("criarMusica",
            "<nome>Teste SOAP</nome><artista>Artista SOAP</artista>");
        String idStr = extrairTag(resp, "id");
        assertFalse(idStr.isEmpty(), "ID não retornado");
        musicaIdCriada = Integer.parseInt(idStr);
        assertTrue(musicaIdCriada > 0);
        System.out.println("  [OK] Criada com id=" + musicaIdCriada);
    }

    @Test @Order(5)
    @DisplayName("atualizarMusica → retorna nome actualizado")
    void atualizarMusica() throws Exception {
        String resp = soap("atualizarMusica",
            "<id>" + musicaIdCriada + "</id>" +
            "<nome>Teste SOAP Actualizado</nome>" +
            "<artista>Novo Artista SOAP</artista>");
        assertEquals("Teste SOAP Actualizado", extrairTag(resp, "nome"));
        System.out.println("  [OK] Actualizada");
    }

    @Test @Order(6)
    @DisplayName("deletarMusica → mensagem de sucesso")
    void deletarMusica() throws Exception {
        String resp = soap("deletarMusica", "<id>" + musicaIdCriada + "</id>");
        assertTrue(resp.contains("deletada"), "Mensagem de sucesso não encontrada");
        System.out.println("  [OK] Apagada id=" + musicaIdCriada);
    }

    // ─── Utilizadores ────────────────────────────────────────────────────────

    @Test @Order(7)
    @DisplayName("listarUsuarios → resposta com elementos")
    void listarUsuarios() throws Exception {
        String resp = soap("listarUsuarios", "");
        assertTrue(resp.contains("<nome>"));
        System.out.println("  [OK] listarUsuarios retornou utilizadores");
    }

    @Test @Order(8)
    @DisplayName("criarUsuario → retorna id gerado")
    void criarUsuario() throws Exception {
        String resp = soap("criarUsuario", "<nome>Teste User SOAP</nome><idade>25</idade>");
        usuarioIdCriado = Integer.parseInt(extrairTag(resp, "id"));
        assertTrue(usuarioIdCriado > 0);
        System.out.println("  [OK] Utilizador criado id=" + usuarioIdCriado);
    }

    @Test @Order(9)
    @DisplayName("deletarUsuario → mensagem de sucesso")
    void deletarUsuario() throws Exception {
        String resp = soap("deletarUsuario", "<id>" + usuarioIdCriado + "</id>");
        assertTrue(resp.contains("deletado"));
        System.out.println("  [OK] Utilizador apagado id=" + usuarioIdCriado);
    }

    // ─── Playlists ───────────────────────────────────────────────────────────

    @Test @Order(10)
    @DisplayName("listarPlaylists → resposta com elementos")
    void listarPlaylists() throws Exception {
        String resp = soap("listarPlaylists", "");
        assertTrue(resp.contains("<nome>"));
        System.out.println("  [OK] listarPlaylists retornou playlists");
    }

    @Test @Order(11)
    @DisplayName("criarPlaylist → retorna id gerado")
    void criarPlaylist() throws Exception {
        String resp = soap("criarPlaylist",
            "<nome>Playlist SOAP</nome><usuario_id>1</usuario_id>");
        playlistIdCriada = Integer.parseInt(extrairTag(resp, "id"));
        assertTrue(playlistIdCriada > 0);
        System.out.println("  [OK] Playlist criada id=" + playlistIdCriada);
    }

    @Test @Order(12)
    @DisplayName("adicionarMusicaNaPlaylist → mensagem de sucesso")
    void adicionarMusicaPlaylist() throws Exception {
        String resp = soap("adicionarMusicaNaPlaylist",
            "<playlist_id>" + playlistIdCriada + "</playlist_id><musica_id>1</musica_id>");
        assertTrue(resp.contains("adicionada"));
        System.out.println("  [OK] Música adicionada à playlist");
    }

    @Test @Order(13)
    @DisplayName("listarMusicasPlaylist → músicas da playlist")
    void listarMusicasPlaylist() throws Exception {
        String resp = soap("listarMusicasPlaylist",
            "<playlist_id>" + playlistIdCriada + "</playlist_id>");
        assertTrue(resp.contains("<nome>"));
        System.out.println("  [OK] Músicas da playlist listadas");
    }

    @Test @Order(14)
    @DisplayName("deletarPlaylist → mensagem de sucesso")
    void deletarPlaylist() throws Exception {
        String resp = soap("deletarPlaylist", "<id>" + playlistIdCriada + "</id>");
        assertTrue(resp.contains("deletada"));
        System.out.println("  [OK] Playlist apagada id=" + playlistIdCriada);
    }
}
