package com.musicaservice.rest.controller;

import com.musicaservice.rest.model.Musica;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/musicas")
public class MusicaController {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    private final RowMapper<Musica> rowMapper = (rs, rowNum) ->
        new Musica(rs.getInt("id"), rs.getString("nome"), rs.getString("artista"));

    @GetMapping
    public List<Musica> listarMusicas() {
        return jdbcTemplate.query("SELECT id, nome, artista FROM musica ORDER BY id", rowMapper);
    }

    @GetMapping("/{id}")
    public ResponseEntity<Musica> buscarMusica(@PathVariable int id) {
        List<Musica> result = jdbcTemplate.query(
            "SELECT id, nome, artista FROM musica WHERE id = ?", rowMapper, id);
        if (result.isEmpty()) return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
        return ResponseEntity.ok(result.get(0));
    }

    @PostMapping
    public ResponseEntity<Musica> criarMusica(@RequestBody Musica musica) {
        int novoId = jdbcTemplate.queryForObject(
            "INSERT INTO musica (nome, artista) VALUES (?, ?) RETURNING id",
            Integer.class, musica.getNome(), musica.getArtista());
        musica.setId(novoId);
        return ResponseEntity.status(HttpStatus.CREATED).body(musica);
    }

    @PutMapping("/{id}")
    public ResponseEntity<Musica> atualizarMusica(@PathVariable int id, @RequestBody Musica musica) {
        int rows = jdbcTemplate.update(
            "UPDATE musica SET nome = ?, artista = ? WHERE id = ?",
            musica.getNome(), musica.getArtista(), id);
        if (rows == 0) return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
        musica.setId(id);
        return ResponseEntity.ok(musica);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Musica> deletarMusica(@PathVariable int id) {
        List<Musica> result = jdbcTemplate.query(
            "SELECT id, nome, artista FROM musica WHERE id = ?", rowMapper, id);
        if (result.isEmpty()) return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
        jdbcTemplate.update("DELETE FROM musica WHERE id = ?", id);
        return ResponseEntity.ok(result.get(0));
    }
}
