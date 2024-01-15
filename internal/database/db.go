package database

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenDB() (*sql.DB, error) {
	ps := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		`localhost`, `5432`, `ivan`, `XXXXXXXX`, `shorturl`)

	db, err := sql.Open("pgx", ps)
	if err != nil {
		return nil, fmt.Errorf("[OpenDB] Не удалось открыть DB: %q", err)
	}

	return db, nil

}
