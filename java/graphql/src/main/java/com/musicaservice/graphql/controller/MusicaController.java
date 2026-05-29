package com.musicaservice.graphql.controller;

import com.musicaservice.graphql.model.Musica;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.graphql.data.method.annotation.Argument;
import org.springframework.graphql.data.method.annotation.MutationMapping;
import org.springframework.graphql.data.method.annotation.QueryMapping;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.stereotype.Controller;

import java.util.List;

@Controller
public class MusicaController {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    private final RowMapper<Musica> rowMapper = (rs, n) ->
        new Musica(rs.getInt("id"), rs.getString("nome"), rs.getString("artista"));

    @QueryMapping
    public List<Musica> musicas() {
        return jdbcTemplate.query("SELECT id, nome, artista FROM musica ORDER BY id", rowMapper);
    }

    @QueryMapping
    public Musica musica(@Argument int id) {
        List<Musica> r = jdbcTemplate.query(
            "SELECT id, nome, artista FROM musica WHERE id = ?", rowMapper, id);
        return r.isEmpty() ? null : r.get(0);
    }

    @MutationMapping
    public Musica criarMusica(@Argument String nome, @Argument String artista) {
        int novoId = jdbcTemplate.queryForObject(
            "INSERT INTO musica (nome, artista) VALUES (?, ?) RETURNING id",
            Integer.class, nome, artista);
        return new Musica(novoId, nome, artista);
    }

    @MutationMapping
    public Musica atualizarMusica(@Argument int id, @Argument String nome, @Argument String artista) {
        jdbcTemplate.update("UPDATE musica SET nome = ?, artista = ? WHERE id = ?", nome, artista, id);
        return new Musica(id, nome, artista);
    }

    @MutationMapping
    public String deletarMusica(@Argument int id) {
        jdbcTemplate.update("DELETE FROM musica WHERE id = ?", id);
        return "Música deletada com sucesso";
    }
}
