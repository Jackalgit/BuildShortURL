package database

import (
	"context"
	"database/sql"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

type DataBase struct {
	Connect *sql.DB
}

func NewDataBase() DataBase {
	ps := config.Config.DatabaseDSN

	db, err := sql.Open("pgx", ps)
	if err != nil {
		log.Printf("[Open DB] Не удалось установить соединение с базой данных: %q", err)
	}

	return DataBase{
		Connect: db,
	}
}

func (db DataBase) AddURL(ctx context.Context, shortURLKey string, originalURL []byte) {

}

func (db DataBase) GetURL(ctx context.Context, shortURLKey string) ([]byte, bool) {

	row := db.Connect.QueryRowContext(
		ctx,
		"SELECT originalURL FROM shorturl WHERE short_key = ?", shortURLKey,
	)

	var originalURL sql.NullString
	err := row.Scan(&originalURL)
	if err != nil {
		log.Printf("[row Scan] Не удалось преобразовать данные: %q", err)
	}

	if originalURL.Valid {
		return []byte(originalURL.String), true
	}
	return nil, false

}
