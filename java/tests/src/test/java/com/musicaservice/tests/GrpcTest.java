package com.musicaservice.tests;

import com.musicaservice.grpc.proto.*;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.junit.jupiter.api.*;

import static org.junit.jupiter.api.Assertions.*;

@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
@DisplayName("gRPC — HOST:8093")
public class GrpcTest {

    private static final String HOST = System.getProperty("test.host", "localhost");
    private static final int PORT = 8093;

    private static ManagedChannel channel;
    private static MusicServiceGrpc.MusicServiceBlockingStub stub;

    private static int musicaIdCriada;
    private static int usuarioIdCriado;
    private static int playlistIdCriada;

    @BeforeAll
    static void conectar() {
        channel = ManagedChannelBuilder.forAddress(HOST, PORT)
            .usePlaintext()
            .build();
        stub = MusicServiceGrpc.newBlockingStub(channel);
        System.out.println("  Conectado a " + HOST + ":" + PORT);
    }

    @AfterAll
    static void desconectar() throws InterruptedException {
        channel.shutdownNow();
    }

    // ─── Músicas ─────────────────────────────────────────────────────────────

    @Test @Order(1)
    @DisplayName("ListMusics → lista não vazia")
    void listarMusicas() {
        MusicListResponse r = stub.listMusics(Empty.newBuilder().build());
        assertTrue(r.getMusicsList().size() > 0);
        System.out.println("  [OK] " + r.getMusicsList().size() + " músicas");
    }

    @Test @Order(2)
    @DisplayName("GetMusic id=1 → retorna música")
    void buscarMusica() {
        Music m = stub.getMusic(IdRequest.newBuilder().setId(1).build());
        assertEquals(1, m.getId());
        System.out.println("  [OK] " + m.getNome() + " — " + m.getArtista());
    }

    @Test @Order(3)
    @DisplayName("CreateMusic → retorna id gerado")
    void criarMusica() {
        Music m = stub.createMusic(CreateMusicRequest.newBuilder()
            .setNome("Teste gRPC").setArtista("gRPC Artist").build());
        musicaIdCriada = m.getId();
        assertTrue(musicaIdCriada > 0);
        System.out.println("  [OK] Criada id=" + musicaIdCriada);
    }

    @Test @Order(4)
    @DisplayName("UpdateMusic → nome actualizado")
    void atualizarMusica() {
        Music m = stub.updateMusic(UpdateMusicRequest.newBuilder()
            .setId(musicaIdCriada)
            .setNome("gRPC Actualizado")
            .setArtista("Novo Artista gRPC")
            .build());
        assertEquals("gRPC Actualizado", m.getNome());
        System.out.println("  [OK] Actualizada");
    }

    @Test @Order(5)
    @DisplayName("DeleteMusic → mensagem de sucesso")
    void deletarMusica() {
        MessageResponse r = stub.deleteMusic(IdRequest.newBuilder().setId(musicaIdCriada).build());
        assertTrue(r.getMessage().contains("deletada"));
        System.out.println("  [OK] Apagada id=" + musicaIdCriada);
    }

    // ─── Utilizadores ────────────────────────────────────────────────────────

    @Test @Order(6)
    @DisplayName("ListUsers → lista não vazia")
    void listarUsuarios() {
        UserListResponse r = stub.listUsers(Empty.newBuilder().build());
        assertTrue(r.getUsersList().size() > 0);
        System.out.println("  [OK] " + r.getUsersList().size() + " utilizadores");
    }

    @Test @Order(7)
    @DisplayName("CreateUser → retorna id gerado")
    void criarUsuario() {
        User u = stub.createUser(CreateUserRequest.newBuilder()
            .setNome("User gRPC").setIdade(28).build());
        usuarioIdCriado = u.getId();
        assertTrue(usuarioIdCriado > 0);
        System.out.println("  [OK] Utilizador criado id=" + usuarioIdCriado);
    }

    @Test @Order(8)
    @DisplayName("DeleteUser → mensagem de sucesso")
    void deletarUsuario() {
        MessageResponse r = stub.deleteUser(IdRequest.newBuilder().setId(usuarioIdCriado).build());
        assertTrue(r.getMessage().contains("deletado"));
        System.out.println("  [OK] Utilizador apagado id=" + usuarioIdCriado);
    }

    // ─── Playlists ───────────────────────────────────────────────────────────

    @Test @Order(9)
    @DisplayName("ListPlaylists → lista não vazia")
    void listarPlaylists() {
        PlaylistListResponse r = stub.listPlaylists(Empty.newBuilder().build());
        assertTrue(r.getPlaylistsList().size() > 0);
        System.out.println("  [OK] " + r.getPlaylistsList().size() + " playlists");
    }

    @Test @Order(10)
    @DisplayName("CreatePlaylist → retorna id gerado")
    void criarPlaylist() {
        Playlist p = stub.createPlaylist(CreatePlaylistRequest.newBuilder()
            .setNome("Playlist gRPC").setUsuarioId(1).build());
        playlistIdCriada = p.getId();
        assertTrue(playlistIdCriada > 0);
        System.out.println("  [OK] Playlist criada id=" + playlistIdCriada);
    }

    @Test @Order(11)
    @DisplayName("AddMusicToPlaylist → mensagem de sucesso")
    void adicionarMusicaPlaylist() {
        MessageResponse r = stub.addMusicToPlaylist(AddMusicToPlaylistRequest.newBuilder()
            .setPlaylistId(playlistIdCriada).setMusicaId(1).build());
        assertTrue(r.getMessage().contains("adicionada"));
        System.out.println("  [OK] Música adicionada");
    }

    @Test @Order(12)
    @DisplayName("ListPlaylistMusics → músicas da playlist")
    void listarMusicasPlaylist() {
        MusicListResponse r = stub.listPlaylistMusics(
            PlaylistMusicsRequest.newBuilder().setPlaylistId(playlistIdCriada).build());
        assertTrue(r.getMusicsList().size() > 0);
        System.out.println("  [OK] " + r.getMusicsList().size() + " música(s) na playlist");
    }

    @Test @Order(13)
    @DisplayName("DeletePlaylist → mensagem de sucesso")
    void deletarPlaylist() {
        MessageResponse r = stub.deletePlaylist(IdRequest.newBuilder().setId(playlistIdCriada).build());
        assertTrue(r.getMessage().contains("deletada"));
        System.out.println("  [OK] Playlist apagada id=" + playlistIdCriada);
    }
}
