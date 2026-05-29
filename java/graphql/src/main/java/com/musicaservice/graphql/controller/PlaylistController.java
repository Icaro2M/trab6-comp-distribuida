package com.musicaservice.graphql.controller;

import com.musicaservice.graphql.model.Musica;
import com.musicaservice.graphql.model.Playlist;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.graphql.data.method.annotation.Argument;
import org.springframework.graphql.data.method.annotation.MutationMapping;
import org.springframework.graphql.data.method.annotation.QueryMapping;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.stereotype.Controller;

import java.util.List;

@Controller
public class PlaylistController {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    private final RowMapper<Playlist> playlistRM = (rs, n) ->
        new Playlist(rs.getInt("id"), rs.getString("nome"), rs.getInt("usuario_id"));

    private final RowMapper<Musica> musicaRM = (rs, n) ->
        new Musica(rs.getInt("id"), rs.getString("nome"), rs.getString("artista"));

    @QueryMapping
    public List<Playlist> playlists() {
        return jdbcTemplate.query("SELECT id, nome, usuario_id FROM playlist ORDER BY id", playlistRM);
    }

    @QueryMapping
    public Playlist playlist(@Argument int id) {
        List<Playlist> r = jdbcTemplate.query(
            "SELECT id, nome, usuario_id FROM playlist WHERE id = ?", playlistRM, id);
        return r.isEmpty() ? null : r.get(0);
    }

    @QueryMapping
    public List<Musica> musicasDaPlaylist(@Argument int playlist_id) {
        return jdbcTemplate.query(
            "SELECT m.id, m.nome, m.artista FROM musica m " +
            "INNER JOIN playlist_musica pm ON pm.musica_id = m.id " +
            "WHERE pm.playlist_id = ? ORDER BY m.id",
            musicaRM, playlist_id);
    }

    @MutationMapping
    public Playlist criarPlaylist(@Argument String nome, @Argument int usuario_id) {
        int novoId = jdbcTemplate.queryForObject(
            "INSERT INTO playlist (nome, usuario_id) VALUES (?, ?) RETURNING id",
            Integer.class, nome, usuario_id);
        return new Playlist(novoId, nome, usuario_id);
    }

    @MutationMapping
    public Playlist atualizarPlaylist(@Argument int id, @Argument String nome, @Argument int usuario_id) {
        jdbcTemplate.update("UPDATE playlist SET nome = ?, usuario_id = ? WHERE id = ?", nome, usuario_id, id);
        return new Playlist(id, nome, usuario_id);
    }

    @MutationMapping
    public String deletarPlaylist(@Argument int id) {
        jdbcTemplate.update("DELETE FROM playlist WHERE id = ?", id);
        return "Playlist deletada com sucesso";
    }

    @MutationMapping
    public String adicionarMusicaNaPlaylist(@Argument int playlist_id, @Argument int musica_id) {
        jdbcTemplate.update(
            "INSERT INTO playlist_musica (playlist_id, musica_id) VALUES (?, ?)",
            playlist_id, musica_id);
        return "Música adicionada na playlist com sucesso";
    }
}
