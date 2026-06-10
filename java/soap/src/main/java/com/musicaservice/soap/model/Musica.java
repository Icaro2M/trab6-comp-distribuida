package com.musicaservice.soap.model;

import jakarta.xml.bind.annotation.XmlAccessType;
import jakarta.xml.bind.annotation.XmlAccessorType;
import jakarta.xml.bind.annotation.XmlType;

@XmlAccessorType(XmlAccessType.FIELD)
@XmlType(name = "Musica")
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
