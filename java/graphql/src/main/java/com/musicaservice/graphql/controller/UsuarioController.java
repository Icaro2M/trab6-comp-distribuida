package com.musicaservice.graphql.controller;

import com.musicaservice.graphql.model.Usuario;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.graphql.data.method.annotation.Argument;
import org.springframework.graphql.data.method.annotation.MutationMapping;
import org.springframework.graphql.data.method.annotation.QueryMapping;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.stereotype.Controller;

import java.util.List;

@Controller
public class UsuarioController {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    private final RowMapper<Usuario> rowMapper = (rs, n) ->
        new Usuario(rs.getInt("id"), rs.getString("nome"), rs.getInt("idade"));

    @QueryMapping
    public List<Usuario> usuarios() {
        return jdbcTemplate.query("SELECT id, nome, idade FROM usuario ORDER BY id", rowMapper);
    }

    @QueryMapping
    public Usuario usuario(@Argument int id) {
        List<Usuario> r = jdbcTemplate.query(
            "SELECT id, nome, idade FROM usuario WHERE id = ?", rowMapper, id);
        return r.isEmpty() ? null : r.get(0);
    }

    @MutationMapping
    public Usuario criarUsuario(@Argument String nome, @Argument int idade) {
        int novoId = jdbcTemplate.queryForObject(
            "INSERT INTO usuario (nome, idade) VALUES (?, ?) RETURNING id",
            Integer.class, nome, idade);
        return new Usuario(novoId, nome, idade);
    }

    @MutationMapping
    public Usuario atualizarUsuario(@Argument int id, @Argument String nome, @Argument int idade) {
        jdbcTemplate.update("UPDATE usuario SET nome = ?, idade = ? WHERE id = ?", nome, idade, id);
        return new Usuario(id, nome, idade);
    }

    @MutationMapping
    public String deletarUsuario(@Argument int id) {
        jdbcTemplate.update("DELETE FROM usuario WHERE id = ?", id);
        return "Usuário deletado com sucesso";
    }
}
