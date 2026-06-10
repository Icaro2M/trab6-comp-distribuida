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
@DisplayName("GraphQL — http://HOST:8092/graphql")
public class GraphqlTest {

    private static final String HOST = System.getProperty("test.host", "localhost");
    private static final String ENDPOINT = "http://" + HOST + ":8092/graphql";
    private static final HttpClient http = HttpClient.newHttpClient();
    private static final ObjectMapper json = new ObjectMapper();

    private static int musicaIdCriada;
    private static int usuarioIdCriado;
    private static int playlistIdCriada;

    // ─── Utilitário ─────────────────────────────────────────────────────────

    private JsonNode gql(String query) throws Exception {
        String body = json.writeValueAsString(java.util.Map.of("query", query));
        HttpResponse<String> r = http.send(
            HttpRequest.newBuilder(URI.create(ENDPOINT))
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(body)).build(),
            HttpResponse.BodyHandlers.ofString());

        assertEquals(200, r.statusCode(), "GraphQL request falhou");
        JsonNode root = json.readTree(r.body());
        assertNull(root.get("errors"), "GraphQL retornou erros: " + root.get("errors"));
        return root.get("data");
    }

    // ─── Músicas ─────────────────────────────────────────────────────────────

    @Test @Order(1)
    @DisplayName("Query musicas → lista não vazia")
    void listarMusicas() throws Exception {
        JsonNode data = gql("{ musicas { id nome artista } }");
        JsonNode musicas = data.get("musicas");
        assertTrue(musicas.isArray() && musicas.size() > 0);
        System.out.println("  [OK] " + musicas.size() + " músicas");
    }

    @Test @Order(2)
    @DisplayName("Query musica(id:1) → retorna música")
    void buscarMusica() throws Exception {
        JsonNode data = gql("{ musica(id: 1) { id nome artista } }");
        assertEquals(1, data.get("musica").get("id").asInt());
        System.out.println("  [OK] " + data.get("musica").get("nome").asText());
    }

    @Test @Order(3)
    @DisplayName("Mutation criarMusica → retorna id")
    void criarMusica() throws Exception {
        JsonNode data = gql(
            "mutation { criarMusica(nome: \"Teste GraphQL\", artista: \"GQL Artist\") { id nome } }");
        musicaIdCriada = data.get("criarMusica").get("id").asInt();
        assertTrue(musicaIdCriada > 0);
        System.out.println("  [OK] Criada id=" + musicaIdCriada);
    }

    @Test @Order(4)
    @DisplayName("Mutation atualizarMusica → nome actualizado")
    void atualizarMusica() throws Exception {
        JsonNode data = gql(
            "mutation { atualizarMusica(id: " + musicaIdCriada +
            ", nome: \"GQL Actualizado\", artista: \"Novo\") { nome } }");
        assertEquals("GQL Actualizado", data.get("atualizarMusica").get("nome").asText());
        System.out.println("  [OK] Actualizada");
    }

    @Test @Order(5)
    @DisplayName("Mutation deletarMusica → mensagem de sucesso")
    void deletarMusica() throws Exception {
        JsonNode data = gql("mutation { deletarMusica(id: " + musicaIdCriada + ") }");
        assertTrue(data.get("deletarMusica").asText().contains("deletada"));
        System.out.println("  [OK] Apagada id=" + musicaIdCriada);
    }

    // ─── Utilizadores ────────────────────────────────────────────────────────

    @Test @Order(6)
    @DisplayName("Query usuarios → lista não vazia")
    void listarUsuarios() throws Exception {
        JsonNode data = gql("{ usuarios { id nome idade } }");
        assertTrue(data.get("usuarios").size() > 0);
        System.out.println("  [OK] " + data.get("usuarios").size() + " utilizadores");
    }

    @Test @Order(7)
    @DisplayName("Mutation criarUsuario → retorna id")
    void criarUsuario() throws Exception {
        JsonNode data = gql(
            "mutation { criarUsuario(nome: \"User GQL\", idade: 30) { id } }");
        usuarioIdCriado = data.get("criarUsuario").get("id").asInt();
        System.out.println("  [OK] Utilizador criado id=" + usuarioIdCriado);
    }

    @Test @Order(8)
    @DisplayName("Mutation deletarUsuario → mensagem de sucesso")
    void deletarUsuario() throws Exception {
        JsonNode data = gql("mutation { deletarUsuario(id: " + usuarioIdCriado + ") }");
        assertTrue(data.get("deletarUsuario").asText().contains("deletado"));
        System.out.println("  [OK] Utilizador apagado id=" + usuarioIdCriado);
    }

    // ─── Playlists ───────────────────────────────────────────────────────────

    @Test @Order(9)
    @DisplayName("Query playlists → lista não vazia")
    void listarPlaylists() throws Exception {
        JsonNode data = gql("{ playlists { id nome usuario_id } }");
        assertTrue(data.get("playlists").size() > 0);
        System.out.println("  [OK] " + data.get("playlists").size() + " playlists");
    }

    @Test @Order(10)
    @DisplayName("Mutation criarPlaylist → retorna id")
    void criarPlaylist() throws Exception {
        JsonNode data = gql(
            "mutation { criarPlaylist(nome: \"Playlist GQL\", usuario_id: 1) { id } }");
        playlistIdCriada = data.get("criarPlaylist").get("id").asInt();
        System.out.println("  [OK] Playlist criada id=" + playlistIdCriada);
    }

    @Test @Order(11)
    @DisplayName("Mutation adicionarMusicaNaPlaylist → mensagem de sucesso")
    void adicionarMusicaPlaylist() throws Exception {
        JsonNode data = gql(
            "mutation { adicionarMusicaNaPlaylist(playlist_id: " + playlistIdCriada +
            ", musica_id: 1) }");
        assertTrue(data.get("adicionarMusicaNaPlaylist").asText().contains("adicionada"));
        System.out.println("  [OK] Música adicionada");
    }

    @Test @Order(12)
    @DisplayName("Query musicasDaPlaylist → músicas da playlist")
    void musicasDaPlaylist() throws Exception {
        JsonNode data = gql(
            "{ musicasDaPlaylist(playlist_id: " + playlistIdCriada + ") { id nome } }");
        assertTrue(data.get("musicasDaPlaylist").size() > 0);
        System.out.println("  [OK] " + data.get("musicasDaPlaylist").size() + " música(s) na playlist");
    }

    @Test @Order(13)
    @DisplayName("Mutation deletarPlaylist → mensagem de sucesso")
    void deletarPlaylist() throws Exception {
        JsonNode data = gql("mutation { deletarPlaylist(id: " + playlistIdCriada + ") }");
        assertTrue(data.get("deletarPlaylist").asText().contains("deletada"));
        System.out.println("  [OK] Playlist apagada id=" + playlistIdCriada);
    }
}
