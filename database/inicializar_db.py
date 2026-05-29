import psycopg2


DB_CONFIG = {
    "host": "localhost",
    "port": 5432,
    "database": "musicas_db",
    "user": "postgres",
    "password": "postgres"
}


def conectar():
    return psycopg2.connect(**DB_CONFIG)


def criar_tabelas():
    conn = conectar()
    cursor = conn.cursor()

    cursor.execute("""
        DROP TABLE IF EXISTS playlist_musica;
        DROP TABLE IF EXISTS playlist;
        DROP TABLE IF EXISTS musica;
        DROP TABLE IF EXISTS usuario;
    """)

    cursor.execute("""
        CREATE TABLE usuario (
            id SERIAL PRIMARY KEY,
            nome VARCHAR(100) NOT NULL,
            idade INT NOT NULL
        );
    """)

    cursor.execute("""
        CREATE TABLE musica (
            id SERIAL PRIMARY KEY,
            nome VARCHAR(100) NOT NULL,
            artista VARCHAR(100) NOT NULL
        );
    """)

    cursor.execute("""
        CREATE TABLE playlist (
            id SERIAL PRIMARY KEY,
            nome VARCHAR(100) NOT NULL,
            usuario_id INT NOT NULL,
            CONSTRAINT fk_playlist_usuario
                FOREIGN KEY (usuario_id)
                REFERENCES usuario(id)
                ON DELETE CASCADE
        );
    """)

    cursor.execute("""
        CREATE TABLE playlist_musica (
            playlist_id INT NOT NULL,
            musica_id INT NOT NULL,

            PRIMARY KEY (playlist_id, musica_id),

            CONSTRAINT fk_playlist_musica_playlist
                FOREIGN KEY (playlist_id)
                REFERENCES playlist(id)
                ON DELETE CASCADE,

            CONSTRAINT fk_playlist_musica_musica
                FOREIGN KEY (musica_id)
                REFERENCES musica(id)
                ON DELETE CASCADE
        );
    """)

    conn.commit()
    cursor.close()
    conn.close()


if __name__ == "__main__":
    criar_tabelas()
    print("Banco inicializado com sucesso.")