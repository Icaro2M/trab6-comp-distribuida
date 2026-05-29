package com.musicaservice.tests;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.*;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

import static org.junit.jupiter.api.Assertions.*;

@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
@DisplayName("REST — http://HOST:8090")
public class RestTest {

    private static final String HOST = System.getProperty("test.host", "localhost");
    private static final String BASE = "http://" + HOST + ":8090";
    private static final HttpClient http = HttpClient.newHttpClient();
    private static final ObjectMapper json = new ObjectMapper();

    private static int musicaIdCriada;
    private static int usuarioIdCriado;
    private static int playlistIdCriada;

    // ─── Utilitários ────────────────────────────────────────────────────────

    private HttpResponse<String> get(String path) throws Exception {
        return http.send(
            HttpRequest.newBuilder(URI.create(BASE + path)).GET().build(),
            HttpResponse.BodyHandlers.ofString());
    }

    private HttpResponse<String> post(String path, String body) throws Exception {
        return http.send(
            HttpRequest.newBuilder(URI.create(BASE + path))
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(body)).build(),
            HttpResponse.BodyHandlers.ofString());
    }

    private HttpResponse<String> put(String path, String body) throws Exception {
        return http.send(
            HttpRequest.newBuilder(URI.create(BASE + path))
                .header("Content-Type", "application/json")
                .PUT(HttpRequest.BodyPublishers.ofString(body)).build(),
            HttpResponse.BodyHandlers.ofString());
    }

    private HttpResponse<String> delete(String path) throws Exception {
        return http.send(
            HttpRequest.newBuilder(URI.create(BASE + path))
                .DELETE().build(),
            HttpResponse.BodyHandlers.ofString());
    }

    // ─── Músicas ─────────────────────────────────────────────────────────────

    @Test @Order(1)
    @DisplayName("GET /musicas → lista não vazia")
    void listarMusicas() throws Exception {
        var r = get("/musicas");
        assertEquals(200, r.statusCode());
        JsonNode body = json.readTree(r.body());
        assertTrue(body.isArray() && body.size() > 0, "Lista de músicas deve ter elementos");
        System.out.println("  [OK] " + body.size() + " músicas encontradas");
    }

    @Test @Order(2)
    @DisplayName("GET /musicas/1 → retorna música com id=1")
    void buscarMusicaPorId() throws Exception {
        var r = get("/musicas/1");
        assertEquals(200, r.statusCode());
        JsonNode body = json.readTree(r.body());
        assertEquals(1, body.get("id").asInt());
        System.out.println("  [OK] Música: " + body.get("nome").asText());
    }

    @Test @Order(3)
    @DisplayName("GET /musicas/999999 → 404")
    void buscarMusicaInexistente() throws Exception {
        var r = get("/musicas/999999");
        assertEquals(404, r.statusCode());
        System.out.println("  [OK] 404 para id inexistente");
    }

    @Test @Order(4)
    @DisplayName("POST /musicas → cria e retorna 201")
    void criarMusica() throws Exception {
        var r = post("/musicas", "{\"nome\":\"Teste REST\",\"artista\":\"Artista Teste\"}");
        assertEquals(201, r.statusCode());
        JsonNode body = json.readTree(r.body());
        musicaIdCriada = body.get("id").asInt();
        assertTrue(musicaIdCriada > 0);
        System.out.println("  [OK] Criada com id=" + musicaIdCriada);
    }

    @Test @Order(5)
    @DisplayName("PUT /musicas/{id} → actualiza e retorna 200")
    void atualizarMusica() throws Exception {
        var r = put("/musicas/" + musicaIdCriada,
            "{\"nome\":\"Teste REST Actualizado\",\"artista\":\"Novo Artista\"}");
        assertEquals(200, r.statusCode());
        JsonNode body = json.readTree(r.body());
        assertEquals("Teste REST Actualizado", body.get("nome").asText());
        System.out.println("  [OK] Actualizada: " + body.get("nome").asText());
    }

    @Test @Order(6)
    @DisplayName("DELETE /musicas/{id} → apaga e retorna 200")
    void deletarMusica() throws Exception {
        var r = delete("/musicas/" + musicaIdCriada);
        assertEquals(200, r.statusCode());
        // Confirma que já não existe
        assertEquals(404, get("/musicas/" + musicaIdCriada).statusCode());
        System.out.println("  [OK] Apagada id=" + musicaIdCriada);
    }

    // ─── Utilizadores ────────────────────────────────────────────────────────

    @Test @Order(7)
    @DisplayName("GET /usuarios → lista não vazia")
    void listarUsuarios() throws Exception {
        var r = get("/usuarios");
        assertEquals(200, r.statusCode());
        JsonNode body = json.readTree(r.body());
        assertTrue(body.isArray() && body.size() > 0);
        System.out.println("  [OK] " + body.size() + " utilizadores encontrados");
    }

    @Test @Order(8)
    @DisplayName("POST /usuarios → cria utilizador")
    void criarUsuario() throws Exception {
        var r = post("/usuarios", "{\"nome\":\"Teste User\",\"idade\":22}");
        assertEquals(201, r.statusCode());
        usuarioIdCriado = json.readTree(r.body()).get("id").asInt();
        System.out.println("  [OK] Utilizador criado id=" + usuarioIdCriado);
    }

    @Test @Order(9)
    @DisplayName("DELETE /usuarios/{id} → apaga utilizador")
    void deletarUsuario() throws Exception {
        assertEquals(200, delete("/usuarios/" + usuarioIdCriado).statusCode());
        System.out.println("  [OK] Utilizador apagado id=" + usuarioIdCriado);
    }

    // ─── Playlists ───────────────────────────────────────────────────────────

    @Test @Order(10)
    @DisplayName("GET /playlists → lista não vazia")
    void listarPlaylists() throws Exception {
        var r = get("/playlists");
        assertEquals(200, r.statusCode());
        JsonNode body = json.readTree(r.body());
        assertTrue(body.isArray() && body.size() > 0);
        System.out.println("  [OK] " + body.size() + " playlists encontradas");
    }

    @Test @Order(11)
    @DisplayName("POST /playlists → cria playlist")
    void criarPlaylist() throws Exception {
        var r = post("/playlists", "{\"nome\":\"Playlist Teste\",\"usuario_id\":1}");
        assertEquals(201, r.statusCode());
        playlistIdCriada = json.readTree(r.body()).get("id").asInt();
        System.out.println("  [OK] Playlist criada id=" + playlistIdCriada);
    }

    @Test @Order(12)
    @DisplayName("POST /playlists/{id}/musicas → adiciona música")
    void adicionarMusicaPlaylist() throws Exception {
        var r = post("/playlists/" + playlistIdCriada + "/musicas", "{\"musica_id\":1}");
        assertEquals(201, r.statusCode());
        System.out.println("  [OK] Música 1 adicionada à playlist " + playlistIdCriada);
    }

    @Test @Order(13)
    @DisplayName("GET /playlists/{id}/musicas → lista músicas da playlist")
    void listarMusicasPlaylist() throws Exception {
        var r = get("/playlists/" + playlistIdCriada + "/musicas");
        assertEquals(200, r.statusCode());
        JsonNode body = json.readTree(r.body());
        assertTrue(body.isArray() && body.size() > 0);
        System.out.println("  [OK] " + body.size() + " música(s) na playlist");
    }

    @Test @Order(14)
    @DisplayName("DELETE /playlists/{id} → apaga playlist")
    void deletarPlaylist() throws Exception {
        assertEquals(200, delete("/playlists/" + playlistIdCriada).statusCode());
        System.out.println("  [OK] Playlist apagada id=" + playlistIdCriada);
    }
}
