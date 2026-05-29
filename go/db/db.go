package db
import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConectarBanco() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=musicas_db sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao abrir conexão com o banco:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}

	return db
}