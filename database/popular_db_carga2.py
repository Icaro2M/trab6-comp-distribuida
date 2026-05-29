import random
import psycopg2
from faker import Faker


DB_CONFIG = {
    "host": "localhost",
    "port": 5432,
    "database": "musicas_db",
    "user": "postgres",
    "password": "postgres"
}


fake = Faker("pt_BR")


QUANTIDADE_USUARIOS = 1000
QUANTIDADE_MUSICAS = 5000
PLAYLISTS_POR_USUARIO = 3
MUSICAS_POR_PLAYLIST = 20


def conectar():
    return psycopg2.connect(**DB_CONFIG)


def limpar_dados(cursor):
    cursor.execute("""
        TRUNCATE TABLE playlist_musica, playlist, musica, usuario
        RESTART IDENTITY CASCADE;
    """)


def inserir_usuarios(cursor):
    usuarios_ids = []

    for _ in range(QUANTIDADE_USUARIOS):
        nome = fake.name()
        idade = random.randint(15, 70)

        cursor.execute("""
            INSERT INTO usuario (nome, idade)
            VALUES (%s, %s)
            RETURNING id;
        """, (nome, idade))

        usuario_id = cursor.fetchone()[0]
        usuarios_ids.append(usuario_id)

    return usuarios_ids


def inserir_musicas(cursor):
    musicas_ids = []

    estilos = [
        "Noite", "Cidade", "Mar", "Luz", "Som",
        "Tempo", "Caminho", "Sonho", "Chuva", "Sol",
        "Horizonte", "Estrada", "Vento", "Fogo", "Lua"
    ]

    artistas = [
        "Banda Aurora", "DJ Horizonte", "Luna Martins",
        "Pedro Vale", "Grupo Atlântico", "Marina Luz",
        "Som Urbano", "Caio Ribeiro", "Vozes do Norte",
        "Ritmo Livre", "Ana Solar", "Banda Eclipse",
        "João Harmonia", "Coletivo Prisma", "Notas do Sul"
    ]

    for i in range(QUANTIDADE_MUSICAS):
        nome = f"{random.choice(estilos)} {fake.word().capitalize()} {i + 1}"
        artista = random.choice(artistas)

        cursor.execute("""
            INSERT INTO musica (nome, artista)
            VALUES (%s, %s)
            RETURNING id;
        """, (nome, artista))

        musica_id = cursor.fetchone()[0]
        musicas_ids.append(musica_id)

    return musicas_ids


def inserir_playlists(cursor, usuarios_ids, musicas_ids):
    for usuario_id in usuarios_ids:
        for i in range(PLAYLISTS_POR_USUARIO):
            nome_playlist = f"Playlist {i + 1} - {fake.word().capitalize()}"

            cursor.execute("""
                INSERT INTO playlist (nome, usuario_id)
                VALUES (%s, %s)
                RETURNING id;
            """, (nome_playlist, usuario_id))

            playlist_id = cursor.fetchone()[0]

            musicas_escolhidas = random.sample(musicas_ids, MUSICAS_POR_PLAYLIST)

            for musica_id in musicas_escolhidas:
                cursor.execute("""
                    INSERT INTO playlist_musica (playlist_id, musica_id)
                    VALUES (%s, %s);
                """, (playlist_id, musica_id))


def popular_banco():
    conn = conectar()
    cursor = conn.cursor()

    limpar_dados(cursor)

    usuarios_ids = inserir_usuarios(cursor)
    musicas_ids = inserir_musicas(cursor)
    inserir_playlists(cursor, usuarios_ids, musicas_ids)

    conn.commit()
    cursor.close()
    conn.close()


if __name__ == "__main__":
    popular_banco()
    print("Banco populado com a carga 2.")