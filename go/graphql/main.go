package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"

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

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variable map[string]interface{} `json:"variables"`
}

var conexao *sql.DB

func main() {
	conexao = banco.ConectarBanco()
	defer conexao.Close()

	schema, err := criarSchema()
	if err != nil {
		log.Fatal("Erro ao criar schema GraphQL:", err)
	}

	http.HandleFunc("/", rotaInicial)
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		graphqlHandler(w, r, schema)
	})

	log.Println("Servidor GraphQL rodando em http://localhost:8082/graphql")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func rotaInicial(w http.ResponseWriter, r *http.Request) {
	responderJSON(w, http.StatusOK, map[string]string{
		"mensagem": "API GraphQL de músicas funcionando",
		"endpoint": "/graphql",
	})
}

func criarSchema() (graphql.Schema, error) {
	musicaType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Musica",
		Fields: graphql.Fields{
			"id":      &graphql.Field{Type: graphql.Int},
			"nome":    &graphql.Field{Type: graphql.String},
			"artista": &graphql.Field{Type: graphql.String},
		},
	})

	usuarioType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Usuario",
		Fields: graphql.Fields{
			"id":    &graphql.Field{Type: graphql.Int},
			"nome":  &graphql.Field{Type: graphql.String},
			"idade": &graphql.Field{Type: graphql.Int},
		},
	})

	playlistType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Playlist",
		Fields: graphql.Fields{
			"id":         &graphql.Field{Type: graphql.Int},
			"nome":       &graphql.Field{Type: graphql.String},
			"usuario_id": &graphql.Field{Type: graphql.Int},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"musicas": &graphql.Field{
				Type: graphql.NewList(musicaType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return listarMusicas()
				},
			},
			"musica": &graphql.Field{
				Type: musicaType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)
					return buscarMusicaPorID(id)
				},
			},
			"usuarios": &graphql.Field{
				Type: graphql.NewList(usuarioType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return listarUsuarios()
				},
			},
			"usuario": &graphql.Field{
				Type: usuarioType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)
					return buscarUsuarioPorID(id)
				},
			},
			"playlists": &graphql.Field{
				Type: graphql.NewList(playlistType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return listarPlaylists()
				},
			},
			"playlist": &graphql.Field{
				Type: playlistType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)
					return buscarPlaylistPorID(id)
				},
			},
			"musicasDaPlaylist": &graphql.Field{
				Type: graphql.NewList(musicaType),
				Args: graphql.FieldConfigArgument{
					"playlist_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					playlistID := params.Args["playlist_id"].(int)
					return listarMusicasDaPlaylist(playlistID)
				},
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"criarMusica": &graphql.Field{
				Type: musicaType,
				Args: graphql.FieldConfigArgument{
					"nome":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"artista": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					nome := params.Args["nome"].(string)
					artista := params.Args["artista"].(string)
					return criarMusica(nome, artista)
				},
			},
			"atualizarMusica": &graphql.Field{
				Type: musicaType,
				Args: graphql.FieldConfigArgument{
					"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"nome":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"artista": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)
					nome := params.Args["nome"].(string)
					artista := params.Args["artista"].(string)
					return atualizarMusica(id, nome, artista)
				},
			},
			"deletarMusica": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)
					err := deletarMusica(id)
					if err != nil {
						return nil, err
					}
					return "Música deletada com sucesso", nil
				},
			},

			"criarUsuario": &graphql.Field{
				Type: usuarioType,
				Args: graphql.FieldConfigArgument{
					"nome":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"idade": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					nome := params.Args["nome"].(string)
					idade := params.Args["idade"].(int)
					return criarUsuario(nome, idade)
				},
			},
			"atualizarUsuario": &graphql.Field{
				Type: usuarioType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"nome":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"idade": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)
					nome := params.Args["nome"].(string)
					idade := params.Args["idade"].(int)
					return atualizarUsuario(id, nome, idade)
				},
			},
			"deletarUsuario": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)
					err := deletarUsuario(id)
					if err != nil {
						return nil, err
					}
					return "Usuário deletado com sucesso", nil
				},
			},

			"criarPlaylist": &graphql.Field{
				Type: playlistType,
				Args: graphql.FieldConfigArgument{
					"nome":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"usuario_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					nome := params.Args["nome"].(string)
					usuarioID := params.Args["usuario_id"].(int)
					return criarPlaylist(nome, usuarioID)
				},
			},
			"atualizarPlaylist": &graphql.Field{
				Type: playlistType,
				Args: graphql.FieldConfigArgument{
					"id":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"nome":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"usuario_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)
					nome := params.Args["nome"].(string)
					usuarioID := params.Args["usuario_id"].(int)
					return atualizarPlaylist(id, nome, usuarioID)
				},
			},
			"deletarPlaylist": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)
					err := deletarPlaylist(id)
					if err != nil {
						return nil, err
					}
					return "Playlist deletada com sucesso", nil
				},
			},
			"adicionarMusicaNaPlaylist": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"playlist_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"musica_id":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					playlistID := params.Args["playlist_id"].(int)
					musicaID := params.Args["musica_id"].(int)

					err := adicionarMusicaNaPlaylist(playlistID, musicaID)
					if err != nil {
						return nil, err
					}

					return "Música adicionada na playlist com sucesso", nil
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}

func graphqlHandler(w http.ResponseWriter, r *http.Request, schema graphql.Schema) {
	if r.Method != http.MethodPost {
		responderJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"erro": "Use POST no endpoint /graphql",
		})
		return
	}

	var request GraphQLRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"erro": "JSON inválido",
		})
		return
	}

	resultado := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  request.Query,
		VariableValues: request.Variable,
	})

	responderJSON(w, http.StatusOK, resultado)
}

func listarMusicas() ([]Musica, error) {
	rows, err := conexao.Query("SELECT id, nome, artista FROM musica ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	musicas := []Musica{}

	for rows.Next() {
		var musica Musica

		err := rows.Scan(&musica.ID, &musica.Nome, &musica.Artista)
		if err != nil {
			return nil, err
		}

		musicas = append(musicas, musica)
	}

	return musicas, nil
}

func buscarMusicaPorID(id int) (*Musica, error) {
	var musica Musica

	err := conexao.QueryRow(
		"SELECT id, nome, artista FROM musica WHERE id = $1",
		id,
	).Scan(&musica.ID, &musica.Nome, &musica.Artista)

	if err != nil {
		return nil, err
	}

	return &musica, nil
}

func criarMusica(nome string, artista string) (*Musica, error) {
	var musica Musica

	musica.Nome = nome
	musica.Artista = artista

	err := conexao.QueryRow(
		"INSERT INTO musica (nome, artista) VALUES ($1, $2) RETURNING id",
		musica.Nome,
		musica.Artista,
	).Scan(&musica.ID)

	if err != nil {
		return nil, err
	}

	return &musica, nil
}

func atualizarMusica(id int, nome string, artista string) (*Musica, error) {
	_, err := conexao.Exec(
		"UPDATE musica SET nome = $1, artista = $2 WHERE id = $3",
		nome,
		artista,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &Musica{
		ID:      id,
		Nome:    nome,
		Artista: artista,
	}, nil
}

func deletarMusica(id int) error {
	_, err := conexao.Exec("DELETE FROM musica WHERE id = $1", id)
	return err
}

func listarUsuarios() ([]Usuario, error) {
	rows, err := conexao.Query("SELECT id, nome, idade FROM usuario ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios := []Usuario{}

	for rows.Next() {
		var usuario Usuario

		err := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Idade)
		if err != nil {
			return nil, err
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func buscarUsuarioPorID(id int) (*Usuario, error) {
	var usuario Usuario

	err := conexao.QueryRow(
		"SELECT id, nome, idade FROM usuario WHERE id = $1",
		id,
	).Scan(&usuario.ID, &usuario.Nome, &usuario.Idade)

	if err != nil {
		return nil, err
	}

	return &usuario, nil
}

func criarUsuario(nome string, idade int) (*Usuario, error) {
	var usuario Usuario

	usuario.Nome = nome
	usuario.Idade = idade

	err := conexao.QueryRow(
		"INSERT INTO usuario (nome, idade) VALUES ($1, $2) RETURNING id",
		usuario.Nome,
		usuario.Idade,
	).Scan(&usuario.ID)

	if err != nil {
		return nil, err
	}

	return &usuario, nil
}

func atualizarUsuario(id int, nome string, idade int) (*Usuario, error) {
	_, err := conexao.Exec(
		"UPDATE usuario SET nome = $1, idade = $2 WHERE id = $3",
		nome,
		idade,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &Usuario{
		ID:    id,
		Nome:  nome,
		Idade: idade,
	}, nil
}

func deletarUsuario(id int) error {
	_, err := conexao.Exec("DELETE FROM usuario WHERE id = $1", id)
	return err
}

func listarPlaylists() ([]Playlist, error) {
	rows, err := conexao.Query("SELECT id, nome, usuario_id FROM playlist ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playlists := []Playlist{}

	for rows.Next() {
		var playlist Playlist

		err := rows.Scan(&playlist.ID, &playlist.Nome, &playlist.UsuarioID)
		if err != nil {
			return nil, err
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func buscarPlaylistPorID(id int) (*Playlist, error) {
	var playlist Playlist

	err := conexao.QueryRow(
		"SELECT id, nome, usuario_id FROM playlist WHERE id = $1",
		id,
	).Scan(&playlist.ID, &playlist.Nome, &playlist.UsuarioID)

	if err != nil {
		return nil, err
	}

	return &playlist, nil
}

func criarPlaylist(nome string, usuarioID int) (*Playlist, error) {
	var playlist Playlist

	playlist.Nome = nome
	playlist.UsuarioID = usuarioID

	err := conexao.QueryRow(
		"INSERT INTO playlist (nome, usuario_id) VALUES ($1, $2) RETURNING id",
		playlist.Nome,
		playlist.UsuarioID,
	).Scan(&playlist.ID)

	if err != nil {
		return nil, err
	}

	return &playlist, nil
}

func atualizarPlaylist(id int, nome string, usuarioID int) (*Playlist, error) {
	_, err := conexao.Exec(
		"UPDATE playlist SET nome = $1, usuario_id = $2 WHERE id = $3",
		nome,
		usuarioID,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &Playlist{
		ID:        id,
		Nome:      nome,
		UsuarioID: usuarioID,
	}, nil
}

func deletarPlaylist(id int) error {
	_, err := conexao.Exec("DELETE FROM playlist WHERE id = $1", id)
	return err
}

func listarMusicasDaPlaylist(playlistID int) ([]Musica, error) {
	rows, err := conexao.Query(`
		SELECT m.id, m.nome, m.artista
		FROM musica m
		INNER JOIN playlist_musica pm ON pm.musica_id = m.id
		WHERE pm.playlist_id = $1
		ORDER BY m.id
	`, playlistID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	musicas := []Musica{}

	for rows.Next() {
		var musica Musica

		err := rows.Scan(&musica.ID, &musica.Nome, &musica.Artista)
		if err != nil {
			return nil, err
		}

		musicas = append(musicas, musica)
	}

	return musicas, nil
}

func adicionarMusicaNaPlaylist(playlistID int, musicaID int) error {
	_, err := conexao.Exec(
		"INSERT INTO playlist_musica (playlist_id, musica_id) VALUES ($1, $2)",
		playlistID,
		musicaID,
	)

	return err
}

func responderJSON(w http.ResponseWriter, status int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dados)
}