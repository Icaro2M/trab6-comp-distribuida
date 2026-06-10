package com.musicaservice.rest.controller;

import com.musicaservice.rest.model.Usuario;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/usuarios")
public class UsuarioController {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    private final RowMapper<Usuario> rowMapper = (rs, rowNum) ->
        new Usuario(rs.getInt("id"), rs.getString("nome"), rs.getInt("idade"));

    @GetMapping
    public List<Usuario> listarUsuarios() {
        return jdbcTemplate.query("SELECT id, nome, idade FROM usuario ORDER BY id", rowMapper);
    }

    @GetMapping("/{id}")
    public ResponseEntity<Usuario> buscarUsuario(@PathVariable int id) {
        List<Usuario> result = jdbcTemplate.query(
            "SELECT id, nome, idade FROM usuario WHERE id = ?", rowMapper, id);
        if (result.isEmpty()) return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
        return ResponseEntity.ok(result.get(0));
    }

    @PostMapping
    public ResponseEntity<Usuario> criarUsuario(@RequestBody Usuario usuario) {
        int novoId = jdbcTemplate.queryForObject(
            "INSERT INTO usuario (nome, idade) VALUES (?, ?) RETURNING id",
            Integer.class, usuario.getNome(), usuario.getIdade());
        usuario.setId(novoId);
        return ResponseEntity.status(HttpStatus.CREATED).body(usuario);
    }

    @PutMapping("/{id}")
    public ResponseEntity<Usuario> atualizarUsuario(@PathVariable int id, @RequestBody Usuario usuario) {
        int rows = jdbcTemplate.update(
            "UPDATE usuario SET nome = ?, idade = ? WHERE id = ?",
            usuario.getNome(), usuario.getIdade(), id);
        if (rows == 0) return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
        usuario.setId(id);
        return ResponseEntity.ok(usuario);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Usuario> deletarUsuario(@PathVariable int id) {
        List<Usuario> result = jdbcTemplate.query(
            "SELECT id, nome, idade FROM usuario WHERE id = ?", rowMapper, id);
        if (result.isEmpty()) return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
        jdbcTemplate.update("DELETE FROM usuario WHERE id = ?", id);
        return ResponseEntity.ok(result.get(0));
    }
}
