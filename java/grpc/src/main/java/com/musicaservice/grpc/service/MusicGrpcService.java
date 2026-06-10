package com.musicaservice.grpc.service;

import com.musicaservice.grpc.proto.*;
import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;

import java.util.List;

@GrpcService
public class MusicGrpcService extends MusicServiceGrpc.MusicServiceImplBase {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    private final RowMapper<Music> musicRM = (rs, n) -> Music.newBuilder()
        .setId(rs.getInt("id"))
        .setNome(rs.getString("nome"))
        .setArtista(rs.getString("artista"))
        .build();

    private final RowMapper<User> userRM = (rs, n) -> User.newBuilder()
        .setId(rs.getInt("id"))
        .setNome(rs.getString("nome"))
        .setIdade(rs.getInt("idade"))
        .build();

    private final RowMapper<Playlist> playlistRM = (rs, n) -> Playlist.newBuilder()
        .setId(rs.getInt("id"))
        .setNome(rs.getString("nome"))
        .setUsuarioId(rs.getInt("usuario_id"))
        .build();

    // ── MUSIC ────────────────────────────────────────────────────────────────

    @Override
    public void listMusics(Empty request, StreamObserver<MusicListResponse> responseObserver) {
        List<Music> musics = jdbcTemplate.query(
            "SELECT id, nome, artista FROM musica ORDER BY id", musicRM);
        responseObserver.onNext(MusicListResponse.newBuilder().addAllMusics(musics).build());
        responseObserver.onCompleted();
    }

    @Override
    public void getMusic(IdRequest request, StreamObserver<Music> responseObserver) {
        List<Music> r = jdbcTemplate.query(
            "SELECT id, nome, artista FROM musica WHERE id = ?", musicRM, request.getId());
        if (!r.isEmpty()) responseObserver.onNext(r.get(0));
        responseObserver.onCompleted();
    }

    @Override
    public void createMusic(CreateMusicRequest request, StreamObserver<Music> responseObserver) {
        int novoId = jdbcTemplate.queryForObject(
            "INSERT INTO musica (nome, artista) VALUES (?, ?) RETURNING id",
            Integer.class, request.getNome(), request.getArtista());
        responseObserver.onNext(Music.newBuilder()
            .setId(novoId).setNome(request.getNome()).setArtista(request.getArtista()).build());
        responseObserver.onCompleted();
    }

    @Override
    public void updateMusic(UpdateMusicRequest request, StreamObserver<Music> responseObserver) {
        jdbcTemplate.update("UPDATE musica SET nome = ?, artista = ? WHERE id = ?",
            request.getNome(), request.getArtista(), request.getId());
        responseObserver.onNext(Music.newBuilder()
            .setId(request.getId()).setNome(request.getNome()).setArtista(request.getArtista()).build());
        responseObserver.onCompleted();
    }

    @Override
    public void deleteMusic(IdRequest request, StreamObserver<MessageResponse> responseObserver) {
        jdbcTemplate.update("DELETE FROM musica WHERE id = ?", request.getId());
        responseObserver.onNext(MessageResponse.newBuilder().setMessage("Música deletada com sucesso").build());
        responseObserver.onCompleted();
    }

    // ── USER ─────────────────────────────────────────────────────────────────

    @Override
    public void listUsers(Empty request, StreamObserver<UserListResponse> responseObserver) {
        List<User> users = jdbcTemplate.query(
            "SELECT id, nome, idade FROM usuario ORDER BY id", userRM);
        responseObserver.onNext(UserListResponse.newBuilder().addAllUsers(users).build());
        responseObserver.onCompleted();
    }

    @Override
    public void getUser(IdRequest request, StreamObserver<User> responseObserver) {
        List<User> r = jdbcTemplate.query(
            "SELECT id, nome, idade FROM usuario WHERE id = ?", userRM, request.getId());
        if (!r.isEmpty()) responseObserver.onNext(r.get(0));
        responseObserver.onCompleted();
    }

    @Override
    public void createUser(CreateUserRequest request, StreamObserver<User> responseObserver) {
        int novoId = jdbcTemplate.queryForObject(
            "INSERT INTO usuario (nome, idade) VALUES (?, ?) RETURNING id",
            Integer.class, request.getNome(), request.getIdade());
        responseObserver.onNext(User.newBuilder()
            .setId(novoId).setNome(request.getNome()).setIdade(request.getIdade()).build());
        responseObserver.onCompleted();
    }

    @Override
    public void updateUser(UpdateUserRequest request, StreamObserver<User> responseObserver) {
        jdbcTemplate.update("UPDATE usuario SET nome = ?, idade = ? WHERE id = ?",
            request.getNome(), request.getIdade(), request.getId());
        responseObserver.onNext(User.newBuilder()
            .setId(request.getId()).setNome(request.getNome()).setIdade(request.getIdade()).build());
        responseObserver.onCompleted();
    }

    @Override
    public void deleteUser(IdRequest request, StreamObserver<MessageResponse> responseObserver) {
        jdbcTemplate.update("DELETE FROM usuario WHERE id = ?", request.getId());
        responseObserver.onNext(MessageResponse.newBuilder().setMessage("Usuário deletado com sucesso").build());
        responseObserver.onCompleted();
    }

    // ── PLAYLIST ─────────────────────────────────────────────────────────────

    @Override
    public void listPlaylists(Empty request, StreamObserver<PlaylistListResponse> responseObserver) {
        List<Playlist> playlists = jdbcTemplate.query(
            "SELECT id, nome, usuario_id FROM playlist ORDER BY id", playlistRM);
        responseObserver.onNext(PlaylistListResponse.newBuilder().addAllPlaylists(playlists).build());
        responseObserver.onCompleted();
    }

    @Override
    public void getPlaylist(IdRequest request, StreamObserver<Playlist> responseObserver) {
        List<Playlist> r = jdbcTemplate.query(
            "SELECT id, nome, usuario_id FROM playlist WHERE id = ?", playlistRM, request.getId());
        if (!r.isEmpty()) responseObserver.onNext(r.get(0));
        responseObserver.onCompleted();
    }

    @Override
    public void createPlaylist(CreatePlaylistRequest request, StreamObserver<Playlist> responseObserver) {
        int novoId = jdbcTemplate.queryForObject(
            "INSERT INTO playlist (nome, usuario_id) VALUES (?, ?) RETURNING id",
            Integer.class, request.getNome(), request.getUsuarioId());
        responseObserver.onNext(Playlist.newBuilder()
            .setId(novoId).setNome(request.getNome()).setUsuarioId(request.getUsuarioId()).build());
        responseObserver.onCompleted();
    }

    @Override
    public void updatePlaylist(UpdatePlaylistRequest request, StreamObserver<Playlist> responseObserver) {
        jdbcTemplate.update("UPDATE playlist SET nome = ?, usuario_id = ? WHERE id = ?",
            request.getNome(), request.getUsuarioId(), request.getId());
        responseObserver.onNext(Playlist.newBuilder()
            .setId(request.getId()).setNome(request.getNome()).setUsuarioId(request.getUsuarioId()).build());
        responseObserver.onCompleted();
    }

    @Override
    public void deletePlaylist(IdRequest request, StreamObserver<MessageResponse> responseObserver) {
        jdbcTemplate.update("DELETE FROM playlist WHERE id = ?", request.getId());
        responseObserver.onNext(MessageResponse.newBuilder().setMessage("Playlist deletada com sucesso").build());
        responseObserver.onCompleted();
    }

    // ── PLAYLIST-MUSIC ────────────────────────────────────────────────────────

    @Override
    public void addMusicToPlaylist(AddMusicToPlaylistRequest request, StreamObserver<MessageResponse> responseObserver) {
        jdbcTemplate.update(
            "INSERT INTO playlist_musica (playlist_id, musica_id) VALUES (?, ?)",
            request.getPlaylistId(), request.getMusicaId());
        responseObserver.onNext(MessageResponse.newBuilder().setMessage("Música adicionada na playlist com sucesso").build());
        responseObserver.onCompleted();
    }

    @Override
    public void listPlaylistMusics(PlaylistMusicsRequest request, StreamObserver<MusicListResponse> responseObserver) {
        List<Music> musics = jdbcTemplate.query(
            "SELECT m.id, m.nome, m.artista FROM musica m " +
            "INNER JOIN playlist_musica pm ON pm.musica_id = m.id " +
            "WHERE pm.playlist_id = ? ORDER BY m.id",
            musicRM, request.getPlaylistId());
        responseObserver.onNext(MusicListResponse.newBuilder().addAllMusics(musics).build());
        responseObserver.onCompleted();
    }
}
