package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	banco "trab6/go/db"
)

type Musica struct {
	ID      int    `json:"id"`
	Nome    string `json:"nome"`
	Artista string `json:"artista"`
}

type Usuario struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Idade int    `json:"idade"`
}

type Playlist struct {
	ID        int    `json:"id"`
	Nome      string `json:"nome"`
	UsuarioID int    `json:"usuario_id"`
}

type AdicionarMusicaPlaylistRequest struct {
	MusicaID int `json:"musica_id"`
}

var conexao *sql.DB

func main() {
	conexao = banco.ConectarBanco()
	defer conexao.Close()

	http.HandleFunc("/", rotaInicial)

	http.HandleFunc("/musicas", musicasHandler)
	http.HandleFunc("/musicas/", musicaPorIDHandler)

	http.HandleFunc("/usuarios", usuariosHandler)
	http.HandleFunc("/usuarios/", usuarioPorIDHandler)

	http.HandleFunc("/playlists", playlistsHandler)
	http.HandleFunc("/playlists/", playlistPorIDHandler)

	log.Println("Servidor REST rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rotaInicial(w http.ResponseWriter, r *http.Request) {
	responderJSON(w, http.StatusOK, map[string]string{
		"mensagem": "API REST de músicas funcionando",
	})
}

func musicasHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listarMusicas(w)
	case http.MethodPost:
		criarMusica(w, r)
	default:
		responderErro(w, http.StatusMethodNotAllowed, "Método não permitido")
	}
}

func musicaPorIDHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := extrairID(r.URL.Path, "/musicas/")
	if !ok {
		responderErro(w, http.StatusBadRequest, "ID inválido")
		return
	}

	switch r.Method {
	case http.MethodGet:
		buscarMusicaPorID(w, id)
	case http.MethodPut:
		atualizarMusica(w, r, id)
	case http.MethodDelete:
		deletarMusica(w, id)
	default:
		responderErro(w, http.StatusMethodNotAllowed, "Método não permitido")
	}
}

func listarMusicas(w http.ResponseWriter) {
	rows, err := conexao.Query("SELECT id, nome, artista FROM musica ORDER BY id")
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao listar músicas")
		return
	}
	defer rows.Close()

	musicas := []Musica{}

	for rows.Next() {
		var musica Musica

		err := rows.Scan(&musica.ID, &musica.Nome, &musica.Artista)
		if err != nil {
			responderErro(w, http.StatusInternalServerError, "Erro ao ler música")
			return
		}

		musicas = append(musicas, musica)
	}

	responderJSON(w, http.StatusOK, musicas)
}

func buscarMusicaPorID(w http.ResponseWriter, id int) {
	var musica Musica

	err := conexao.QueryRow(
		"SELECT id, nome, artista FROM musica WHERE id = $1",
		id,
	).Scan(&musica.ID, &musica.Nome, &musica.Artista)

	if err == sql.ErrNoRows {
		responderErro(w, http.StatusNotFound, "Música não encontrada")
		return
	}

	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao buscar música")
		return
	}

	responderJSON(w, http.StatusOK, musica)
}

func criarMusica(w http.ResponseWriter, r *http.Request) {
	var musica Musica

	err := json.NewDecoder(r.Body).Decode(&musica)
	if err != nil {
		responderErro(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	err = conexao.QueryRow(
		"INSERT INTO musica (nome, artista) VALUES ($1, $2) RETURNING id",
		musica.Nome,
		musica.Artista,
	).Scan(&musica.ID)

	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao criar música")
		return
	}

	responderJSON(w, http.StatusCreated, musica)
}

func atualizarMusica(w http.ResponseWriter, r *http.Request, id int) {
	var musica Musica

	err := json.NewDecoder(r.Body).Decode(&musica)
	if err != nil {
		responderErro(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	resultado, err := conexao.Exec(
		"UPDATE musica SET nome = $1, artista = $2 WHERE id = $3",
		musica.Nome,
		musica.Artista,
		id,
	)

	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao atualizar música")
		return
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		responderErro(w, http.StatusNotFound, "Música não encontrada")
		return
	}

	musica.ID = id
	responderJSON(w, http.StatusOK, musica)
}

func deletarMusica(w http.ResponseWriter, id int) {
	resultado, err := conexao.Exec("DELETE FROM musica WHERE id = $1", id)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao deletar música")
		return
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		responderErro(w, http.StatusNotFound, "Música não encontrada")
		return
	}

	responderJSON(w, http.StatusOK, map[string]string{
		"mensagem": "Música deletada com sucesso",
	})
}

func usuariosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listarUsuarios(w)
	case http.MethodPost:
		criarUsuario(w, r)
	default:
		responderErro(w, http.StatusMethodNotAllowed, "Método não permitido")
	}
}

func usuarioPorIDHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := extrairID(r.URL.Path, "/usuarios/")
	if !ok {
		responderErro(w, http.StatusBadRequest, "ID inválido")
		return
	}

	switch r.Method {
	case http.MethodGet:
		buscarUsuarioPorID(w, id)
	case http.MethodPut:
		atualizarUsuario(w, r, id)
	case http.MethodDelete:
		deletarUsuario(w, id)
	default:
		responderErro(w, http.StatusMethodNotAllowed, "Método não permitido")
	}
}

func listarUsuarios(w http.ResponseWriter) {
	rows, err := conexao.Query("SELECT id, nome, idade FROM usuario ORDER BY id")
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao listar usuários")
		return
	}
	defer rows.Close()

	usuarios := []Usuario{}

	for rows.Next() {
		var usuario Usuario

		err := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Idade)
		if err != nil {
			responderErro(w, http.StatusInternalServerError, "Erro ao ler usuário")
			return
		}

		usuarios = append(usuarios, usuario)
	}

	responderJSON(w, http.StatusOK, usuarios)
}

func buscarUsuarioPorID(w http.ResponseWriter, id int) {
	var usuario Usuario

	err := conexao.QueryRow(
		"SELECT id, nome, idade FROM usuario WHERE id = $1",
		id,
	).Scan(&usuario.ID, &usuario.Nome, &usuario.Idade)

	if err == sql.ErrNoRows {
		responderErro(w, http.StatusNotFound, "Usuário não encontrado")
		return
	}

	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao buscar usuário")
		return
	}

	responderJSON(w, http.StatusOK, usuario)
}

func criarUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario Usuario

	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		responderErro(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	err = conexao.QueryRow(
		"INSERT INTO usuario (nome, idade) VALUES ($1, $2) RETURNING id",
		usuario.Nome,
		usuario.Idade,
	).Scan(&usuario.ID)

	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao criar usuário")
		return
	}

	responderJSON(w, http.StatusCreated, usuario)
}

func atualizarUsuario(w http.ResponseWriter, r *http.Request, id int) {
	var usuario Usuario

	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		responderErro(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	resultado, err := conexao.Exec(
		"UPDATE usuario SET nome = $1, idade = $2 WHERE id = $3",
		usuario.Nome,
		usuario.Idade,
		id,
	)

	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao atualizar usuário")
		return
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		responderErro(w, http.StatusNotFound, "Usuário não encontrado")
		return
	}

	usuario.ID = id
	responderJSON(w, http.StatusOK, usuario)
}

func deletarUsuario(w http.ResponseWriter, id int) {
	resultado, err := conexao.Exec("DELETE FROM usuario WHERE id = $1", id)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao deletar usuário")
		return
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		responderErro(w, http.StatusNotFound, "Usuário não encontrado")
		return
	}

	responderJSON(w, http.StatusOK, map[string]string{
		"mensagem": "Usuário deletado com sucesso",
	})
}

func playlistsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listarPlaylists(w)
	case http.MethodPost:
		criarPlaylist(w, r)
	default:
		responderErro(w, http.StatusMethodNotAllowed, "Método não permitido")
	}
}

func playlistPorIDHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/musicas") {
		id, ok := extrairIDPlaylistMusicas(r.URL.Path)
		if !ok {
			responderErro(w, http.StatusBadRequest, "ID inválido")
			return
		}

		switch r.Method {
		case http.MethodGet:
			listarMusicasDaPlaylist(w, id)
		case http.MethodPost:
			adicionarMusicaNaPlaylist(w, r, id)
		default:
			responderErro(w, http.StatusMethodNotAllowed, "Método não permitido")
		}

		return
	}

	id, ok := extrairID(r.URL.Path, "/playlists/")
	if !ok {
		responderErro(w, http.StatusBadRequest, "ID inválido")
		return
	}

	switch r.Method {
	case http.MethodGet:
		buscarPlaylistPorID(w, id)
	case http.MethodPut:
		atualizarPlaylist(w, r, id)
	case http.MethodDelete:
		deletarPlaylist(w, id)
	default:
		responderErro(w, http.StatusMethodNotAllowed, "Método não permitido")
	}
}

func listarPlaylists(w http.ResponseWriter) {
	rows, err := conexao.Query("SELECT id, nome, usuario_id FROM playlist ORDER BY id")
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao listar playlists")
		return
	}
	defer rows.Close()

	playlists := []Playlist{}

	for rows.Next() {
		var playlist Playlist

		err := rows.Scan(&playlist.ID, &playlist.Nome, &playlist.UsuarioID)
		if err != nil {
			responderErro(w, http.StatusInternalServerError, "Erro ao ler playlist")
			return
		}

		playlists = append(playlists, playlist)
	}

	responderJSON(w, http.StatusOK, playlists)
}

func buscarPlaylistPorID(w http.ResponseWriter, id int) {
	var playlist Playlist

	err := conexao.QueryRow(
		"SELECT id, nome, usuario_id FROM playlist WHERE id = $1",
		id,
	).Scan(&playlist.ID, &playlist.Nome, &playlist.UsuarioID)

	if err == sql.ErrNoRows {
		responderErro(w, http.StatusNotFound, "Playlist não encontrada")
		return
	}

	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao buscar playlist")
		return
	}

	responderJSON(w, http.StatusOK, playlist)
}

func criarPlaylist(w http.ResponseWriter, r *http.Request) {
	var playlist Playlist

	err := json.NewDecoder(r.Body).Decode(&playlist)
	if err != nil {
		responderErro(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	err = conexao.QueryRow(
		"INSERT INTO playlist (nome, usuario_id) VALUES ($1, $2) RETURNING id",
		playlist.Nome,
		playlist.UsuarioID,
	).Scan(&playlist.ID)

	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao criar playlist")
		return
	}

	responderJSON(w, http.StatusCreated, playlist)
}

func atualizarPlaylist(w http.ResponseWriter, r *http.Request, id int) {
	var playlist Playlist

	err := json.NewDecoder(r.Body).Decode(&playlist)
	if err != nil {
		responderErro(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	resultado, err := conexao.Exec(
		"UPDATE playlist SET nome = $1, usuario_id = $2 WHERE id = $3",
		playlist.Nome,
		playlist.UsuarioID,
		id,
	)

	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao atualizar playlist")
		return
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		responderErro(w, http.StatusNotFound, "Playlist não encontrada")
		return
	}

	playlist.ID = id
	responderJSON(w, http.StatusOK, playlist)
}

func deletarPlaylist(w http.ResponseWriter, id int) {
	resultado, err := conexao.Exec("DELETE FROM playlist WHERE id = $1", id)
	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao deletar playlist")
		return
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		responderErro(w, http.StatusNotFound, "Playlist não encontrada")
		return
	}

	responderJSON(w, http.StatusOK, map[string]string{
		"mensagem": "Playlist deletada com sucesso",
	})
}

func listarMusicasDaPlaylist(w http.ResponseWriter, playlistID int) {
	rows, err := conexao.Query(`
		SELECT m.id, m.nome, m.artista
		FROM musica m
		INNER JOIN playlist_musica pm ON pm.musica_id = m.id
		WHERE pm.playlist_id = $1
		ORDER BY m.id
	`, playlistID)

	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao listar músicas da playlist")
		return
	}
	defer rows.Close()

	musicas := []Musica{}

	for rows.Next() {
		var musica Musica

		err := rows.Scan(&musica.ID, &musica.Nome, &musica.Artista)
		if err != nil {
			responderErro(w, http.StatusInternalServerError, "Erro ao ler música")
			return
		}

		musicas = append(musicas, musica)
	}

	responderJSON(w, http.StatusOK, musicas)
}

func adicionarMusicaNaPlaylist(w http.ResponseWriter, r *http.Request, playlistID int) {
	var request AdicionarMusicaPlaylistRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		responderErro(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	_, err = conexao.Exec(
		"INSERT INTO playlist_musica (playlist_id, musica_id) VALUES ($1, $2)",
		playlistID,
		request.MusicaID,
	)

	if err != nil {
		responderErro(w, http.StatusInternalServerError, "Erro ao adicionar música na playlist")
		return
	}

	responderJSON(w, http.StatusCreated, map[string]string{
		"mensagem": "Música adicionada na playlist com sucesso",
	})
}

func extrairID(path string, prefixo string) (int, bool) {
	idTexto := strings.TrimPrefix(path, prefixo)
	idTexto = strings.Trim(idTexto, "/")

	if idTexto == "" {
		return 0, false
	}

	id, err := strconv.Atoi(idTexto)
	if err != nil {
		return 0, false
	}

	return id, true
}

func extrairIDPlaylistMusicas(path string) (int, bool) {
	partes := strings.Split(strings.Trim(path, "/"), "/")

	if len(partes) != 3 {
		return 0, false
	}

	if partes[0] != "playlists" || partes[2] != "musicas" {
		return 0, false
	}

	id, err := strconv.Atoi(partes[1])
	if err != nil {
		return 0, false
	}

	return id, true
}

func responderJSON(w http.ResponseWriter, status int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dados)
}

func responderErro(w http.ResponseWriter, status int, mensagem string) {
	responderJSON(w, status, map[string]string{
		"erro": mensagem,
	})
}