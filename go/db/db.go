package db
import (
	"database/sql"
	"log"
	"time"

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

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}
