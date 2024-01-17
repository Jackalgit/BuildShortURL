package database

import (
	"context"
	"database/sql"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

type DataBase struct {
	Connect *sql.DB
}

func NewDataBase(ctx context.Context) DataBase {
	ps := config.Config.DatabaseDSN
	query := `CREATE TABLE IF NOT EXISTS storage_url(id int primary key, shortURLKey text, originalURL text)`

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", ps)
	if err != nil {
		log.Printf("[Open DB] Не удалось установить соединение с базой данных: %q", err)
	}

	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("[Create DB] Не удалось создать таблицу в база данных: %q", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Ошибка в получении количества строк: %q", err)
	}
	log.Printf("Количество строк: %d", rows)

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
