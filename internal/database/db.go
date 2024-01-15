package database

import (
	"database/sql"
	"fmt"
)

func OpenDB() *sql.DB {
	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		`localhost`, `ivan`, `992036`, `shorturl`)

	db, err := sql.Open("pgx", ps)
	if err != nil {
		fmt.Errorf("[OpenDB] Не удалось открыть DB: %q", err)
	}
	defer db.Close()

	return db

}
