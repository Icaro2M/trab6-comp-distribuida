package main

import (
	"context"
	"database/sql"
	"log"
	"net"

	musicpb "trab6/go/grpc/proto"

	banco "trab6/go/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type servidor struct {
	musicpb.UnimplementedMusicServiceServer
	conexao *sql.DB
}

func main() {
	conexao := banco.ConectarBanco()
	defer conexao.Close()

	listener, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatal("Erro ao abrir porta 8083:", err)
	}

	grpcServer := grpc.NewServer()

	musicpb.RegisterMusicServiceServer(grpcServer, &servidor{
		conexao: conexao,
	})

	log.Println("Servidor gRPC rodando na porta 8083")

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Erro ao iniciar servidor gRPC:", err)
	}
}

func (s *servidor) ListMusics(ctx context.Context, req *musicpb.Empty) (*musicpb.MusicListResponse, error) {
	rows, err := s.conexao.Query("SELECT id, nome, artista FROM musica ORDER BY id")
	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao listar músicas")
	}
	defer rows.Close()

	musicas := []*musicpb.Music{}

	for rows.Next() {
		var musica musicpb.Music

		err := rows.Scan(&musica.Id, &musica.Nome, &musica.Artista)
		if err != nil {
			return nil, status.Error(codes.Internal, "Erro ao ler música")
		}

		musicas = append(musicas, &musica)
	}

	return &musicpb.MusicListResponse{
		Musics: musicas,
	}, nil
}

func (s *servidor) GetMusic(ctx context.Context, req *musicpb.IdRequest) (*musicpb.Music, error) {
	var musica musicpb.Music

	err := s.conexao.QueryRow(
		"SELECT id, nome, artista FROM musica WHERE id = $1",
		req.Id,
	).Scan(&musica.Id, &musica.Nome, &musica.Artista)

	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, "Música não encontrada")
	}

	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao buscar música")
	}

	return &musica, nil
}

func (s *servidor) CreateMusic(ctx context.Context, req *musicpb.CreateMusicRequest) (*musicpb.Music, error) {
	var musica musicpb.Music

	musica.Nome = req.Nome
	musica.Artista = req.Artista

	err := s.conexao.QueryRow(
		"INSERT INTO musica (nome, artista) VALUES ($1, $2) RETURNING id",
		musica.Nome,
		musica.Artista,
	).Scan(&musica.Id)

	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao criar música")
	}

	return &musica, nil
}

func (s *servidor) UpdateMusic(ctx context.Context, req *musicpb.UpdateMusicRequest) (*musicpb.Music, error) {
	resultado, err := s.conexao.Exec(
		"UPDATE musica SET nome = $1, artista = $2 WHERE id = $3",
		req.Nome,
		req.Artista,
		req.Id,
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao atualizar música")
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		return nil, status.Error(codes.NotFound, "Música não encontrada")
	}

	return &musicpb.Music{
		Id:      req.Id,
		Nome:    req.Nome,
		Artista: req.Artista,
	}, nil
}

func (s *servidor) DeleteMusic(ctx context.Context, req *musicpb.IdRequest) (*musicpb.MessageResponse, error) {
	resultado, err := s.conexao.Exec("DELETE FROM musica WHERE id = $1", req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao deletar música")
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		return nil, status.Error(codes.NotFound, "Música não encontrada")
	}

	return &musicpb.MessageResponse{
		Message: "Música deletada com sucesso",
	}, nil
}

func (s *servidor) ListUsers(ctx context.Context, req *musicpb.Empty) (*musicpb.UserListResponse, error) {
	rows, err := s.conexao.Query("SELECT id, nome, idade FROM usuario ORDER BY id")
	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao listar usuários")
	}
	defer rows.Close()

	usuarios := []*musicpb.User{}

	for rows.Next() {
		var usuario musicpb.User

		err := rows.Scan(&usuario.Id, &usuario.Nome, &usuario.Idade)
		if err != nil {
			return nil, status.Error(codes.Internal, "Erro ao ler usuário")
		}

		usuarios = append(usuarios, &usuario)
	}

	return &musicpb.UserListResponse{
		Users: usuarios,
	}, nil
}

func (s *servidor) GetUser(ctx context.Context, req *musicpb.IdRequest) (*musicpb.User, error) {
	var usuario musicpb.User

	err := s.conexao.QueryRow(
		"SELECT id, nome, idade FROM usuario WHERE id = $1",
		req.Id,
	).Scan(&usuario.Id, &usuario.Nome, &usuario.Idade)

	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, "Usuário não encontrado")
	}

	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao buscar usuário")
	}

	return &usuario, nil
}

func (s *servidor) CreateUser(ctx context.Context, req *musicpb.CreateUserRequest) (*musicpb.User, error) {
	var usuario musicpb.User

	usuario.Nome = req.Nome
	usuario.Idade = req.Idade

	err := s.conexao.QueryRow(
		"INSERT INTO usuario (nome, idade) VALUES ($1, $2) RETURNING id",
		usuario.Nome,
		usuario.Idade,
	).Scan(&usuario.Id)

	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao criar usuário")
	}

	return &usuario, nil
}

func (s *servidor) UpdateUser(ctx context.Context, req *musicpb.UpdateUserRequest) (*musicpb.User, error) {
	resultado, err := s.conexao.Exec(
		"UPDATE usuario SET nome = $1, idade = $2 WHERE id = $3",
		req.Nome,
		req.Idade,
		req.Id,
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao atualizar usuário")
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		return nil, status.Error(codes.NotFound, "Usuário não encontrado")
	}

	return &musicpb.User{
		Id:    req.Id,
		Nome:  req.Nome,
		Idade: req.Idade,
	}, nil
}

func (s *servidor) DeleteUser(ctx context.Context, req *musicpb.IdRequest) (*musicpb.MessageResponse, error) {
	resultado, err := s.conexao.Exec("DELETE FROM usuario WHERE id = $1", req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao deletar usuário")
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		return nil, status.Error(codes.NotFound, "Usuário não encontrado")
	}

	return &musicpb.MessageResponse{
		Message: "Usuário deletado com sucesso",
	}, nil
}

func (s *servidor) ListPlaylists(ctx context.Context, req *musicpb.Empty) (*musicpb.PlaylistListResponse, error) {
	rows, err := s.conexao.Query("SELECT id, nome, usuario_id FROM playlist ORDER BY id")
	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao listar playlists")
	}
	defer rows.Close()

	playlists := []*musicpb.Playlist{}

	for rows.Next() {
		var playlist musicpb.Playlist

		err := rows.Scan(&playlist.Id, &playlist.Nome, &playlist.UsuarioId)
		if err != nil {
			return nil, status.Error(codes.Internal, "Erro ao ler playlist")
		}

		playlists = append(playlists, &playlist)
	}

	return &musicpb.PlaylistListResponse{
		Playlists: playlists,
	}, nil
}

func (s *servidor) GetPlaylist(ctx context.Context, req *musicpb.IdRequest) (*musicpb.Playlist, error) {
	var playlist musicpb.Playlist

	err := s.conexao.QueryRow(
		"SELECT id, nome, usuario_id FROM playlist WHERE id = $1",
		req.Id,
	).Scan(&playlist.Id, &playlist.Nome, &playlist.UsuarioId)

	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, "Playlist não encontrada")
	}

	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao buscar playlist")
	}

	return &playlist, nil
}

func (s *servidor) CreatePlaylist(ctx context.Context, req *musicpb.CreatePlaylistRequest) (*musicpb.Playlist, error) {
	var playlist musicpb.Playlist

	playlist.Nome = req.Nome
	playlist.UsuarioId = req.UsuarioId

	err := s.conexao.QueryRow(
		"INSERT INTO playlist (nome, usuario_id) VALUES ($1, $2) RETURNING id",
		playlist.Nome,
		playlist.UsuarioId,
	).Scan(&playlist.Id)

	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao criar playlist")
	}

	return &playlist, nil
}

func (s *servidor) UpdatePlaylist(ctx context.Context, req *musicpb.UpdatePlaylistRequest) (*musicpb.Playlist, error) {
	resultado, err := s.conexao.Exec(
		"UPDATE playlist SET nome = $1, usuario_id = $2 WHERE id = $3",
		req.Nome,
		req.UsuarioId,
		req.Id,
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao atualizar playlist")
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		return nil, status.Error(codes.NotFound, "Playlist não encontrada")
	}

	return &musicpb.Playlist{
		Id:        req.Id,
		Nome:      req.Nome,
		UsuarioId: req.UsuarioId,
	}, nil
}

func (s *servidor) DeletePlaylist(ctx context.Context, req *musicpb.IdRequest) (*musicpb.MessageResponse, error) {
	resultado, err := s.conexao.Exec("DELETE FROM playlist WHERE id = $1", req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao deletar playlist")
	}

	linhasAfetadas, _ := resultado.RowsAffected()
	if linhasAfetadas == 0 {
		return nil, status.Error(codes.NotFound, "Playlist não encontrada")
	}

	return &musicpb.MessageResponse{
		Message: "Playlist deletada com sucesso",
	}, nil
}

func (s *servidor) AddMusicToPlaylist(ctx context.Context, req *musicpb.AddMusicToPlaylistRequest) (*musicpb.MessageResponse, error) {
	_, err := s.conexao.Exec(
		"INSERT INTO playlist_musica (playlist_id, musica_id) VALUES ($1, $2)",
		req.PlaylistId,
		req.MusicaId,
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao adicionar música na playlist")
	}

	return &musicpb.MessageResponse{
		Message: "Música adicionada na playlist com sucesso",
	}, nil
}

func (s *servidor) ListPlaylistMusics(ctx context.Context, req *musicpb.PlaylistMusicsRequest) (*musicpb.MusicListResponse, error) {
	rows, err := s.conexao.Query(`
		SELECT m.id, m.nome, m.artista
		FROM musica m
		INNER JOIN playlist_musica pm ON pm.musica_id = m.id
		WHERE pm.playlist_id = $1
		ORDER BY m.id
	`, req.PlaylistId)

	if err != nil {
		return nil, status.Error(codes.Internal, "Erro ao listar músicas da playlist")
	}
	defer rows.Close()

	musicas := []*musicpb.Music{}

	for rows.Next() {
		var musica musicpb.Music

		err := rows.Scan(&musica.Id, &musica.Nome, &musica.Artista)
		if err != nil {
			return nil, status.Error(codes.Internal, "Erro ao ler música")
		}

		musicas = append(musicas, &musica)
	}

	return &musicpb.MusicListResponse{
		Musics: musicas,
	}, nil
}