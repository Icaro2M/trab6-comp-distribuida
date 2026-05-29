# Trabalho 6 - Computação Distribuída

Este projeto simula um serviço de músicas utilizando um banco de dados PostgreSQL e quatro formas diferentes de acesso aos dados:

* REST
* SOAP
* GraphQL
* gRPC

A parte em Go implementa os quatro serviços acessando o mesmo banco de dados. O banco contém dados de músicas, usuários e playlists.

## Estrutura geral

```txt
trab6-comp-distribuida/
├─ docker-compose.yaml
├─ README.md
├─ database/
│  ├─ inicializar_db.py
│  ├─ popular_db_carga1.py
│  ├─ popular_db_carga2.py
│  └─ requirements.txt
├─ go/
│  ├─ db/
│  │  └─ db.go
│  ├─ rest/
│  │  └─ main.go
│  ├─ soap/
│  │  └─ main.go
│  ├─ graphql/
│  │  └─ main.go
│  └─ grpc/
│     ├─ main.go
│     └─ proto/
│        ├─ music.proto
│        ├─ music.pb.go
│        └─ music_grpc.pb.go
└─ java/
   ├─ rest/
   ├─ graphql/
   ├─ grpc/
   └─ soap/
```

## Banco de dados

O banco utilizado é PostgreSQL, executado via Docker.

As tabelas utilizadas são:

* `usuario`
* `musica`
* `playlist`
* `playlist_musica`

A tabela `playlist_musica` representa a relação entre playlists e músicas.

## Como executar o banco

Na raiz do projeto, execute:

```powershell
docker compose up -d
```

Depois inicialize as tabelas:

```powershell
python database/inicializar_db.py
```

Para popular o banco com uma carga pequena:

```powershell
python database/popular_db_carga1.py
```

Para popular o banco com uma carga maior:

```powershell
python database/popular_db_carga2.py
```

## Como executar os serviços Go

Entre na pasta `go`:

```powershell
cd go
```

Execute cada serviço em um terminal separado.

### REST

```powershell
go run ./rest
```

Endereço:

```txt
http://localhost:8080
```

### SOAP

```powershell
go run ./soap
```

Endereço:

```txt
http://localhost:8081/soap
```

### GraphQL

```powershell
go run ./graphql
```

Endereço:

```txt
http://localhost:8082/graphql
```

### gRPC

```powershell
go run ./grpc
```

Endereço:

```txt
localhost:8083
```

## Portas utilizadas

```txt
REST     -> http://localhost:8080
SOAP     -> http://localhost:8081/soap
GraphQL  -> http://localhost:8082/graphql
gRPC     -> localhost:8083
```

---

# Exemplos de requisições

Os testes podem ser feitos pelo Postman.

## REST

### Listar músicas

```http
GET http://localhost:8080/musicas
```

### Buscar música por ID

```http
GET http://localhost:8080/musicas/1
```

### Criar música

```http
POST http://localhost:8080/musicas
Content-Type: application/json

{
  "nome": "Musica REST",
  "artista": "Artista REST"
}
```

### Atualizar música

```http
PUT http://localhost:8080/musicas/1
Content-Type: application/json

{
  "nome": "Musica REST Atualizada",
  "artista": "Artista REST Atualizado"
}
```

### Deletar música

```http
DELETE http://localhost:8080/musicas/1
```

### Listar usuários

```http
GET http://localhost:8080/usuarios
```

### Buscar usuário por ID

```http
GET http://localhost:8080/usuarios/1
```

### Criar usuário

```http
POST http://localhost:8080/usuarios
Content-Type: application/json

{
  "nome": "Usuario REST",
  "idade": 22
}
```

### Atualizar usuário

```http
PUT http://localhost:8080/usuarios/1
Content-Type: application/json

{
  "nome": "Usuario REST Atualizado",
  "idade": 23
}
```

### Deletar usuário

```http
DELETE http://localhost:8080/usuarios/1
```

### Listar playlists

```http
GET http://localhost:8080/playlists
```

### Buscar playlist por ID

```http
GET http://localhost:8080/playlists/1
```

### Criar playlist

```http
POST http://localhost:8080/playlists
Content-Type: application/json

{
  "nome": "Playlist REST",
  "usuario_id": 1
}
```

### Atualizar playlist

```http
PUT http://localhost:8080/playlists/1
Content-Type: application/json

{
  "nome": "Playlist REST Atualizada",
  "usuario_id": 1
}
```

### Deletar playlist

```http
DELETE http://localhost:8080/playlists/1
```

### Adicionar música na playlist

```http
POST http://localhost:8080/playlists/1/musicas
Content-Type: application/json

{
  "musica_id": 1
}
```

### Listar músicas de uma playlist

```http
GET http://localhost:8080/playlists/1/musicas
```

---

## SOAP

Todas as requisições SOAP utilizam:

```txt
POST http://localhost:8081/soap
Content-Type: text/xml
```

### Listar músicas

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ListarMusicasRequest/>
  </soap:Body>
</soap:Envelope>
```

### Buscar música por ID

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <BuscarMusicaRequest>
      <id>1</id>
    </BuscarMusicaRequest>
  </soap:Body>
</soap:Envelope>
```

### Criar música

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CriarMusicaRequest>
      <nome>Musica SOAP</nome>
      <artista>Artista SOAP</artista>
    </CriarMusicaRequest>
  </soap:Body>
</soap:Envelope>
```

### Atualizar música

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AtualizarMusicaRequest>
      <id>1</id>
      <nome>Musica SOAP Atualizada</nome>
      <artista>Artista SOAP Atualizado</artista>
    </AtualizarMusicaRequest>
  </soap:Body>
</soap:Envelope>
```

### Deletar música

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DeletarMusicaRequest>
      <id>1</id>
    </DeletarMusicaRequest>
  </soap:Body>
</soap:Envelope>
```

### Listar usuários

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ListarUsuariosRequest/>
  </soap:Body>
</soap:Envelope>
```

### Buscar usuário por ID

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <BuscarUsuarioRequest>
      <id>1</id>
    </BuscarUsuarioRequest>
  </soap:Body>
</soap:Envelope>
```

### Criar usuário

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CriarUsuarioRequest>
      <nome>Usuario SOAP</nome>
      <idade>23</idade>
    </CriarUsuarioRequest>
  </soap:Body>
</soap:Envelope>
```

### Atualizar usuário

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AtualizarUsuarioRequest>
      <id>1</id>
      <nome>Usuario SOAP Atualizado</nome>
      <idade>24</idade>
    </AtualizarUsuarioRequest>
  </soap:Body>
</soap:Envelope>
```

### Deletar usuário

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DeletarUsuarioRequest>
      <id>1</id>
    </DeletarUsuarioRequest>
  </soap:Body>
</soap:Envelope>
```

### Listar playlists

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ListarPlaylistsRequest/>
  </soap:Body>
</soap:Envelope>
```

### Buscar playlist por ID

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <BuscarPlaylistRequest>
      <id>1</id>
    </BuscarPlaylistRequest>
  </soap:Body>
</soap:Envelope>
```

### Criar playlist

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CriarPlaylistRequest>
      <nome>Playlist SOAP</nome>
      <usuario_id>1</usuario_id>
    </CriarPlaylistRequest>
  </soap:Body>
</soap:Envelope>
```

### Atualizar playlist

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AtualizarPlaylistRequest>
      <id>1</id>
      <nome>Playlist SOAP Atualizada</nome>
      <usuario_id>1</usuario_id>
    </AtualizarPlaylistRequest>
  </soap:Body>
</soap:Envelope>
```

### Deletar playlist

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DeletarPlaylistRequest>
      <id>1</id>
    </DeletarPlaylistRequest>
  </soap:Body>
</soap:Envelope>
```

### Adicionar música na playlist

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AdicionarMusicaPlaylistRequest>
      <playlist_id>1</playlist_id>
      <musica_id>1</musica_id>
    </AdicionarMusicaPlaylistRequest>
  </soap:Body>
</soap:Envelope>
```

### Listar músicas de uma playlist

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ListarMusicasPlaylistRequest>
      <playlist_id>1</playlist_id>
    </ListarMusicasPlaylistRequest>
  </soap:Body>
</soap:Envelope>
```

---

## GraphQL

Endpoint:

```txt
POST http://localhost:8082/graphql
```

No Postman, selecione a opção **Body > GraphQL**.

### Listar músicas

```graphql
query {
  musicas {
    id
    nome
    artista
  }
}
```

### Buscar música por ID

```graphql
query {
  musica(id: 1) {
    id
    nome
    artista
  }
}
```

### Criar música

```graphql
mutation {
  criarMusica(nome: "Musica GraphQL", artista: "Artista GraphQL") {
    id
    nome
    artista
  }
}
```

### Atualizar música

```graphql
mutation {
  atualizarMusica(id: 1, nome: "Musica GraphQL Atualizada", artista: "Artista GraphQL Atualizado") {
    id
    nome
    artista
  }
}
```

### Deletar música

```graphql
mutation {
  deletarMusica(id: 1)
}
```

### Listar usuários

```graphql
query {
  usuarios {
    id
    nome
    idade
  }
}
```

### Buscar usuário por ID

```graphql
query {
  usuario(id: 1) {
    id
    nome
    idade
  }
}
```

### Criar usuário

```graphql
mutation {
  criarUsuario(nome: "Usuario GraphQL", idade: 25) {
    id
    nome
    idade
  }
}
```

### Atualizar usuário

```graphql
mutation {
  atualizarUsuario(id: 1, nome: "Usuario GraphQL Atualizado", idade: 26) {
    id
    nome
    idade
  }
}
```

### Deletar usuário

```graphql
mutation {
  deletarUsuario(id: 1)
}
```

### Listar playlists

```graphql
query {
  playlists {
    id
    nome
    usuario_id
  }
}
```

### Buscar playlist por ID

```graphql
query {
  playlist(id: 1) {
    id
    nome
    usuario_id
  }
}
```

### Criar playlist

```graphql
mutation {
  criarPlaylist(nome: "Playlist GraphQL", usuario_id: 1) {
    id
    nome
    usuario_id
  }
}
```

### Atualizar playlist

```graphql
mutation {
  atualizarPlaylist(id: 1, nome: "Playlist GraphQL Atualizada", usuario_id: 1) {
    id
    nome
    usuario_id
  }
}
```

### Deletar playlist

```graphql
mutation {
  deletarPlaylist(id: 1)
}
```

### Adicionar música na playlist

```graphql
mutation {
  adicionarMusicaNaPlaylist(playlist_id: 1, musica_id: 1)
}
```

### Listar músicas de uma playlist

```graphql
query {
  musicasDaPlaylist(playlist_id: 1) {
    id
    nome
    artista
  }
}
```

---

## gRPC

O serviço gRPC roda em:

```txt
localhost:8083
```

No Postman:

```txt
New > gRPC Request
```

Depois coloque o endereço:

```txt
localhost:8083
```

Importe o arquivo `.proto`:

```txt
go/grpc/proto/music.proto
```

Serviço:

```txt
music.MusicService
```

### Listar músicas

Método:

```txt
ListMusics
```

Body:

```json
{}
```

### Buscar música por ID

Método:

```txt
GetMusic
```

Body:

```json
{
  "id": 1
}
```

### Criar música

Método:

```txt
CreateMusic
```

Body:

```json
{
  "nome": "Musica gRPC",
  "artista": "Artista gRPC"
}
```

### Atualizar música

Método:

```txt
UpdateMusic
```

Body:

```json
{
  "id": 1,
  "nome": "Musica gRPC Atualizada",
  "artista": "Artista gRPC Atualizado"
}
```

### Deletar música

Método:

```txt
DeleteMusic
```

Body:

```json
{
  "id": 1
}
```

### Listar usuários

Método:

```txt
ListUsers
```

Body:

```json
{}
```

### Buscar usuário por ID

Método:

```txt
GetUser
```

Body:

```json
{
  "id": 1
}
```

### Criar usuário

Método:

```txt
CreateUser
```

Body:

```json
{
  "nome": "Usuario gRPC",
  "idade": 24
}
```

### Atualizar usuário

Método:

```txt
UpdateUser
```

Body:

```json
{
  "id": 1,
  "nome": "Usuario gRPC Atualizado",
  "idade": 25
}
```

### Deletar usuário

Método:

```txt
DeleteUser
```

Body:

```json
{
  "id": 1
}
```

### Listar playlists

Método:

```txt
ListPlaylists
```

Body:

```json
{}
```

### Buscar playlist por ID

Método:

```txt
GetPlaylist
```

Body:

```json
{
  "id": 1
}
```

### Criar playlist

Método:

```txt
CreatePlaylist
```

Body:

```json
{
  "nome": "Playlist gRPC",
  "usuario_id": 1
}
```

### Atualizar playlist

Método:

```txt
UpdatePlaylist
```

Body:

```json
{
  "id": 1,
  "nome": "Playlist gRPC Atualizada",
  "usuario_id": 1
}
```

### Deletar playlist

Método:

```txt
DeletePlaylist
```

Body:

```json
{
  "id": 1
}
```

### Adicionar música na playlist

Método:

```txt
AddMusicToPlaylist
```

Body:

```json
{
  "playlist_id": 1,
  "musica_id": 1
}
```

### Listar músicas de uma playlist

Método:

```txt
ListPlaylistMusics
```

Body:

```json
{
  "playlist_id": 1
}
```

---

## Observação sobre os IDs

Os exemplos utilizam IDs como `1`, mas esses valores dependem dos dados existentes no banco.

Após rodar o script `popular_db_carga1.py`, já existirão usuários, músicas e playlists cadastrados. Caso uma requisição retorne erro de item não encontrado, verifique primeiro os registros existentes usando as rotas de listagem.

## Ordem recomendada para testes

```txt
1. Listar usuários
2. Listar músicas
3. Listar playlists
4. Criar uma nova música
5. Criar um novo usuário
6. Criar uma nova playlist
7. Adicionar uma música a uma playlist
8. Listar músicas da playlist
```
