package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/genda/genda-api/pkg/config"
	_ "github.com/lib/pq"
)

type PostgresDB *sql.DB

func NewConnection() (*sql.DB, error) {
	logger := log.New(os.Stdout, "[internal.storage.postgres] ", log.LstdFlags|log.Lmicroseconds)
	conf := config.New()

	host := conf.PostgresHost
	port := conf.PostgresPort
	user := conf.PostgresUser
	password := conf.PostgresPassword
	dbname := conf.PostgresDatabase
	sslmode := conf.PostgresSsl

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening postgres connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	logger.Println("Connected to postgres")
	return db, nil
}
