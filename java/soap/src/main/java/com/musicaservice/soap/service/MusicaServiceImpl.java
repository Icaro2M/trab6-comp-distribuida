package com.musicaservice.soap.service;

import com.musicaservice.soap.model.Musica;
import com.musicaservice.soap.model.Playlist;
import com.musicaservice.soap.model.Usuario;
import jakarta.jws.WebService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@WebService(endpointInterface = "com.musicaservice.soap.service.MusicaService",
            serviceName = "MusicaService",
            targetNamespace = "http://service.soap.musicaservice.com/")
public class MusicaServiceImpl implements MusicaService {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    private final RowMapper<Musica> musicaRM = (rs, n) ->
        new Musica(rs.getInt("id"), rs.getString("nome"), rs.getString("artista"));

    private final RowMapper<Usuario> usuarioRM = (rs, n) ->
        new Usuario(rs.getInt("id"), rs.getString("nome"), rs.getInt("idade"));

    private final RowMapper<Playlist> playlistRM = (rs, n) ->
        new Playlist(rs.getInt("id"), rs.getString("nome"), rs.getInt("usuario_id"));

    // ── MUSICA ──────────────────────────────────────────────────────────────

    @Override
    public List<Musica> listarMusicas() {
        return jdbcTemplate.query("SELECT id, nome, artista FROM musica ORDER BY id", musicaRM);
    }

    @Override
    public Musica buscarMusica(int id) {
        List<Musica> r = jdbcTemplate.query(
            "SELECT id, nome, artista FROM musica WHERE id = ?", musicaRM, id);
        return r.isEmpty() ? null : r.get(0);
    }

    @Override
    public Musica criarMusica(String nome, String artista) {
        int novoId = jdbcTemplate.queryForObject(
            "INSERT INTO musica (nome, artista) VALUES (?, ?) RETURNING id",
            Integer.class, nome, artista);
        return new Musica(novoId, nome, artista);
    }

    @Override
    public Musica atualizarMusica(int id, String nome, String artista) {
        jdbcTemplate.update("UPDATE musica SET nome = ?, artista = ? WHERE id = ?", nome, artista, id);
        return new Musica(id, nome, artista);
    }

    @Override
    public String deletarMusica(int id) {
        jdbcTemplate.update("DELETE FROM musica WHERE id = ?", id);
        return "Música deletada com sucesso";
    }

    // ── USUARIO ─────────────────────────────────────────────────────────────

    @Override
    public List<Usuario> listarUsuarios() {
        return jdbcTemplate.query("SELECT id, nome, idade FROM usuario ORDER BY id", usuarioRM);
    }

    @Override
    public Usuario buscarUsuario(int id) {
        List<Usuario> r = jdbcTemplate.query(
            "SELECT id, nome, idade FROM usuario WHERE id = ?", usuarioRM, id);
        return r.isEmpty() ? null : r.get(0);
    }

    @Override
    public Usuario criarUsuario(String nome, int idade) {
        int novoId = jdbcTemplate.queryForObject(
            "INSERT INTO usuario (nome, idade) VALUES (?, ?) RETURNING id",
            Integer.class, nome, idade);
        return new Usuario(novoId, nome, idade);
    }

    @Override
    public Usuario atualizarUsuario(int id, String nome, int idade) {
        jdbcTemplate.update("UPDATE usuario SET nome = ?, idade = ? WHERE id = ?", nome, idade, id);
        return new Usuario(id, nome, idade);
    }

    @Override
    public String deletarUsuario(int id) {
        jdbcTemplate.update("DELETE FROM usuario WHERE id = ?", id);
        return "Usuário deletado com sucesso";
    }

    // ── PLAYLIST ─────────────────────────────────────────────────────────────

    @Override
    public List<Playlist> listarPlaylists() {
        return jdbcTemplate.query("SELECT id, nome, usuario_id FROM playlist ORDER BY id", playlistRM);
    }

    @Override
    public Playlist buscarPlaylist(int id) {
        List<Playlist> r = jdbcTemplate.query(
            "SELECT id, nome, usuario_id FROM playlist WHERE id = ?", playlistRM, id);
        return r.isEmpty() ? null : r.get(0);
    }

    @Override
    public Playlist criarPlaylist(String nome, int usuario_id) {
        int novoId = jdbcTemplate.queryForObject(
            "INSERT INTO playlist (nome, usuario_id) VALUES (?, ?) RETURNING id",
            Integer.class, nome, usuario_id);
        return new Playlist(novoId, nome, usuario_id);
    }

    @Override
    public Playlist atualizarPlaylist(int id, String nome, int usuario_id) {
        jdbcTemplate.update("UPDATE playlist SET nome = ?, usuario_id = ? WHERE id = ?", nome, usuario_id, id);
        return new Playlist(id, nome, usuario_id);
    }

    @Override
    public String deletarPlaylist(int id) {
        jdbcTemplate.update("DELETE FROM playlist WHERE id = ?", id);
        return "Playlist deletada com sucesso";
    }

    // ── PLAYLIST-MUSICA ──────────────────────────────────────────────────────

    @Override
    public List<Musica> listarMusicasPlaylist(int playlist_id) {
        return jdbcTemplate.query(
            "SELECT m.id, m.nome, m.artista FROM musica m " +
            "INNER JOIN playlist_musica pm ON pm.musica_id = m.id " +
            "WHERE pm.playlist_id = ? ORDER BY m.id",
            musicaRM, playlist_id);
    }

    @Override
    public String adicionarMusicaNaPlaylist(int playlist_id, int musica_id) {
        jdbcTemplate.update(
            "INSERT INTO playlist_musica (playlist_id, musica_id) VALUES (?, ?)",
            playlist_id, musica_id);
        return "Música adicionada na playlist com sucesso";
    }
}
