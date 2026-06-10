package com.musicaservice.soap.service;

import com.musicaservice.soap.model.Musica;
import com.musicaservice.soap.model.Playlist;
import com.musicaservice.soap.model.Usuario;
import jakarta.jws.WebMethod;
import jakarta.jws.WebParam;
import jakarta.jws.WebService;

import java.util.List;

@WebService(name = "MusicaService", targetNamespace = "http://service.soap.musicaservice.com/")
public interface MusicaService {

    @WebMethod List<Musica> listarMusicas();
    @WebMethod Musica buscarMusica(@WebParam(name = "id") int id);
    @WebMethod Musica criarMusica(@WebParam(name = "nome") String nome, @WebParam(name = "artista") String artista);
    @WebMethod Musica atualizarMusica(@WebParam(name = "id") int id, @WebParam(name = "nome") String nome, @WebParam(name = "artista") String artista);
    @WebMethod String deletarMusica(@WebParam(name = "id") int id);

    @WebMethod List<Usuario> listarUsuarios();
    @WebMethod Usuario buscarUsuario(@WebParam(name = "id") int id);
    @WebMethod Usuario criarUsuario(@WebParam(name = "nome") String nome, @WebParam(name = "idade") int idade);
    @WebMethod Usuario atualizarUsuario(@WebParam(name = "id") int id, @WebParam(name = "nome") String nome, @WebParam(name = "idade") int idade);
    @WebMethod String deletarUsuario(@WebParam(name = "id") int id);

    @WebMethod List<Playlist> listarPlaylists();
    @WebMethod Playlist buscarPlaylist(@WebParam(name = "id") int id);
    @WebMethod Playlist criarPlaylist(@WebParam(name = "nome") String nome, @WebParam(name = "usuario_id") int usuario_id);
    @WebMethod Playlist atualizarPlaylist(@WebParam(name = "id") int id, @WebParam(name = "nome") String nome, @WebParam(name = "usuario_id") int usuario_id);
    @WebMethod String deletarPlaylist(@WebParam(name = "id") int id);

    @WebMethod List<Musica> listarMusicasPlaylist(@WebParam(name = "playlist_id") int playlist_id);
    @WebMethod String adicionarMusicaNaPlaylist(@WebParam(name = "playlist_id") int playlist_id, @WebParam(name = "musica_id") int musica_id);
}
