package com.musicaservice.rest.controller;

import com.musicaservice.rest.model.Musica;
import com.musicaservice.rest.model.Playlist;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/playlists")
public class PlaylistController {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    private final RowMapper<Playlist> playlistRowMapper = (rs, rowNum) ->
        new Playlist(rs.getInt("id"), rs.getString("nome"), rs.getInt("usuario_id"));

    private final RowMapper<Musica> musicaRowMapper = (rs, rowNum) ->
        new Musica(rs.getInt("id"), rs.getString("nome"), rs.getString("artista"));

    @GetMapping
    public List<Playlist> listarPlaylists() {
        return jdbcTemplate.query("SELECT id, nome, usuario_id FROM playlist ORDER BY id", playlistRowMapper);
    }

    @GetMapping("/{id}")
    public ResponseEntity<Playlist> buscarPlaylist(@PathVariable int id) {
        List<Playlist> result = jdbcTemplate.query(
            "SELECT id, nome, usuario_id FROM playlist WHERE id = ?", playlistRowMapper, id);
        if (result.isEmpty()) return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
        return ResponseEntity.ok(result.get(0));
    }

    @PostMapping
    public ResponseEntity<Playlist> criarPlaylist(@RequestBody Playlist playlist) {
        int novoId = jdbcTemplate.queryForObject(
            "INSERT INTO playlist (nome, usuario_id) VALUES (?, ?) RETURNING id",
            Integer.class, playlist.getNome(), playlist.getUsuario_id());
        playlist.setId(novoId);
        return ResponseEntity.status(HttpStatus.CREATED).body(playlist);
    }

    @PutMapping("/{id}")
    public ResponseEntity<Playlist> atualizarPlaylist(@PathVariable int id, @RequestBody Playlist playlist) {
        int rows = jdbcTemplate.update(
            "UPDATE playlist SET nome = ?, usuario_id = ? WHERE id = ?",
            playlist.getNome(), playlist.getUsuario_id(), id);
        if (rows == 0) return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
        playlist.setId(id);
        return ResponseEntity.ok(playlist);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Playlist> deletarPlaylist(@PathVariable int id) {
        List<Playlist> result = jdbcTemplate.query(
            "SELECT id, nome, usuario_id FROM playlist WHERE id = ?", playlistRowMapper, id);
        if (result.isEmpty()) return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
        jdbcTemplate.update("DELETE FROM playlist WHERE id = ?", id);
        return ResponseEntity.ok(result.get(0));
    }

    @GetMapping("/{id}/musicas")
    public ResponseEntity<List<Musica>> listarMusicasPlaylist(@PathVariable int id) {
        List<Playlist> pl = jdbcTemplate.query(
            "SELECT id, nome, usuario_id FROM playlist WHERE id = ?", playlistRowMapper, id);
        if (pl.isEmpty()) return ResponseEntity.status(HttpStatus.NOT_FOUND).build();

        List<Musica> musicas = jdbcTemplate.query(
            "SELECT m.id, m.nome, m.artista FROM musica m " +
            "INNER JOIN playlist_musica pm ON pm.musica_id = m.id " +
            "WHERE pm.playlist_id = ? ORDER BY m.id",
            musicaRowMapper, id);
        return ResponseEntity.ok(musicas);
    }

    @PostMapping("/{id}/musicas")
    public ResponseEntity<Void> adicionarMusicaPlaylist(
            @PathVariable int id,
            @RequestBody Map<String, Integer> body) {
        Integer musicaId = body.get("musica_id");
        if (musicaId == null) return ResponseEntity.status(HttpStatus.BAD_REQUEST).build();
        jdbcTemplate.update(
            "INSERT INTO playlist_musica (playlist_id, musica_id) VALUES (?, ?)", id, musicaId);
        return ResponseEntity.status(HttpStatus.CREATED).build();
    }
}
