package com.musicaservice.rest.model;

public class Musica {
    private int id;
    private String nome;
    private String artista;

    public Musica() {}

    public Musica(int id, String nome, String artista) {
        this.id = id;
        this.nome = nome;
        this.artista = artista;
    }

    public int getId() { return id; }
    public void setId(int id) { this.id = id; }
    public String getNome() { return nome; }
    public void setNome(String nome) { this.nome = nome; }
    public String getArtista() { return artista; }
    public void setArtista(String artista) { this.artista = artista; }
}
