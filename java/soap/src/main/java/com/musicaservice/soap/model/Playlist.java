package com.musicaservice.soap.model;

import jakarta.xml.bind.annotation.XmlAccessType;
import jakarta.xml.bind.annotation.XmlAccessorType;
import jakarta.xml.bind.annotation.XmlType;

@XmlAccessorType(XmlAccessType.FIELD)
@XmlType(name = "Playlist")
public class Playlist {
    private int id;
    private String nome;
    private int usuario_id;

    public Playlist() {}

    public Playlist(int id, String nome, int usuario_id) {
        this.id = id;
        this.nome = nome;
        this.usuario_id = usuario_id;
    }

    public int getId() { return id; }
    public void setId(int id) { this.id = id; }
    public String getNome() { return nome; }
    public void setNome(String nome) { this.nome = nome; }
    public int getUsuario_id() { return usuario_id; }
    public void setUsuario_id(int usuario_id) { this.usuario_id = usuario_id; }
}
