package main

import (
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	banco "trab6/go/db"
)

type Musica struct {
	ID      int    `xml:"id"`
	Nome    string `xml:"nome"`
	Artista string `xml:"artista"`
}

type Usuario struct {
	ID    int    `xml:"id"`
	Nome  string `xml:"nome"`
	Idade int    `xml:"idade"`
}

type Playlist struct {
	ID        int    `xml:"id"`
	Nome      string `xml:"nome"`
	UsuarioID int    `xml:"usuario_id"`
}

type EnvelopeEntrada struct {
	XMLName xml.Name     `xml:"Envelope"`
	Body    CorpoEntrada `xml:"Body"`
}

type CorpoEntrada struct {
	Conteudo string `xml:",innerxml"`
}

type ListarMusicasRequest struct {
	XMLName xml.Name `xml:"ListarMusicasRequest"`
}

type BuscarMusicaRequest struct {
	XMLName xml.Name `xml:"BuscarMusicaRequest"`
	ID      int      `xml:"id"`
}

type CriarMusicaRequest struct {
	XMLName xml.Name `xml:"CriarMusicaRequest"`
	Nome    string   `xml:"nome"`
	Artista string   `xml:"artista"`
}

type AtualizarMusicaRequest struct {
	XMLName xml.Name `xml:"AtualizarMusicaRequest"`
	ID      int      `xml:"id"`
	Nome    string   `xml:"nome"`
	Artista string   `xml:"artista"`
}

type DeletarMusicaRequest struct {
	XMLName xml.Name `xml:"DeletarMusicaRequest"`
	ID      int      `xml:"id"`
}

type ListarUsuariosRequest struct {
	XMLName xml.Name `xml:"ListarUsuariosRequest"`
}

type BuscarUsuarioRequest struct {
	XMLName xml.Name `xml:"BuscarUsuarioRequest"`
	ID      int      `xml:"id"`
}

type CriarUsuarioRequest struct {
	XMLName xml.Name `xml:"CriarUsuarioRequest"`
	Nome    string   `xml:"nome"`
	Idade   int      `xml:"idade"`
}

type AtualizarUsuarioRequest struct {
	XMLName xml.Name `xml:"AtualizarUsuarioRequest"`
	ID      int      `xml:"id"`
	Nome    string   `xml:"nome"`
	Idade   int      `xml:"idade"`
}

type DeletarUsuarioRequest struct {
	XMLName xml.Name `xml:"DeletarUsuarioRequest"`
	ID      int      `xml:"id"`
}

type ListarPlaylistsRequest struct {
	XMLName xml.Name `xml:"ListarPlaylistsRequest"`
}

type BuscarPlaylistRequest struct {
	XMLName xml.Name `xml:"BuscarPlaylistRequest"`
	ID      int      `xml:"id"`
}

type CriarPlaylistRequest struct {
	XMLName   xml.Name `xml:"CriarPlaylistRequest"`
	Nome      string   `xml:"nome"`
	UsuarioID int      `xml:"usuario_id"`
}

type AtualizarPlaylistRequest struct {
	XMLName   xml.Name `xml:"AtualizarPlaylistRequest"`
	ID        int      `xml:"id"`
	Nome      string   `xml:"nome"`
	UsuarioID int      `xml:"usuario_id"`
}

type DeletarPlaylistRequest struct {
	XMLName xml.Name `xml:"DeletarPlaylistRequest"`
	ID      int      `xml:"id"`
}

type ListarMusicasPlaylistRequest struct {
	XMLName    xml.Name `xml:"ListarMusicasPlaylistRequest"`
	PlaylistID int      `xml:"playlist_id"`
}

type AdicionarMusicaPlaylistRequest struct {
	XMLName    xml.Name `xml:"AdicionarMusicaPlaylistRequest"`
	PlaylistID int      `xml:"playlist_id"`
	MusicaID   int      `xml:"musica_id"`
}

type MusicasResponse struct {
	XMLName xml.Name `xml:"ListarMusicasResponse"`
	Musicas []Musica `xml:"musicas>musica"`
}

type MusicaResponse struct {
	XMLName xml.Name `xml:"MusicaResponse"`
	Musica  Musica   `xml:"musica"`
}

type UsuariosResponse struct {
	XMLName  xml.Name  `xml:"ListarUsuariosResponse"`
	Usuarios []Usuario `xml:"usuarios>usuario"`
}

type UsuarioResponse struct {
	XMLName xml.Name `xml:"UsuarioResponse"`
	Usuario Usuario  `xml:"usuario"`
}

type PlaylistsResponse struct {
	XMLName   xml.Name   `xml:"ListarPlaylistsResponse"`
	Playlists []Playlist `xml:"playlists>playlist"`
}

type PlaylistResponse struct {
	XMLName  xml.Name `xml:"PlaylistResponse"`
	Playlist Playlist `xml:"playlist"`
}

type MensagemResponse struct {
	XMLName  xml.Name `xml:"MensagemResponse"`
	Mensagem string   `xml:"mensagem"`
}

type ErroResponse struct {
	XMLName xml.Name `xml:"ErroResponse"`
	Erro    string   `xml:"erro"`
}

var conexao *sql.DB

func main() {
	conexao = banco.ConectarBanco()
	defer conexao.Close()

	http.HandleFunc("/soap", soapHandler)

	log.Println("Servidor SOAP rodando em http://localhost:8081/soap")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func soapHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responderErroSOAP(w, "Use POST para enviar requisições SOAP")
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		responderErroSOAP(w, "Erro ao ler corpo da requisição")
		return
	}

	var envelope EnvelopeEntrada
	err = xml.Unmarshal(bodyBytes, &envelope)
	if err != nil {
		responderErroSOAP(w, "XML inválido")
		return
	}

	conteudo := envelope.Body.Conteudo

	switch {
	case strings.Contains(conteudo, "ListarMusicasRequest"):
		listarMusicasSOAP(w)

	case strings.Contains(conteudo, "BuscarMusicaRequest"):
		var req BuscarMusicaRequest
		xml.Unmarshal([]byte(conteudo), &req)
		buscarMusicaSOAP(w, req.ID)

	case strings.Contains(conteudo, "CriarMusicaRequest"):
		var req CriarMusicaRequest
		xml.Unmarshal([]byte(conteudo), &req)
		criarMusicaSOAP(w, req)

	case strings.Contains(conteudo, "AtualizarMusicaRequest"):
		var req AtualizarMusicaRequest
		xml.Unmarshal([]byte(conteudo), &req)
		atualizarMusicaSOAP(w, req)

	case strings.Contains(conteudo, "DeletarMusicaRequest"):
		var req DeletarMusicaRequest
		xml.Unmarshal([]byte(conteudo), &req)
		deletarMusicaSOAP(w, req.ID)

	case strings.Contains(conteudo, "ListarUsuariosRequest"):
		listarUsuariosSOAP(w)

	case strings.Contains(conteudo, "BuscarUsuarioRequest"):
		var req BuscarUsuarioRequest
		xml.Unmarshal([]byte(conteudo), &req)
		buscarUsuarioSOAP(w, req.ID)

	case strings.Contains(conteudo, "CriarUsuarioRequest"):
		var req CriarUsuarioRequest
		xml.Unmarshal([]byte(conteudo), &req)
		criarUsuarioSOAP(w, req)

	case strings.Contains(conteudo, "AtualizarUsuarioRequest"):
		var req AtualizarUsuarioRequest
		xml.Unmarshal([]byte(conteudo), &req)
		atualizarUsuarioSOAP(w, req)

	case strings.Contains(conteudo, "DeletarUsuarioRequest"):
		var req DeletarUsuarioRequest
		xml.Unmarshal([]byte(conteudo), &req)
		deletarUsuarioSOAP(w, req.ID)

	case strings.Contains(conteudo, "ListarPlaylistsRequest"):
		listarPlaylistsSOAP(w)

	case strings.Contains(conteudo, "BuscarPlaylistRequest"):
		var req BuscarPlaylistRequest
		xml.Unmarshal([]byte(conteudo), &req)
		buscarPlaylistSOAP(w, req.ID)

	case strings.Contains(conteudo, "CriarPlaylistRequest"):
		var req CriarPlaylistRequest
		xml.Unmarshal([]byte(conteudo), &req)
		criarPlaylistSOAP(w, req)

	case strings.Contains(conteudo, "AtualizarPlaylistRequest"):
		var req AtualizarPlaylistRequest
		xml.Unmarshal([]byte(conteudo), &req)
		atualizarPlaylistSOAP(w, req)

	case strings.Contains(conteudo, "DeletarPlaylistRequest"):
		var req DeletarPlaylistRequest
		xml.Unmarshal([]byte(conteudo), &req)
		deletarPlaylistSOAP(w, req.ID)

	case strings.Contains(conteudo, "ListarMusicasPlaylistRequest"):
		var req ListarMusicasPlaylistRequest
		xml.Unmarshal([]byte(conteudo), &req)
		listarMusicasPlaylistSOAP(w, req.PlaylistID)

	case strings.Contains(conteudo, "AdicionarMusicaPlaylistRequest"):
		var req AdicionarMusicaPlaylistRequest
		xml.Unmarshal([]byte(conteudo), &req)
		adicionarMusicaPlaylistSOAP(w, req)

	default:
		responderErroSOAP(w, "Operação SOAP não reconhecida")
	}
}

func listarMusicasSOAP(w http.ResponseWriter) {
	rows, err := conexao.Query("SELECT id, nome, artista FROM musica ORDER BY id")
	if err != nil {
		responderErroSOAP(w, "Erro ao listar músicas")
		return
	}
	defer rows.Close()

	musicas := []Musica{}

	for rows.Next() {
		var musica Musica
		err := rows.Scan(&musica.ID, &musica.Nome, &musica.Artista)
		if err != nil {
			responderErroSOAP(w, "Erro ao ler música")
			return
		}
		musicas = append(musicas, musica)
	}

	responderSOAP(w, MusicasResponse{Musicas: musicas})
}

func buscarMusicaSOAP(w http.ResponseWriter, id int) {
	var musica Musica

	err := conexao.QueryRow(
		"SELECT id, nome, artista FROM musica WHERE id = $1",
		id,
	).Scan(&musica.ID, &musica.Nome, &musica.Artista)

	if err == sql.ErrNoRows {
		responderErroSOAP(w, "Música não encontrada")
		return
	}

	if err != nil {
		responderErroSOAP(w, "Erro ao buscar música")
		return
	}

	responderSOAP(w, MusicaResponse{Musica: musica})
}

func criarMusicaSOAP(w http.ResponseWriter, req CriarMusicaRequest) {
	var musica Musica

	musica.Nome = req.Nome
	musica.Artista = req.Artista

	err := conexao.QueryRow(
		"INSERT INTO musica (nome, artista) VALUES ($1, $2) RETURNING id",
		musica.Nome,
		musica.Artista,
	).Scan(&musica.ID)

	if err != nil {
		responderErroSOAP(w, "Erro ao criar música")
		return
	}

	responderSOAP(w, MusicaResponse{Musica: musica})
}

func atualizarMusicaSOAP(w http.ResponseWriter, req AtualizarMusicaRequest) {
	resultado, err := conexao.Exec(
		"UPDATE musica SET nome = $1, artista = $2 WHERE id = $3",
		req.Nome,
		req.Artista,
		req.ID,
	)

	if err != nil {
		responderErroSOAP(w, "Erro ao atualizar música")
		return
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		responderErroSOAP(w, "Música não encontrada")
		return
	}

	responderSOAP(w, MusicaResponse{
		Musica: Musica{
			ID:      req.ID,
			Nome:    req.Nome,
			Artista: req.Artista,
		},
	})
}

func deletarMusicaSOAP(w http.ResponseWriter, id int) {
	resultado, err := conexao.Exec("DELETE FROM musica WHERE id = $1", id)
	if err != nil {
		responderErroSOAP(w, "Erro ao deletar música")
		return
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		responderErroSOAP(w, "Música não encontrada")
		return
	}

	responderSOAP(w, MensagemResponse{Mensagem: "Música deletada com sucesso"})
}

func listarUsuariosSOAP(w http.ResponseWriter) {
	rows, err := conexao.Query("SELECT id, nome, idade FROM usuario ORDER BY id")
	if err != nil {
		responderErroSOAP(w, "Erro ao listar usuários")
		return
	}
	defer rows.Close()

	usuarios := []Usuario{}

	for rows.Next() {
		var usuario Usuario
		err := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Idade)
		if err != nil {
			responderErroSOAP(w, "Erro ao ler usuário")
			return
		}
		usuarios = append(usuarios, usuario)
	}

	responderSOAP(w, UsuariosResponse{Usuarios: usuarios})
}

func buscarUsuarioSOAP(w http.ResponseWriter, id int) {
	var usuario Usuario

	err := conexao.QueryRow(
		"SELECT id, nome, idade FROM usuario WHERE id = $1",
		id,
	).Scan(&usuario.ID, &usuario.Nome, &usuario.Idade)

	if err == sql.ErrNoRows {
		responderErroSOAP(w, "Usuário não encontrado")
		return
	}

	if err != nil {
		responderErroSOAP(w, "Erro ao buscar usuário")
		return
	}

	responderSOAP(w, UsuarioResponse{Usuario: usuario})
}

func criarUsuarioSOAP(w http.ResponseWriter, req CriarUsuarioRequest) {
	var usuario Usuario

	usuario.Nome = req.Nome
	usuario.Idade = req.Idade

	err := conexao.QueryRow(
		"INSERT INTO usuario (nome, idade) VALUES ($1, $2) RETURNING id",
		usuario.Nome,
		usuario.Idade,
	).Scan(&usuario.ID)

	if err != nil {
		responderErroSOAP(w, "Erro ao criar usuário")
		return
	}

	responderSOAP(w, UsuarioResponse{Usuario: usuario})
}

func atualizarUsuarioSOAP(w http.ResponseWriter, req AtualizarUsuarioRequest) {
	resultado, err := conexao.Exec(
		"UPDATE usuario SET nome = $1, idade = $2 WHERE id = $3",
		req.Nome,
		req.Idade,
		req.ID,
	)

	if err != nil {
		responderErroSOAP(w, "Erro ao atualizar usuário")
		return
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		responderErroSOAP(w, "Usuário não encontrado")
		return
	}

	responderSOAP(w, UsuarioResponse{
		Usuario: Usuario{
			ID:    req.ID,
			Nome:  req.Nome,
			Idade: req.Idade,
		},
	})
}

func deletarUsuarioSOAP(w http.ResponseWriter, id int) {
	resultado, err := conexao.Exec("DELETE FROM usuario WHERE id = $1", id)
	if err != nil {
		responderErroSOAP(w, "Erro ao deletar usuário")
		return
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		responderErroSOAP(w, "Usuário não encontrado")
		return
	}

	responderSOAP(w, MensagemResponse{Mensagem: "Usuário deletado com sucesso"})
}

func listarPlaylistsSOAP(w http.ResponseWriter) {
	rows, err := conexao.Query("SELECT id, nome, usuario_id FROM playlist ORDER BY id")
	if err != nil {
		responderErroSOAP(w, "Erro ao listar playlists")
		return
	}
	defer rows.Close()

	playlists := []Playlist{}

	for rows.Next() {
		var playlist Playlist
		err := rows.Scan(&playlist.ID, &playlist.Nome, &playlist.UsuarioID)
		if err != nil {
			responderErroSOAP(w, "Erro ao ler playlist")
			return
		}
		playlists = append(playlists, playlist)
	}

	responderSOAP(w, PlaylistsResponse{Playlists: playlists})
}

func buscarPlaylistSOAP(w http.ResponseWriter, id int) {
	var playlist Playlist

	err := conexao.QueryRow(
		"SELECT id, nome, usuario_id FROM playlist WHERE id = $1",
		id,
	).Scan(&playlist.ID, &playlist.Nome, &playlist.UsuarioID)

	if err == sql.ErrNoRows {
		responderErroSOAP(w, "Playlist não encontrada")
		return
	}

	if err != nil {
		responderErroSOAP(w, "Erro ao buscar playlist")
		return
	}

	responderSOAP(w, PlaylistResponse{Playlist: playlist})
}

func criarPlaylistSOAP(w http.ResponseWriter, req CriarPlaylistRequest) {
	var playlist Playlist

	playlist.Nome = req.Nome
	playlist.UsuarioID = req.UsuarioID

	err := conexao.QueryRow(
		"INSERT INTO playlist (nome, usuario_id) VALUES ($1, $2) RETURNING id",
		playlist.Nome,
		playlist.UsuarioID,
	).Scan(&playlist.ID)

	if err != nil {
		responderErroSOAP(w, "Erro ao criar playlist")
		return
	}

	responderSOAP(w, PlaylistResponse{Playlist: playlist})
}

func atualizarPlaylistSOAP(w http.ResponseWriter, req AtualizarPlaylistRequest) {
	resultado, err := conexao.Exec(
		"UPDATE playlist SET nome = $1, usuario_id = $2 WHERE id = $3",
		req.Nome,
		req.UsuarioID,
		req.ID,
	)

	if err != nil {
		responderErroSOAP(w, "Erro ao atualizar playlist")
		return
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		responderErroSOAP(w, "Playlist não encontrada")
		return
	}

	responderSOAP(w, PlaylistResponse{
		Playlist: Playlist{
			ID:        req.ID,
			Nome:      req.Nome,
			UsuarioID: req.UsuarioID,
		},
	})
}

func deletarPlaylistSOAP(w http.ResponseWriter, id int) {
	resultado, err := conexao.Exec("DELETE FROM playlist WHERE id = $1", id)
	if err != nil {
		responderErroSOAP(w, "Erro ao deletar playlist")
		return
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		responderErroSOAP(w, "Playlist não encontrada")
		return
	}

	responderSOAP(w, MensagemResponse{Mensagem: "Playlist deletada com sucesso"})
}

func listarMusicasPlaylistSOAP(w http.ResponseWriter, playlistID int) {
	rows, err := conexao.Query(`
		SELECT m.id, m.nome, m.artista
		FROM musica m
		INNER JOIN playlist_musica pm ON pm.musica_id = m.id
		WHERE pm.playlist_id = $1
		ORDER BY m.id
	`, playlistID)

	if err != nil {
		responderErroSOAP(w, "Erro ao listar músicas da playlist")
		return
	}
	defer rows.Close()

	musicas := []Musica{}

	for rows.Next() {
		var musica Musica
		err := rows.Scan(&musica.ID, &musica.Nome, &musica.Artista)
		if err != nil {
			responderErroSOAP(w, "Erro ao ler música")
			return
		}
		musicas = append(musicas, musica)
	}

	responderSOAP(w, MusicasResponse{Musicas: musicas})
}

func adicionarMusicaPlaylistSOAP(w http.ResponseWriter, req AdicionarMusicaPlaylistRequest) {
	_, err := conexao.Exec(
		"INSERT INTO playlist_musica (playlist_id, musica_id) VALUES ($1, $2)",
		req.PlaylistID,
		req.MusicaID,
	)

	if err != nil {
		responderErroSOAP(w, "Erro ao adicionar música na playlist")
		return
	}

	responderSOAP(w, MensagemResponse{Mensagem: "Música adicionada na playlist com sucesso"})
}

func responderSOAP(w http.ResponseWriter, resposta interface{}) {
	xmlResposta, err := xml.MarshalIndent(resposta, "", "  ")
	if err != nil {
		responderErroSOAP(w, "Erro ao gerar XML de resposta")
		return
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
	w.Write([]byte(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">`))
	w.Write([]byte(`<soap:Body>`))
	w.Write(xmlResposta)
	w.Write([]byte(`</soap:Body>`))
	w.Write([]byte(`</soap:Envelope>`))
}

func responderErroSOAP(w http.ResponseWriter, mensagem string) {
	resposta := ErroResponse{Erro: mensagem}

	xmlResposta, _ := xml.MarshalIndent(resposta, "", "  ")

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
	w.Write([]byte(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">`))
	w.Write([]byte(`<soap:Body>`))
	w.Write(xmlResposta)
	w.Write([]byte(`</soap:Body>`))
	w.Write([]byte(`</soap:Envelope>`))
}

func converterParaInt(valor string) int {
	numero, _ := strconv.Atoi(valor)
	return numero
}